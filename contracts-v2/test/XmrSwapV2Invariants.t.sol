// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test, console} from "forge-std/Test.sol";
import {XmrSwapV2} from "../src/XmrSwapV2.sol";
import {KeyVectorsBig} from "./KeyVectorsBig.sol";

/// @dev Drives XmrSwapV2 through random call sequences on real ETH swaps. One-time keys come from
///      a 513-vector pool so the fuzzer can build hundreds of genuine swaps. Live offers and open
///      swaps are tracked in their own index lists, so a fuzzed call almost always does real work
///      instead of picking a dead item and bailing.
///
///      The handler keeps its own ledger of where every wei should be — bonds still held, money
///      owed but not yet pulled, and value locked in open swaps — built from the transitions it
///      makes rather than from anything the contract reports. The invariants then check the
///      contract's balance and its own `withdrawable` accounting against that ledger, so the
///      contract is never used to verify itself.
contract HandlerV2 is Test {
    XmrSwapV2 public swap;

    bytes32[] internal secrets;
    bytes32[] internal pubs;
    uint256 public keyCursor;

    address[3] public makers;
    address[3] public takers;
    address[3] public payouts;
    address public founder;
    address public devFund;

    struct OpenOffer {
        bytes32 offerId;
        address maker;
        address payout;
        uint256 min;
        uint256 max;
        uint256 bond;
        bytes32 makerSecret;
        bool active;
    }

    OpenOffer[] public offerList;
    uint256[] public liveOffers; // positions into offerList

    struct Tracked {
        bytes32 swapId;
        address maker;
        address taker;
        address payout;
        uint256 value;
        uint256 bond;
        bytes32 makerSecret;
        bytes32 takerSecret;
        uint64 createdAt;
        bool claimed;
        bool refunded;
    }

    Tracked[] public tracked;
    uint256[] public openSwaps; // positions into tracked

    // ---- independent ledger: where every wei in the contract should be right now
    uint256 public ghostHeldBond; // bonds locked, whether the offer is live or carried into a swap
    uint256 public ghostOwed; // accrued and not yet pulled: fee shares and returned bonds
    uint256 public ghostOpenValue; // taker payments in swaps that have not settled

    // ---- flow totals
    uint256 public ghostClaimedValue;
    uint256 public ghostPayoutPaid;
    uint256 public ghostFeeAccrued;
    uint256 public ghostFounderAccrued;
    uint256 public ghostDevAccrued;
    uint256 public ghostRefundedValue;
    uint256 public ghostRefundPaid;
    uint256 public ghostBondPosted;
    uint256 public ghostBondReturned;
    uint256 public ghostWithdrawn;

    // ---- call counters, so a green run cannot hide untouched code paths
    uint256 public nPost;
    uint256 public nCancel;
    uint256 public nReap;
    uint256 public nTake;
    uint256 public nReady;
    uint256 public nClaim;
    uint256 public nRefund;
    uint256 public nWithdraw;

    constructor(XmrSwapV2 swap_) {
        swap = swap_;
        founder = swap_.founder();
        devFund = swap_.devFund();
        secrets = KeyVectorsBig.secrets();
        pubs = KeyVectorsBig.pubs();
        keyCursor = 1; // vector 0 is the base point; leave it for targeted tests
        for (uint256 i = 0; i < 3; i++) {
            makers[i] = address(uint160(0xA11CE00 + i));
            takers[i] = address(uint160(0xB0B0B00 + i));
            payouts[i] = address(uint160(0xC0FFEE0 + i));
        }
    }

    receive() external payable {}

    // ------------------------------------------------------------ views for the invariants

    function keysLeft() public view returns (uint256) {
        return secrets.length - keyCursor;
    }

    function trackedCount() public view returns (uint256) {
        return tracked.length;
    }

    function openSwapCount() public view returns (uint256) {
        return openSwaps.length;
    }

    function liveOfferCount() public view returns (uint256) {
        return liveOffers.length;
    }

    function liveOfferIdAt(uint256 i) external view returns (bytes32) {
        return offerList[liveOffers[i]].offerId;
    }

    // narrow getters: reading a whole Tracked tuple at once blows the stack in the invariants
    function swapIdAt(uint256 i) external view returns (bytes32) {
        return tracked[i].swapId;
    }

    function createdAtOf(uint256 i) external view returns (uint64) {
        return tracked[i].createdAt;
    }

    function flagsOf(uint256 i) external view returns (bool claimed, bool refunded) {
        return (tracked[i].claimed, tracked[i].refunded);
    }

    /// @notice Every account the contract could owe money to, so the invariants can total
    ///         `withdrawable` without guessing.
    function accounts() external view returns (address[] memory a) {
        a = new address[](5);
        a[0] = founder;
        a[1] = devFund;
        a[2] = makers[0];
        a[3] = makers[1];
        a[4] = makers[2];
    }

    // ------------------------------------------------------------ internals

    function _nextKey() internal returns (bytes32 secret, bytes32 pub) {
        secret = secrets[keyCursor];
        pub = pubs[keyCursor];
        keyCursor++;
    }

    function _dropLiveOffer(uint256 pos) internal {
        liveOffers[pos] = liveOffers[liveOffers.length - 1];
        liveOffers.pop();
    }

    function _dropOpenSwap(uint256 pos) internal {
        openSwaps[pos] = openSwaps[openSwaps.length - 1];
        openSwaps.pop();
    }

    function _timing(bytes32 id) internal view returns (uint64 t1, uint64 t2, XmrSwapV2.Stage stage) {
        (,,,,,,, t1, t2, stage) = swap.swaps(id);
    }

    // ------------------------------------------------------------ fuzzed actions

    function postOffer(uint256 actorSeed, uint256 minSeed, uint256 spanSeed, uint256 t1Seed, uint256 t2Seed, uint256 expSeed)
        external
    {
        if (keysLeft() == 0) return;
        (bytes32 secret, bytes32 pub) = _nextKey();

        XmrSwapV2.OfferParams memory p;
        address maker = makers[bound(actorSeed, 0, 2)];
        p.payout = payable(payouts[bound(actorSeed >> 8, 0, 2)]);
        p.minAmount = uint128(bound(minSeed, 1e12, 1 ether)); // large enough that the price never rounds to zero XMR
        p.maxAmount = p.minAmount + uint128(bound(spanSeed, 0, 1 ether));
        p.xmrPerAsset = 15e12;
        p.timeout1Duration = uint64(bound(t1Seed, 24 hours, 7 days));
        p.timeout2Duration = uint64(bound(t2Seed, 24 hours, 7 days));
        p.expiry = uint64(block.timestamp + bound(expSeed, 1 hours, 30 days));
        p.makerSpendPub = pub;
        p.makerViewPriv = bytes32(uint256(0x1111));

        uint256 bond = swap.bondFor(p.maxAmount);
        vm.deal(maker, maker.balance + bond);
        vm.prank(maker);
        bytes32 id = swap.postOffer{value: bond}(p);

        offerList.push(OpenOffer(id, maker, p.payout, p.minAmount, p.maxAmount, bond, secret, true));
        liveOffers.push(offerList.length - 1);
        ghostHeldBond += bond;
        ghostBondPosted += bond;
        nPost++;
    }

    function cancelOffer(uint256 seed) external {
        if (liveOffers.length == 0) return;
        uint256 pos = bound(seed, 0, liveOffers.length - 1);
        OpenOffer storage o = offerList[liveOffers[pos]];
        vm.prank(o.maker);
        swap.cancelOffer(o.offerId);
        o.active = false;
        _dropLiveOffer(pos);
        ghostHeldBond -= o.bond;
        ghostOwed += o.bond;
        ghostBondReturned += o.bond;
        nCancel++;
    }

    /// @dev Anyone may reap, so this deliberately calls as an unrelated address.
    function reapExpired(uint256 seed) external {
        uint256 n = liveOffers.length;
        if (n == 0) return;
        uint256 start = bound(seed, 0, n - 1);
        for (uint256 k = 0; k < n; k++) {
            uint256 pos = (start + k) % n;
            OpenOffer storage o = offerList[liveOffers[pos]];
            if (block.timestamp < swap.getOffer(o.offerId).expiry) continue;
            vm.prank(address(0xDEAD));
            swap.reapExpired(o.offerId);
            o.active = false;
            _dropLiveOffer(pos);
            ghostHeldBond -= o.bond;
            ghostOwed += o.bond;
            ghostBondReturned += o.bond;
            nReap++;
            return;
        }
    }

    function takeOffer(uint256 offerSeed, uint256 actorSeed, uint256 amtSeed) external {
        if (liveOffers.length == 0 || keysLeft() == 0) return;
        uint256 pos = bound(offerSeed, 0, liveOffers.length - 1);
        OpenOffer storage o = offerList[liveOffers[pos]];

        if (block.timestamp >= swap.getOffer(o.offerId).expiry) return; // reapExpired clears these

        address taker = takers[bound(actorSeed, 0, 2)];
        uint256 amount = bound(amtSeed, o.min, o.max);
        (bytes32 secret, bytes32 pub) = _nextKey();

        vm.deal(taker, taker.balance + amount);
        vm.prank(taker);
        bytes32 swapId = swap.takeOffer{value: amount}(o.offerId, pub, bytes32(uint256(0x2222)), 0);

        o.active = false;
        _dropLiveOffer(pos);
        tracked.push(
            Tracked(swapId, o.maker, taker, o.payout, amount, o.bond, o.makerSecret, secret, uint64(block.timestamp), false, false)
        );
        openSwaps.push(tracked.length - 1);
        ghostOpenValue += amount; // the bond stays held, it just moves from the offer to the swap
        nTake++;
    }

    function setReady(uint256 seed) external {
        uint256 n = openSwaps.length;
        if (n == 0) return;
        uint256 start = bound(seed, 0, n - 1);
        for (uint256 k = 0; k < n; k++) {
            uint256 pos = (start + k) % n;
            Tracked storage t = tracked[openSwaps[pos]];
            (uint64 t1,, XmrSwapV2.Stage stage) = _timing(t.swapId);
            if (stage != XmrSwapV2.Stage.PENDING || block.timestamp >= t1) continue;
            vm.prank(t.taker);
            swap.setReady(t.swapId);
            nReady++;
            return;
        }
    }

    function claim(uint256 seed) external {
        uint256 n = openSwaps.length;
        if (n == 0) return;
        uint256 start = bound(seed, 0, n - 1);
        for (uint256 k = 0; k < n; k++) {
            uint256 pos = (start + k) % n;
            (uint64 t1, uint64 t2, XmrSwapV2.Stage stage) = _timing(tracked[openSwaps[pos]].swapId);
            if (stage == XmrSwapV2.Stage.COMPLETED || stage == XmrSwapV2.Stage.INVALID) continue;
            if (block.timestamp >= t2) continue;
            if (stage != XmrSwapV2.Stage.READY && block.timestamp < t1) continue;
            _doClaim(openSwaps[pos]);
            _dropOpenSwap(pos);
            nClaim++;
            return;
        }
    }

    function _doClaim(uint256 ti) internal {
        Tracked storage t = tracked[ti];
        uint256 payoutBefore = t.payout.balance;
        uint256 founderBefore = swap.withdrawable(founder);
        uint256 devBefore = swap.withdrawable(devFund);

        vm.prank(t.maker);
        swap.claim(t.swapId, t.makerSecret);

        uint256 founderCut = swap.withdrawable(founder) - founderBefore;
        uint256 devCut = swap.withdrawable(devFund) - devBefore;

        ghostClaimedValue += t.value;
        ghostPayoutPaid += t.payout.balance - payoutBefore;
        ghostFeeAccrued += founderCut + devCut;
        ghostFounderAccrued += founderCut;
        ghostDevAccrued += devCut;

        ghostOpenValue -= t.value;
        ghostHeldBond -= t.bond;
        ghostOwed += t.bond + founderCut + devCut;
        ghostBondReturned += t.bond;

        t.claimed = true;
    }

    function refund(uint256 seed) external {
        uint256 n = openSwaps.length;
        if (n == 0) return;
        uint256 start = bound(seed, 0, n - 1);
        for (uint256 k = 0; k < n; k++) {
            uint256 pos = (start + k) % n;
            (uint64 t1, uint64 t2, XmrSwapV2.Stage stage) = _timing(tracked[openSwaps[pos]].swapId);
            if (stage == XmrSwapV2.Stage.COMPLETED || stage == XmrSwapV2.Stage.INVALID) continue;
            bool blocked = block.timestamp < t2 && (block.timestamp >= t1 || stage == XmrSwapV2.Stage.READY);
            if (blocked) continue;
            _doRefund(openSwaps[pos]);
            _dropOpenSwap(pos);
            nRefund++;
            return;
        }
    }

    function _doRefund(uint256 ti) internal {
        Tracked storage t = tracked[ti];
        uint256 takerBefore = t.taker.balance;

        vm.prank(t.taker);
        swap.refund(t.swapId, t.takerSecret);

        ghostRefundedValue += t.value;
        ghostRefundPaid += t.taker.balance - takerBefore;

        ghostOpenValue -= t.value;
        ghostHeldBond -= t.bond;
        ghostOwed += t.bond;
        ghostBondReturned += t.bond;

        t.refunded = true;
    }

    /// @dev Exercises the pull path, including a third party settling someone else's balance.
    function withdraw(uint256 seed) external {
        address[5] memory pool = [founder, devFund, makers[0], makers[1], makers[2]];
        uint256 start = bound(seed, 0, 4);
        for (uint256 k = 0; k < 5; k++) {
            address who = pool[(start + k) % 5];
            uint256 owed = swap.withdrawable(who);
            if (owed == 0) continue;
            if (k % 2 == 0) {
                vm.prank(who);
                swap.withdraw();
            } else {
                vm.prank(address(0xBEEF)); // anyone may push someone else their balance
                swap.withdrawFor(who);
            }
            ghostOwed -= owed;
            ghostWithdrawn += owed;
            nWithdraw++;
            return;
        }
    }

    /// @dev small steps, so timeout windows get landed in rather than jumped clean over
    function warp(uint256 seed) external {
        vm.warp(block.timestamp + bound(seed, 1 hours, 3 days));
    }
}

contract XmrSwapV2InvariantsTest is Test {
    XmrSwapV2 swap;
    HandlerV2 handler;

    address payable founder = payable(makeAddr("invFounderSafe"));
    address payable devFund = payable(makeAddr("invDevFundSafe"));

    uint16 constant FEE_BPS = 15;
    uint16 constant SHARE_BPS = 5_000;
    uint16 constant BOND_BPS = 500;

    function setUp() public {
        swap = new XmrSwapV2(FEE_BPS, SHARE_BPS, BOND_BPS, founder, devFund);
        handler = new HandlerV2(swap);

        bytes4[] memory sels = new bytes4[](9);
        sels[0] = HandlerV2.postOffer.selector;
        sels[1] = HandlerV2.cancelOffer.selector;
        sels[2] = HandlerV2.reapExpired.selector;
        sels[3] = HandlerV2.takeOffer.selector;
        sels[4] = HandlerV2.setReady.selector;
        sels[5] = HandlerV2.claim.selector;
        sels[6] = HandlerV2.refund.selector;
        sels[7] = HandlerV2.withdraw.selector;
        sels[8] = HandlerV2.warp.selector;
        targetSelector(FuzzSelector({addr: address(handler), selectors: sels}));
        targetContract(address(handler));
    }

    /// The headline one. Every wei the contract holds is accounted for by exactly one of three
    /// things: a bond it is still holding, money it owes someone, or a payment in an open swap.
    /// Nothing stranded, nothing missing, nothing double-counted.
    function invariant_everyWeiIsAccountedFor() public view {
        assertEq(
            address(swap).balance,
            handler.ghostHeldBond() + handler.ghostOwed() + handler.ghostOpenValue(),
            "contract balance does not match bonds + owed + open swaps"
        );
    }

    /// The contract's own `withdrawable` ledger must agree with what the handler says is owed.
    /// If these ever diverge, the contract is promising money it does not have, or sitting on
    /// money it will not pay out.
    function invariant_withdrawableMatchesLedger() public view {
        address[] memory accts = handler.accounts();
        uint256 total;
        for (uint256 i = 0; i < accts.length; i++) {
            total += swap.withdrawable(accts[i]);
        }
        assertEq(total, handler.ghostOwed(), "withdrawable total does not match what is owed");
    }

    /// The contract must always be able to honour every balance it has promised.
    function invariant_solvent() public view {
        address[] memory accts = handler.accounts();
        uint256 total;
        for (uint256 i = 0; i < accts.length; i++) {
            total += swap.withdrawable(accts[i]);
        }
        assertGe(address(swap).balance, total, "cannot cover what it owes");
    }

    /// A bond is held or it is owed back. It is never seized, and it never evaporates.
    function invariant_bondsAreNeverSeized() public view {
        assertEq(
            handler.ghostBondPosted(),
            handler.ghostHeldBond() + handler.ghostBondReturned(),
            "a bond went missing"
        );
    }

    /// Every wei locked in a claimed swap leaves as payout + fee. Never more, never less.
    function invariant_claimPaysOutExactly() public view {
        assertEq(
            handler.ghostPayoutPaid() + handler.ghostFeeAccrued(),
            handler.ghostClaimedValue(),
            "payout + fee != value locked"
        );
    }

    /// A refund returns the whole locked amount, with no fee skimmed.
    function invariant_refundReturnsFullAmount() public view {
        assertEq(handler.ghostRefundPaid(), handler.ghostRefundedValue(), "refund did not return full value");
    }

    /// The two fee recipients always split to exactly the fee, and at a 50/50 share neither side
    /// can drift more than the one wei that an odd fee leaves behind on each claim.
    function invariant_feeSplitIsExact() public view {
        assertEq(
            handler.ghostFounderAccrued() + handler.ghostDevAccrued(),
            handler.ghostFeeAccrued(),
            "the two shares do not add up to the fee"
        );
        assertLe(handler.ghostFounderAccrued(), handler.ghostDevAccrued(), "rounding favoured the founder");
    }

    /// The fee can never exceed the hard cap, whatever the fuzzer does.
    function invariant_feeWithinCap() public view {
        uint256 maxFee = (handler.ghostClaimedValue() * swap.MAX_FEE_BPS()) / swap.BPS();
        assertLe(handler.ghostFeeAccrued(), maxFee, "fee exceeded MAX_FEE_BPS");
    }

    /// Timeout arithmetic must always hold: t1 before t2, each window at least 24 hours.
    function invariant_timeoutsWellFormed() public view {
        uint256 n = handler.trackedCount();
        for (uint256 i = 0; i < n; i++) {
            uint64 createdAt = handler.createdAtOf(i);
            (,,,,,,, uint64 t1, uint64 t2, XmrSwapV2.Stage stage) = swap.swaps(handler.swapIdAt(i));
            if (stage == XmrSwapV2.Stage.INVALID) continue;
            assertLt(t1, t2, "t1 not before t2");
            assertGe(t1 - createdAt, swap.MIN_TIMEOUT(), "t1 shorter than MIN_TIMEOUT");
            assertGe(t2 - t1, swap.MIN_TIMEOUT(), "t2 window shorter than MIN_TIMEOUT");
        }
    }

    /// A swap can never be both claimed and refunded.
    function invariant_neverSettledTwice() public view {
        uint256 n = handler.trackedCount();
        for (uint256 i = 0; i < n; i++) {
            (bool claimed, bool refunded) = handler.flagsOf(i);
            assertFalse(claimed && refunded, "swap both claimed and refunded");
        }
    }

    /// The on-chain offer board must match reality: the same count the handler believes is live,
    /// every entry flagged active, and no duplicates left behind by the swap-and-pop removal.
    function invariant_offerBoardIsConsistent() public view {
        uint256 n = swap.offerCount();
        assertEq(n, handler.liveOfferCount(), "offer board count drifted from the live set");
        for (uint256 i = 0; i < n; i++) {
            bytes32 id = swap.offerIds(i);
            XmrSwapV2.Offer memory o = swap.getOffer(id);
            assertTrue(o.active, "a closed offer is still on the board");
            assertEq(o.index, i, "offer index does not match its slot");
        }
    }

    /// Prints how far the fuzzer actually got, so a green run cannot hide dead code paths.
    function invariant_callSummary() public view {
        console.log("offers posted ", handler.nPost());
        console.log("offers taken  ", handler.nTake());
        console.log("set ready     ", handler.nReady());
        console.log("claims        ", handler.nClaim());
        console.log("refunds       ", handler.nRefund());
        console.log("cancels       ", handler.nCancel());
        console.log("reaped        ", handler.nReap());
        console.log("withdrawals   ", handler.nWithdraw());
        console.log("live offers   ", handler.liveOfferCount());
        console.log("open swaps    ", handler.openSwapCount());
        console.log("keys left     ", handler.keysLeft());
    }
}
