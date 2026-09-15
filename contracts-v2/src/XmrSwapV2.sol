// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Ed25519} from "./lib/Ed25519.sol";

/// @notice Minimal Chainlink-compatible price feed interface.
/// @dev A compatible feed reports how much Monero one ETH buys. Any `decimals()` is accepted and
///      normalised to 12 (4.65 XMR per ETH == 4.65e12 piconero per 1e18 wei).
interface IAggregatorV3 {
    function decimals() external view returns (uint8);
    function latestRoundData()
        external
        view
        returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);
}

/// @title XmrSwapV2
/// @notice Trustless Monero <-> ETH atomic swaps with the offer board on-chain.
///
/// Roles: the *maker* holds Monero and posts an offer. The *taker* pays ETH by taking it,
/// which locks the payment here. The maker sends Monero to a one-time address controlled
/// jointly by both parties' one-time keys. Once the taker sees the Monero it marks the swap
/// ready; the maker collects the ETH by revealing its one-time Monero secret, and the taker
/// uses that secret to spend the Monero. If anything stalls, timeouts return funds.
///
/// The contract has no owner and no upgrade path. Fee, split and bond are fixed at deployment.
///
/// Changes from v1:
///  - the fee splits between two immutable recipients and *accrues* rather than being pushed,
///    so a recipient can never brick a claim;
///  - offers are enumerable from storage, so discovery no longer depends on log scanning;
///  - makers post a bond to list, held and always returned, never seized;
///  - minimum timeout is 24h, long enough to survive an L2 sequencer stall;
///  - ETH only. An ERC-20 issuer can blacklist a contract and strand every swap inside it;
///  - offers may optionally price from a feed instead of a fixed number.
contract XmrSwapV2 is ReentrancyGuard {
    // ---------------------------------------------------------------- constants

    uint16 public constant BPS = 10_000;
    uint16 public constant MAX_FEE_BPS = 100; // hard cap: 1%
    uint16 public constant MAX_BOND_BPS = 2_000; // hard cap: 20% of an offer's ceiling
    /// @dev 24h. Base has one sequencer; forcing a transaction in from L1 takes ~12h, and a
    ///      shorter window can close before a censored user can act.
    uint64 public constant MIN_TIMEOUT = 24 hours;
    /// @dev Units of `xmrPerAsset`: piconero per 1e18 wei.
    uint8 public constant PRICE_DECIMALS = 12;

    /// @notice Total fee taken from the taker's payment when the maker claims, in basis points.
    uint16 public immutable feeBps;
    /// @notice Share of that fee going to `founder`, in basis points. The remainder goes to `devFund`.
    uint16 public immutable founderShareBps;
    /// @notice Bond a maker locks to list an offer, in basis points of the offer's ceiling.
    uint16 public immutable bondBps;
    /// @notice Receives `founderShareBps` of every fee. Immutable: it can never be redirected.
    address payable public immutable founder;
    /// @notice Receives the rest of every fee. Immutable.
    address payable public immutable devFund;

    // ---------------------------------------------------------------- types

    enum Stage {
        INVALID,
        PENDING,
        READY,
        COMPLETED
    }

    /// @dev One offer = one swap. Keys are one-time; the maker posts a fresh offer after each take.
    struct Offer {
        bytes32 id;
        address maker; // posted it, paid the bond, may cancel and may claim
        address payable payout; // receives the payment; any wallet the maker chooses
        uint128 minAmount; // smallest payment a taker may lock, in wei
        uint128 maxAmount; // largest payment a taker may lock, in wei
        uint256 xmrPerAsset; // fixed price, piconero per 1e18 wei. Zero means price from `oracle`.
        uint64 expiry; // cannot be taken at or after this timestamp
        uint64 timeout1Duration; // from take until the maker may claim without ready
        uint64 timeout2Duration; // after timeout1 until the taker may refund regardless
        bytes32 makerSpendPub; // maker's one-time Monero public spend key (Monero byte order)
        bytes32 makerViewPriv; // maker's one-time Monero private view key
        uint256 bond; // locked by the maker, always returned, never seized
        address oracle; // zero for a fixed price
        uint256 oracleRatioBps; // price = feed * ratio / BPS, then + offset
        int256 oracleOffset; // in piconero per 1e18 wei
        uint64 oracleMaxAge; // reject a feed reading older than this
        uint256 oracleMaxPrice; // maker's ceiling on the resolved price; zero means none
        uint256 index; // position in `offerIds`
        bool active;
    }

    struct Swap {
        address payable taker; // may refund; receives the refund
        address maker; // may claim; gets the bond back either way
        address payable payout; // receives the payment
        uint256 value; // wei locked
        uint256 bond; // the maker's bond, carried through settlement
        bytes32 makerSpendPub; // claim commitment: claim(secret) must derive this key
        bytes32 takerSpendPub; // refund commitment: refund(secret) must derive this key
        uint64 timeout1;
        uint64 timeout2;
        Stage stage;
    }

    /// @dev Grouped so `postOffer` does not exhaust the stack.
    struct OfferParams {
        uint128 minAmount;
        uint128 maxAmount;
        uint256 xmrPerAsset;
        uint64 expiry;
        uint64 timeout1Duration;
        uint64 timeout2Duration;
        bytes32 makerSpendPub;
        bytes32 makerViewPriv;
        address payable payout;
        address oracle;
        uint256 oracleRatioBps;
        int256 oracleOffset;
        uint64 oracleMaxAge;
        uint256 oracleMaxPrice;
    }

    // ---------------------------------------------------------------- state

    /// @dev internal: a 19-field auto-getter exhausts the stack. Read with `getOffer`.
    mapping(bytes32 offerId => Offer) internal offers;
    mapping(bytes32 swapId => Swap) public swaps;
    /// @dev One-time keys must never be reused: once a swap reveals a secret, a second swap
    ///      built on the same key would let the counterparty spend the Monero without paying.
    mapping(bytes32 spendPub => bool) public keyUsed;
    /// @dev Money owed but not yet pushed: fee shares and returned bonds. Pulling instead of
    ///      pushing means nothing at a recipient can make a claim or a refund fail.
    mapping(address account => uint256) public withdrawable;

    /// @notice Every live offer, in no particular order. Read with `listOffers`.
    bytes32[] public offerIds;

    uint256 private _nonce;

    // ---------------------------------------------------------------- events

    event OfferPosted(
        bytes32 indexed offerId,
        address indexed maker,
        uint256 minAmount,
        uint256 maxAmount,
        uint256 xmrPerAsset,
        address oracle,
        uint256 bond,
        uint64 expiry
    );
    event OfferClosed(bytes32 indexed offerId, bool taken);
    event SwapCreated(
        bytes32 indexed swapId,
        bytes32 indexed offerId,
        address indexed taker,
        uint256 value,
        uint256 xmrPiconero,
        bytes32 takerSpendPub,
        bytes32 takerViewPriv,
        uint64 timeout1,
        uint64 timeout2
    );
    event Ready(bytes32 indexed swapId);
    event Claimed(bytes32 indexed swapId, bytes32 secret, uint256 payout, uint256 fee);
    event Refunded(bytes32 indexed swapId, bytes32 secret);
    event Withdrawn(address indexed account, uint256 amount);

    // ---------------------------------------------------------------- errors

    error FeeTooHigh();
    error BondTooHigh();
    error BadShare();
    error ZeroAddress();
    error ZeroKey();
    error BadAmounts();
    error TimeoutTooShort();
    error OfferNotActive();
    error OfferExpired();
    error OfferNotExpired();
    error NotMaker();
    error NotTaker();
    error AmountOutOfRange();
    error WrongValue();
    error WrongBond();
    error InvalidSwap();
    error SwapCompleted();
    error AlreadyReady();
    error TooEarlyToClaim();
    error TooLateToClaim();
    error NotTimeToRefund();
    error InvalidSecret();
    error KeyAlreadyUsed();
    error BadOracle();
    error OracleStale();
    error PriceTooHigh();
    error NotEnoughXmr();
    error NothingToWithdraw();
    error TransferFailed();

    // ---------------------------------------------------------------- constructor

    constructor(
        uint16 feeBps_,
        uint16 founderShareBps_,
        uint16 bondBps_,
        address payable founder_,
        address payable devFund_
    ) {
        if (feeBps_ > MAX_FEE_BPS) revert FeeTooHigh();
        if (bondBps_ > MAX_BOND_BPS) revert BondTooHigh();
        if (founderShareBps_ > BPS) revert BadShare();
        if (feeBps_ != 0 && (founder_ == address(0) || devFund_ == address(0))) revert ZeroAddress();
        feeBps = feeBps_;
        founderShareBps = founderShareBps_;
        bondBps = bondBps_;
        founder = founder_;
        devFund = devFund_;
    }

    // ---------------------------------------------------------------- offers

    /// @notice Post a single-use offer to sell Monero for ETH. Send exactly `bondFor(maxAmount)`.
    /// @dev The bond is held while the offer is live and returned on cancel, expiry or settlement.
    ///      It is never paid to anyone else: this contract cannot see Monero, so it cannot tell a
    ///      non-delivering maker from a taker who bailed out early.
    function postOffer(OfferParams calldata p) external payable nonReentrant returns (bytes32 offerId) {
        if (p.payout == address(0)) revert ZeroAddress();
        if (p.makerSpendPub == bytes32(0) || p.makerViewPriv == bytes32(0)) revert ZeroKey();
        if (p.minAmount == 0 || p.maxAmount < p.minAmount) revert BadAmounts();
        if (p.timeout1Duration < MIN_TIMEOUT || p.timeout2Duration < MIN_TIMEOUT) revert TimeoutTooShort();
        if (p.expiry <= block.timestamp) revert OfferExpired();
        if (keyUsed[p.makerSpendPub]) revert KeyAlreadyUsed();

        // Exactly one pricing source.
        if (p.oracle == address(0)) {
            if (p.xmrPerAsset == 0) revert BadAmounts();
        } else {
            if (p.xmrPerAsset != 0 || p.oracleRatioBps == 0 || p.oracleMaxAge == 0) revert BadOracle();
            // Fail now rather than at take time if the feed is unusable.
            _readOracle(p.oracle, p.oracleRatioBps, p.oracleOffset, p.oracleMaxAge);
        }

        uint256 required = bondFor(p.maxAmount);
        if (msg.value != required) revert WrongBond();

        keyUsed[p.makerSpendPub] = true;
        offerId = keccak256(abi.encodePacked(block.chainid, address(this), msg.sender, _nonce++));

        Offer storage o = offers[offerId];
        o.id = offerId;
        o.maker = msg.sender;
        o.payout = p.payout;
        o.minAmount = p.minAmount;
        o.maxAmount = p.maxAmount;
        o.xmrPerAsset = p.xmrPerAsset;
        o.expiry = p.expiry;
        o.timeout1Duration = p.timeout1Duration;
        o.timeout2Duration = p.timeout2Duration;
        o.makerSpendPub = p.makerSpendPub;
        o.makerViewPriv = p.makerViewPriv;
        o.bond = required;
        o.oracle = p.oracle;
        o.oracleRatioBps = p.oracleRatioBps;
        o.oracleOffset = p.oracleOffset;
        o.oracleMaxAge = p.oracleMaxAge;
        o.oracleMaxPrice = p.oracleMaxPrice;
        o.index = offerIds.length;
        o.active = true;
        offerIds.push(offerId);

        emit OfferPosted(
            offerId, msg.sender, p.minAmount, p.maxAmount, p.xmrPerAsset, p.oracle, required, p.expiry
        );
    }

    /// @notice Withdraw an offer that has not been taken. The bond becomes withdrawable.
    function cancelOffer(bytes32 offerId) external {
        Offer storage o = offers[offerId];
        if (!o.active) revert OfferNotActive();
        if (o.maker != msg.sender) revert NotMaker();
        _closeOffer(o, false);
        withdrawable[o.maker] += o.bond;
    }

    /// @notice Clear an expired offer off the board. Anyone may call; the bond still goes to the maker.
    function reapExpired(bytes32 offerId) external {
        Offer storage o = offers[offerId];
        if (!o.active) revert OfferNotActive();
        if (block.timestamp < o.expiry) revert OfferNotExpired();
        _closeOffer(o, false);
        withdrawable[o.maker] += o.bond;
    }

    // ---------------------------------------------------------------- swaps

    /// @notice Take an offer by locking `msg.value` wei against it.
    /// @param minXmrPiconero the least Monero the taker will accept for that payment. Protects the
    ///        taker when the offer is priced from a feed that can move between quote and mining.
    function takeOffer(
        bytes32 offerId,
        bytes32 takerSpendPub,
        bytes32 takerViewPriv,
        uint256 minXmrPiconero
    ) external payable nonReentrant returns (bytes32 swapId) {
        Offer storage o = offers[offerId];
        if (!o.active) revert OfferNotActive();
        if (block.timestamp >= o.expiry) revert OfferExpired();
        if (takerSpendPub == bytes32(0) || takerViewPriv == bytes32(0)) revert ZeroKey();
        if (msg.value < o.minAmount || msg.value > o.maxAmount) revert AmountOutOfRange();
        if (keyUsed[takerSpendPub]) revert KeyAlreadyUsed();

        uint256 price = _priceOf(o);
        uint256 xmrPiconero = (msg.value * price) / 1e18;
        if (xmrPiconero < minXmrPiconero || xmrPiconero == 0) revert NotEnoughXmr();

        keyUsed[takerSpendPub] = true;

        swapId = keccak256(abi.encodePacked(offerId, msg.sender, takerSpendPub));
        if (swaps[swapId].stage != Stage.INVALID) revert InvalidSwap();

        uint64 t1 = uint64(block.timestamp) + o.timeout1Duration;
        uint64 t2 = t1 + o.timeout2Duration;

        swaps[swapId] = Swap({
            taker: payable(msg.sender),
            maker: o.maker,
            payout: o.payout,
            value: msg.value,
            bond: o.bond,
            makerSpendPub: o.makerSpendPub,
            takerSpendPub: takerSpendPub,
            timeout1: t1,
            timeout2: t2,
            stage: Stage.PENDING
        });

        _closeOffer(o, true);

        emit SwapCreated(swapId, offerId, msg.sender, msg.value, xmrPiconero, takerSpendPub, takerViewPriv, t1, t2);
    }

    /// @notice Taker confirms the Monero arrived. From now until timeout2 the maker may claim
    ///         and the taker may not refund.
    function setReady(bytes32 swapId) external {
        Swap storage s = swaps[swapId];
        if (s.stage == Stage.INVALID) revert InvalidSwap();
        if (s.stage == Stage.COMPLETED) revert SwapCompleted();
        if (s.stage == Stage.READY) revert AlreadyReady();
        if (s.taker != msg.sender) revert NotTaker();
        if (block.timestamp >= s.timeout1) revert TooLateToClaim();
        s.stage = Stage.READY;
        emit Ready(swapId);
    }

    /// @notice Maker collects the payment by revealing its one-time Monero secret spend key.
    ///         Allowed once the taker set ready, or after timeout1; never after timeout2.
    function claim(bytes32 swapId, bytes32 secret) external nonReentrant {
        Swap storage s = swaps[swapId];
        if (s.stage == Stage.INVALID) revert InvalidSwap();
        if (s.stage == Stage.COMPLETED) revert SwapCompleted();
        if (s.maker != msg.sender) revert NotMaker();
        if (block.timestamp < s.timeout1 && s.stage != Stage.READY) revert TooEarlyToClaim();
        if (block.timestamp >= s.timeout2) revert TooLateToClaim();
        if (!_secretMatches(secret, s.makerSpendPub)) revert InvalidSecret();

        s.stage = Stage.COMPLETED;

        uint256 fee = (s.value * feeBps) / BPS;
        uint256 payout = s.value - fee;

        // Fees and the bond accrue. Only the maker's own payout address is pushed, so nothing
        // outside the maker's control can make this revert.
        if (fee != 0) {
            uint256 founderCut = (fee * founderShareBps) / BPS;
            withdrawable[founder] += founderCut;
            withdrawable[devFund] += fee - founderCut;
        }
        withdrawable[s.maker] += s.bond;

        _send(s.payout, payout);
        emit Claimed(swapId, secret, payout, fee);
    }

    /// @notice Taker takes the payment back by revealing its one-time Monero secret spend key,
    ///         which lets the maker recover any Monero it locked. No fee is charged on a refund,
    ///         and the maker's bond is returned in full.
    function refund(bytes32 swapId, bytes32 secret) external nonReentrant {
        Swap storage s = swaps[swapId];
        if (s.stage == Stage.INVALID) revert InvalidSwap();
        if (s.stage == Stage.COMPLETED) revert SwapCompleted();
        if (s.taker != msg.sender) revert NotTaker();
        if (block.timestamp < s.timeout2 && (block.timestamp >= s.timeout1 || s.stage == Stage.READY)) {
            revert NotTimeToRefund();
        }
        if (!_secretMatches(secret, s.takerSpendPub)) revert InvalidSecret();

        s.stage = Stage.COMPLETED;
        withdrawable[s.maker] += s.bond;

        _send(s.taker, s.value);
        emit Refunded(swapId, secret);
    }

    // ---------------------------------------------------------------- money out

    /// @notice Take whatever is owed to you: fee shares and returned bonds.
    function withdraw() external nonReentrant {
        _withdrawTo(msg.sender);
    }

    /// @notice Push whatever is owed to `account`. Anyone may call, so the dev fund can be
    ///         settled without the Safe having to send a transaction itself.
    function withdrawFor(address account) external nonReentrant {
        _withdrawTo(account);
    }

    // ---------------------------------------------------------------- views

    /// @notice Number of live offers.
    function offerCount() external view returns (uint256) {
        return offerIds.length;
    }

    /// @notice A page of live offers, straight from storage. No log scanning, so an offer stays
    ///         findable for its whole life rather than only while it is recent.
    function listOffers(uint256 offset, uint256 count) external view returns (Offer[] memory page) {
        uint256 n = offerIds.length;
        if (offset >= n) return new Offer[](0);
        uint256 end = offset + count;
        if (end > n) end = n;
        page = new Offer[](end - offset);
        for (uint256 i = offset; i < end; i++) {
            page[i - offset] = offers[offerIds[i]];
        }
    }

    /// @notice Read a single offer, live or settled.
    function getOffer(bytes32 offerId) external view returns (Offer memory) {
        return offers[offerId];
    }

    /// @notice The bond required to list an offer whose ceiling is `maxAmount`.
    function bondFor(uint256 maxAmount) public view returns (uint256) {
        return (maxAmount * bondBps) / BPS;
    }

    /// @notice The price an offer would resolve to right now, in piconero per 1e18 wei.
    function priceOf(bytes32 offerId) external view returns (uint256) {
        return _priceOf(offers[offerId]);
    }

    /// @notice True if `secret` (a Monero private spend key, Monero byte order) derives `pubKey`.
    function secretMatches(bytes32 secret, bytes32 pubKey) external view returns (bool) {
        return _secretMatches(secret, pubKey);
    }

    // ---------------------------------------------------------------- internals

    function _priceOf(Offer storage o) internal view returns (uint256) {
        if (o.oracle == address(0)) return o.xmrPerAsset;
        uint256 price = _readOracle(o.oracle, o.oracleRatioBps, o.oracleOffset, o.oracleMaxAge);
        if (o.oracleMaxPrice != 0 && price > o.oracleMaxPrice) revert PriceTooHigh();
        return price;
    }

    /// @dev Reads a feed, normalises it to PRICE_DECIMALS, then applies the offer's ratio and offset.
    function _readOracle(address oracle, uint256 ratioBps, int256 offset, uint64 maxAge)
        internal
        view
        returns (uint256 price)
    {
        (, int256 answer,, uint256 updatedAt,) = IAggregatorV3(oracle).latestRoundData();
        if (answer <= 0 || updatedAt == 0) revert BadOracle();
        if (block.timestamp > updatedAt + maxAge) revert OracleStale();

        uint256 raw = uint256(answer);
        uint8 d = IAggregatorV3(oracle).decimals();
        if (d < PRICE_DECIMALS) raw *= 10 ** (PRICE_DECIMALS - d);
        else if (d > PRICE_DECIMALS) raw /= 10 ** (d - PRICE_DECIMALS);

        price = (raw * ratioBps) / BPS;
        if (offset >= 0) price += uint256(offset);
        else {
            uint256 down = uint256(-offset);
            if (down >= price) revert BadOracle();
            price -= down;
        }
        if (price == 0) revert BadOracle();
    }

    function _closeOffer(Offer storage o, bool taken) internal {
        o.active = false;
        uint256 i = o.index;
        uint256 last = offerIds.length - 1;
        if (i != last) {
            bytes32 moved = offerIds[last];
            offerIds[i] = moved;
            offers[moved].index = i;
        }
        offerIds.pop();
        emit OfferClosed(o.id, taken);
    }

    function _withdrawTo(address account) internal {
        uint256 amount = withdrawable[account];
        if (amount == 0) revert NothingToWithdraw();
        withdrawable[account] = 0;
        _send(payable(account), amount);
        emit Withdrawn(account, amount);
    }

    function _secretMatches(bytes32 secret, bytes32 pubKey) internal view returns (bool) {
        // Monero encodes scalars little-endian; the library wants a plain integer.
        uint256 s = Ed25519.changeEndianness(uint256(secret));
        if (s == 0) return false;
        (uint256 x, uint256 y) = Ed25519.scalarMultBase(s);
        return bytes32(Ed25519.compressPoint(x, y)) == pubKey;
    }

    function _send(address payable to, uint256 amount) internal {
        (bool ok,) = to.call{value: amount}("");
        if (!ok) revert TransferFailed();
    }
}
