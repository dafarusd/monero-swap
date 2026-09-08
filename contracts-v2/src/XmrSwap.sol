// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Ed25519} from "./lib/Ed25519.sol";

/// @title XmrSwap
/// @notice Trustless Monero <-> ETH / ERC-20 atomic swaps with the offer board on-chain.
///
/// Roles: the *maker* holds Monero and posts an offer. The *taker* pays ETH (or a token)
/// by taking the offer, which locks the payment in this contract. The maker then sends
/// Monero to a one-time address controlled jointly by both parties' one-time keys.
/// Once the taker sees the Monero, it marks the swap ready; the maker collects the
/// payment by revealing its one-time Monero secret, and the taker uses that secret to
/// spend the Monero. If anything stalls, timeouts return funds to their owners.
///
/// The contract has no owner and no upgrade path. The fee is fixed at deployment.
contract XmrSwap is ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ---------------------------------------------------------------- constants

    uint16 public constant BPS = 10_000;
    uint16 public constant MAX_FEE_BPS = 100; // hard cap: 1%
    uint64 public constant MIN_TIMEOUT = 1 hours;

    /// @notice Fee taken from the taker's payment when the maker claims, in basis points.
    uint16 public immutable feeBps;
    /// @notice Wallet that receives the fee.
    address payable public immutable feeRecipient;

    // ---------------------------------------------------------------- types

    enum Stage {
        INVALID,
        PENDING,
        READY,
        COMPLETED
    }

    /// @dev One offer = one swap. Keys are one-time; the maker posts a fresh offer after each take.
    struct Offer {
        address maker;              // posted the offer; may cancel it and claim the swap
        address payable payout;     // receives the payment (any wallet the maker chooses)
        address asset;              // address(0) = ETH, otherwise ERC-20
        uint128 minAmount;          // smallest payment the taker may lock, in asset base units
        uint128 maxAmount;          // largest payment the taker may lock
        uint256 xmrPerAsset;        // piconero the taker receives per 1e18 base units of asset
        uint64 expiry;              // offer cannot be taken after this timestamp
        uint64 timeout1Duration;    // seconds from take until the maker may claim without ready
        uint64 timeout2Duration;    // seconds after timeout1 until the taker may refund regardless
        bytes32 makerSpendPub;      // maker's one-time Monero public spend key (Monero byte order)
        bytes32 makerViewPriv;      // maker's one-time Monero private view key
        bool active;
    }

    struct Swap {
        address payable taker;      // may refund; receives the refund
        address maker;              // may claim
        address payable payout;     // receives the payment
        address asset;
        uint256 value;              // amount locked
        bytes32 makerSpendPub;      // claim commitment: claim(secret) must derive this key
        bytes32 takerSpendPub;      // refund commitment: refund(secret) must derive this key
        uint64 timeout1;
        uint64 timeout2;
        Stage stage;
    }

    // ---------------------------------------------------------------- state

    mapping(bytes32 offerId => Offer) public offers;
    mapping(bytes32 swapId => Swap) public swaps;
    /// @dev One-time keys must never be reused: once a swap reveals a secret, a second swap
    ///      built on the same key would let the counterparty spend the Monero without paying.
    mapping(bytes32 spendPub => bool) public keyUsed;
    uint256 private _nonce;

    // ---------------------------------------------------------------- events

    /// @dev Keys, timeouts and payout are in `offers(offerId)`.
    event OfferPosted(
        bytes32 indexed offerId,
        address indexed maker,
        address indexed asset,
        uint256 minAmount,
        uint256 maxAmount,
        uint256 xmrPerAsset,
        uint64 expiry
    );
    event OfferCancelled(bytes32 indexed offerId);
    /// @dev Maker key material and asset are in the offer (still readable after the take).
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

    // ---------------------------------------------------------------- errors

    error FeeTooHigh();
    error ZeroAddress();
    error ZeroKey();
    error BadAmounts();
    error TimeoutTooShort();
    error OfferNotActive();
    error OfferExpired();
    error NotMaker();
    error NotTaker();
    error AmountOutOfRange();
    error WrongValue();
    error InvalidSwap();
    error SwapCompleted();
    error AlreadyReady();
    error TooEarlyToClaim();
    error TooLateToClaim();
    error NotTimeToRefund();
    error InvalidSecret();
    error KeyAlreadyUsed();

    // ---------------------------------------------------------------- constructor

    constructor(uint16 feeBps_, address payable feeRecipient_) {
        if (feeBps_ > MAX_FEE_BPS) revert FeeTooHigh();
        if (feeBps_ != 0 && feeRecipient_ == address(0)) revert ZeroAddress();
        feeBps = feeBps_;
        feeRecipient = feeRecipient_;
    }

    // ---------------------------------------------------------------- offers

    /// @notice Post a single-use offer to sell Monero for `asset`.
    function postOffer(
        address asset,
        uint128 minAmount,
        uint128 maxAmount,
        uint256 xmrPerAsset,
        uint64 expiry,
        uint64 timeout1Duration,
        uint64 timeout2Duration,
        bytes32 makerSpendPub,
        bytes32 makerViewPriv,
        address payable payout
    ) external returns (bytes32 offerId) {
        if (payout == address(0)) revert ZeroAddress();
        if (makerSpendPub == bytes32(0) || makerViewPriv == bytes32(0)) revert ZeroKey();
        if (minAmount == 0 || maxAmount < minAmount || xmrPerAsset == 0) revert BadAmounts();
        if (timeout1Duration < MIN_TIMEOUT || timeout2Duration < MIN_TIMEOUT) revert TimeoutTooShort();
        if (expiry <= block.timestamp) revert OfferExpired();
        if (keyUsed[makerSpendPub]) revert KeyAlreadyUsed();
        keyUsed[makerSpendPub] = true;

        offerId = keccak256(abi.encodePacked(block.chainid, address(this), msg.sender, _nonce++));
        offers[offerId] = Offer({
            maker: msg.sender,
            payout: payout,
            asset: asset,
            minAmount: minAmount,
            maxAmount: maxAmount,
            xmrPerAsset: xmrPerAsset,
            expiry: expiry,
            timeout1Duration: timeout1Duration,
            timeout2Duration: timeout2Duration,
            makerSpendPub: makerSpendPub,
            makerViewPriv: makerViewPriv,
            active: true
        });

        emit OfferPosted(offerId, msg.sender, asset, minAmount, maxAmount, xmrPerAsset, expiry);
    }

    /// @notice Withdraw an offer that has not been taken.
    function cancelOffer(bytes32 offerId) external {
        Offer storage o = offers[offerId];
        if (!o.active) revert OfferNotActive();
        if (o.maker != msg.sender) revert NotMaker();
        o.active = false;
        emit OfferCancelled(offerId);
    }

    // ---------------------------------------------------------------- swaps

    /// @notice Take an offer by locking `amount` of its asset. For ETH, send `amount` as msg.value.
    /// @param takerSpendPub  taker's one-time Monero public spend key (Monero byte order)
    /// @param takerViewPriv  taker's one-time Monero private view key
    function takeOffer(
        bytes32 offerId,
        uint256 amount,
        bytes32 takerSpendPub,
        bytes32 takerViewPriv
    ) external payable nonReentrant returns (bytes32 swapId) {
        Offer storage o = offers[offerId];
        if (!o.active) revert OfferNotActive();
        if (block.timestamp >= o.expiry) revert OfferExpired();
        if (takerSpendPub == bytes32(0) || takerViewPriv == bytes32(0)) revert ZeroKey();
        if (amount < o.minAmount || amount > o.maxAmount) revert AmountOutOfRange();
        if (keyUsed[takerSpendPub]) revert KeyAlreadyUsed();
        keyUsed[takerSpendPub] = true;

        if (o.asset == address(0)) {
            if (msg.value != amount) revert WrongValue();
        } else {
            if (msg.value != 0) revert WrongValue();
            IERC20(o.asset).safeTransferFrom(msg.sender, address(this), amount);
        }

        o.active = false;

        uint64 t1 = uint64(block.timestamp) + o.timeout1Duration;
        uint64 t2 = t1 + o.timeout2Duration;
        swapId = keccak256(abi.encodePacked(offerId, msg.sender, takerSpendPub));
        if (swaps[swapId].stage != Stage.INVALID) revert InvalidSwap();

        swaps[swapId] = Swap({
            taker: payable(msg.sender),
            maker: o.maker,
            payout: o.payout,
            asset: o.asset,
            value: amount,
            makerSpendPub: o.makerSpendPub,
            takerSpendPub: takerSpendPub,
            timeout1: t1,
            timeout2: t2,
            stage: Stage.PENDING
        });

        uint256 xmrPiconero = (amount * o.xmrPerAsset) / 1e18;

        emit SwapCreated(swapId, offerId, msg.sender, amount, xmrPiconero, takerSpendPub, takerViewPriv, t1, t2);
    }

    /// @notice Taker confirms the Monero arrived. From now until timeout2 the maker may claim
    ///         and the taker may not refund.
    function setReady(bytes32 swapId) external {
        Swap storage s = swaps[swapId];
        if (s.stage == Stage.INVALID) revert InvalidSwap();
        if (s.stage == Stage.COMPLETED) revert SwapCompleted();
        if (s.stage == Stage.READY) revert AlreadyReady();
        if (s.taker != msg.sender) revert NotTaker();
        if (block.timestamp >= s.timeout1) revert TooLateToClaim(); // pointless after t1; keep state simple
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

        _pay(s.asset, s.payout, payout);
        if (fee != 0) _pay(s.asset, feeRecipient, fee);

        emit Claimed(swapId, secret, payout, fee);
    }

    /// @notice Taker takes the payment back by revealing its one-time Monero secret spend key,
    ///         which lets the maker recover any Monero it locked.
    ///         Allowed before timeout1 unless ready was set, and always after timeout2.
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
        _pay(s.asset, s.taker, s.value);

        emit Refunded(swapId, secret);
    }

    // ---------------------------------------------------------------- views

    /// @notice True if `secret` (a Monero private spend key, Monero byte order) derives `pubKey`.
    function secretMatches(bytes32 secret, bytes32 pubKey) external view returns (bool) {
        return _secretMatches(secret, pubKey);
    }

    // ---------------------------------------------------------------- internals

    function _secretMatches(bytes32 secret, bytes32 pubKey) internal view returns (bool) {
        // Monero encodes scalars little-endian; the library wants a plain integer.
        uint256 s = Ed25519.changeEndianness(uint256(secret));
        if (s == 0) return false;
        (uint256 x, uint256 y) = Ed25519.scalarMultBase(s);
        return bytes32(Ed25519.compressPoint(x, y)) == pubKey;
    }

    function _pay(address asset, address payable to, uint256 amount) internal {
        if (asset == address(0)) {
            (bool ok, ) = to.call{value: amount}("");
            require(ok, "ETH transfer failed");
        } else {
            IERC20(asset).safeTransfer(to, amount);
        }
    }
}
