// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test, console} from "forge-std/Test.sol";
import {XmrSwap} from "../src/XmrSwap.sol";
import {KeyVectorsBig} from "./KeyVectorsBig.sol";

/// @dev Drives XmrSwap through random call sequences on real ETH swaps. One-time keys come from a
///      513-vector pool so the fuzzer can build hundreds of genuine swaps. Active offers and open
///      swaps are tracked in their own index lists, so a fuzzed call almost always does real work
///      instead of picking a dead item and bailing. Ghost totals are kept independently of the
///      contract, so the invariants check the money without trusting its own accounting.
contract Handler is Test {
    XmrSwap public swap;

    bytes32[] internal secrets;
    bytes32[] internal pubs;
    uint256 public keyCursor;

    address[3] public makers;
    address[3] public takers;
    address[3] public payouts;

    struct OpenOffer {
        bytes32 offerId;
        address maker;
        address payout;
        uint256 min;
        uint256 max;
        bytes32 makerSecret;
        bool active;
    }
    OpenOffer[] public offerList;
    uint256[] public activeOffers; // positions into offerList

    struct Tracked {
        bytes32 swapId;
        address maker;
        address taker;
        address payout;
        uint256 value;
        bytes32 makerSecret;
        bytes32 takerSecret;
        uint64 createdAt;
        bool settled;
        bool claimed;
        bool refunded;
    }
    Tracked[] public tracked;
    uint256[] public openSwaps; // positions into tracked

    // ghost accounting, maintained independently of the contract
    uint256 public ghostClaimedValue;
    uint256 public ghostPayoutPaid;
    uint256 public ghostFeePaid;
    uint256 public ghostRefundedValue;
    uint256 public ghostRefundPaid;

    // call counters, so a green run cannot hide untouched code paths
    uint256 public nPost;
    uint256 public nCancel;
    uint256 public nTake;
    uint256 public nReady;
    uint256 public nClaim;
    uint256 public nRefund;

    constructor(XmrSwap swap_) {
        swap = swap_;
        secrets = KeyVectorsBig.secrets();
        pubs = KeyVectorsBig.pubs();
        // skip vector 0 (scalar 1, the base point) so it stays free for targeted tests
        keyCursor = 1;
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

    // ------------------------------------------------------------ internals

    function _nextKey() internal returns (bytes32 secret, bytes32 pub) {
        secret = secrets[keyCursor];
        pub = pubs[keyCursor];
        keyCursor++;
    }

    function _dropActiveOffer(uint256 pos) internal {
        activeOffers[pos] = activeOffers[activeOffers.length - 1];
        activeOffers.pop();
    }

    function _dropOpenSwap(uint256 pos) internal {
        openSwaps[pos] = openSwaps[openSwaps.length - 1];
        openSwaps.pop();
    }

    function _offerExpiry(bytes32 id) internal view returns (uint64 e) {
        (, , , , , , e, , , , , ) = swap.offers(id);
    }

    function _swapTiming(bytes32 id) internal view returns (uint64 t1, uint64 t2, XmrSwap.Stage stage) {
        (, , , , , , , t1, t2, stage) = swap.swaps(id);
    }

    // ------------------------------------------------------------ fuzzed actions

    struct PostParams {
        address maker;
        address payout;
        uint256 minA;
        uint256 maxA;
        uint64 t1;
        uint64 t2;
        uint64 expiry;
        bytes32 pub;
    }

    function postOffer(uint256 actorSeed, uint256 minSeed, uint256 spanSeed, uint256 t1Seed, uint256 t2Seed, uint256 expSeed)
        external
    {
        if (keysLeft() == 0) return;
        (bytes32 secret, bytes32 pub) = _nextKey();

        PostParams memory p;
        p.maker = makers[bound(actorSeed, 0, 2)];
        p.payout = payouts[bound(actorSeed >> 8, 0, 2)];
        p.minA = bound(minSeed, 1, 1 ether);
        p.maxA = p.minA + bound(spanSeed, 0, 1 ether);
        p.t1 = uint64(bound(t1Seed, 1 hours, 7 days));
        p.t2 = uint64(bound(t2Seed, 1 hours, 7 days));
        p.expiry = uint64(block.timestamp + bound(expSeed, 1 hours, 30 days));
        p.pub = pub;

        offerList.push(OpenOffer(_post(p), p.maker, p.payout, p.minA, p.maxA, secret, true));
        activeOffers.push(offerList.length - 1);
        nPost++;
    }

    /// @dev params travel in memory so the ten-argument call does not blow the stack
    function _post(PostParams memory p) internal returns (bytes32) {
        vm.prank(p.maker);
        return swap.postOffer(
            address(0),
            uint128(p.minA),
            uint128(p.maxA),
            15e12,
            p.expiry,
            p.t1,
            p.t2,
            p.pub,
            bytes32(uint256(0x1111)),
            payable(p.payout)
        );
    }

    function cancelOffer(uint256 seed) external {
        if (activeOffers.length == 0) return;
        uint256 pos = bound(seed, 0, activeOffers.length - 1);
        OpenOffer storage o = offerList[activeOffers[pos]];
        vm.prank(o.maker);
        swap.cancelOffer(o.offerId);
        o.active = false;
        _dropActiveOffer(pos);
        nCancel++;
    }

    function takeOffer(uint256 offerSeed, uint256 actorSeed, uint256 amtSeed) external {
        if (activeOffers.length == 0 || keysLeft() == 0) return;
        uint256 pos = bound(offerSeed, 0, activeOffers.length - 1);
        OpenOffer storage o = offerList[activeOffers[pos]];

        if (block.timestamp >= _offerExpiry(o.offerId)) {
            // an expired offer can never be taken; retire it so calls are not wasted on it
            o.active = false;
            _dropActiveOffer(pos);
            return;
        }

        address taker = takers[bound(actorSeed, 0, 2)];
        uint256 amount = bound(amtSeed, o.min, o.max);
        (bytes32 secret, bytes32 pub) = _nextKey();

        vm.deal(taker, amount);
        vm.prank(taker);
        bytes32 swapId = swap.takeOffer{value: amount}(o.offerId, amount, pub, bytes32(uint256(0x2222)));

        o.active = false;
        _dropActiveOffer(pos);
        tracked.push(
            Tracked(swapId, o.maker, taker, o.payout, amount, o.makerSecret, secret, uint64(block.timestamp), false, false, false)
        );
        openSwaps.push(tracked.length - 1);
        nTake++;
    }

    function setReady(uint256 seed) external {
        uint256 n = openSwaps.length;
        if (n == 0) return;
        uint256 start = bound(seed, 0, n - 1);
        for (uint256 k = 0; k < n; k++) {
            uint256 pos = (start + k) % n;
            Tracked storage t = tracked[openSwaps[pos]];
            (uint64 t1, , XmrSwap.Stage stage) = _swapTiming(t.swapId);
            if (stage != XmrSwap.Stage.PENDING || block.timestamp >= t1) continue;
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
            (uint64 t1, uint64 t2, XmrSwap.Stage stage) = _swapTiming(tracked[openSwaps[pos]].swapId);
            if (stage == XmrSwap.Stage.COMPLETED || stage == XmrSwap.Stage.INVALID) continue;
            if (block.timestamp >= t2) continue;
            if (stage != XmrSwap.Stage.READY && block.timestamp < t1) continue;
            _doClaim(openSwaps[pos]);
            _dropOpenSwap(pos);
            nClaim++;
            return;
        }
    }

    function _doClaim(uint256 ti) internal {
        Tracked storage t = tracked[ti];
        uint256 payoutBefore = t.payout.balance;
        uint256 feeBefore = swap.feeRecipient().balance;
        vm.prank(t.maker);
        swap.claim(t.swapId, t.makerSecret);
        ghostClaimedValue += t.value;
        ghostPayoutPaid += t.payout.balance - payoutBefore;
        ghostFeePaid += swap.feeRecipient().balance - feeBefore;
        t.settled = true;
        t.claimed = true;
    }

    function refund(uint256 seed) external {
        uint256 n = openSwaps.length;
        if (n == 0) return;
        uint256 start = bound(seed, 0, n - 1);
        for (uint256 k = 0; k < n; k++) {
            uint256 pos = (start + k) % n;
            (uint64 t1, uint64 t2, XmrSwap.Stage stage) = _swapTiming(tracked[openSwaps[pos]].swapId);
            if (stage == XmrSwap.Stage.COMPLETED || stage == XmrSwap.Stage.INVALID) continue;
            bool blocked = block.timestamp < t2 && (block.timestamp >= t1 || stage == XmrSwap.Stage.READY);
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
        t.settled = true;
        t.refunded = true;
    }

    /// @dev small steps, so timeout windows get landed in rather than jumped clean over
    function warp(uint256 seed) external {
        vm.warp(block.timestamp + bound(seed, 1 minutes, 3 days));
    }
}

contract XmrSwapInvariantsTest is Test {
    XmrSwap swap;
    Handler handler;
    address payable feeWallet = payable(makeAddr("invFeeWallet"));

    uint16 constant FEE_BPS = 15;

    function setUp() public {
        swap = new XmrSwap(FEE_BPS, feeWallet);
        handler = new Handler(swap);
        // without this the fuzzer would also call forge-std's own public functions
        bytes4[] memory sels = new bytes4[](7);
        sels[0] = Handler.postOffer.selector;
        sels[1] = Handler.cancelOffer.selector;
        sels[2] = Handler.takeOffer.selector;
        sels[3] = Handler.setReady.selector;
        sels[4] = Handler.claim.selector;
        sels[5] = Handler.refund.selector;
        sels[6] = Handler.warp.selector;
        targetSelector(FuzzSelector({addr: address(handler), selectors: sels}));
        targetContract(address(handler));
    }

    /// Money in the contract must equal exactly the sum of swaps still open. Nothing stranded,
    /// nothing missing, no open swap left unbacked.
    function invariant_ethBacksEveryOpenSwap() public view {
        uint256 open;
        uint256 n = handler.trackedCount();
        for (uint256 i = 0; i < n; i++) {
            (, , , , uint256 value, , , , , XmrSwap.Stage stage) = swap.swaps(handler.swapIdAt(i));
            if (stage == XmrSwap.Stage.PENDING || stage == XmrSwap.Stage.READY) open += value;
        }
        assertEq(address(swap).balance, open, "contract ETH does not match open swaps");
    }

    /// Every wei locked in a claimed swap comes back out as payout + fee. Never more, never less.
    function invariant_claimPaysOutExactly() public view {
        assertEq(
            handler.ghostPayoutPaid() + handler.ghostFeePaid(),
            handler.ghostClaimedValue(),
            "payout + fee != value locked"
        );
    }

    /// A refund returns the whole locked amount, with no fee skimmed.
    function invariant_refundReturnsFullAmount() public view {
        assertEq(handler.ghostRefundPaid(), handler.ghostRefundedValue(), "refund did not return full value");
    }

    /// The fee can never exceed the hard cap, whatever the fuzzer does.
    function invariant_feeWithinCap() public view {
        uint256 maxFee = (handler.ghostClaimedValue() * swap.MAX_FEE_BPS()) / swap.BPS();
        assertLe(handler.ghostFeePaid(), maxFee, "fee exceeded MAX_FEE_BPS");
    }

    /// Timeout arithmetic must always hold: t1 before t2, each window at least the minimum.
    function invariant_timeoutsWellFormed() public view {
        uint256 n = handler.trackedCount();
        for (uint256 i = 0; i < n; i++) {
            uint64 createdAt = handler.createdAtOf(i);
            (, , , , , , , uint64 t1, uint64 t2, XmrSwap.Stage stage) = swap.swaps(handler.swapIdAt(i));
            if (stage == XmrSwap.Stage.INVALID) continue;
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

    /// Prints how far the fuzzer actually got, so a green run cannot hide dead code paths.
    function invariant_callSummary() public view {
        console.log("offers posted ", handler.nPost());
        console.log("offers taken  ", handler.nTake());
        console.log("set ready     ", handler.nReady());
        console.log("claims        ", handler.nClaim());
        console.log("refunds       ", handler.nRefund());
        console.log("cancels       ", handler.nCancel());
        console.log("open swaps    ", handler.openSwapCount());
        console.log("keys left     ", handler.keysLeft());
    }
}
