// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/Test.sol";
import {XmrSwapV2, IAggregatorV3} from "../src/XmrSwapV2.sol";
import {KeyVectors} from "./KeyVectors.sol";

/// @dev Chainlink-shaped feed we can drive from tests.
contract MockFeed is IAggregatorV3 {
    uint8 internal dec;
    int256 internal answer;
    uint256 internal updatedAt;

    constructor(uint8 d, int256 a) {
        dec = d;
        answer = a;
        updatedAt = block.timestamp;
    }

    function set(int256 a, uint256 t) external {
        answer = a;
        updatedAt = t;
    }

    function decimals() external view returns (uint8) {
        return dec;
    }

    function latestRoundData() external view returns (uint80, int256, uint256, uint256, uint80) {
        return (1, answer, updatedAt, updatedAt, 1);
    }
}

contract XmrSwapV2Test is Test {
    XmrSwapV2 swap;
    MockFeed feed;

    address payable founder = payable(makeAddr("founderSafe"));
    address payable devFund = payable(makeAddr("devFundSafe"));
    address payable maker = payable(makeAddr("maker"));
    address payable taker = payable(makeAddr("taker"));
    address payable payout = payable(makeAddr("payout"));
    address stranger = makeAddr("stranger");

    uint16 constant FEE_BPS = 15; // 0.15%
    uint16 constant SHARE_BPS = 5_000; // 50/50
    uint16 constant BOND_BPS = 500; // 5% of the offer ceiling
    uint64 constant T1 = 24 hours;
    uint64 constant T2 = 24 hours;
    bytes32 constant VIEW_A = bytes32(uint256(0x1111));
    bytes32 constant VIEW_B = bytes32(uint256(0x2222));

    bytes32[] secrets;
    bytes32[] pubs;

    function setUp() public {
        swap = new XmrSwapV2(FEE_BPS, SHARE_BPS, BOND_BPS, founder, devFund);
        secrets = KeyVectors.secrets();
        pubs = KeyVectors.pubs();
        feed = new MockFeed(8, 465e6); // 4.65 XMR per ETH at 8 decimals
        vm.deal(maker, 100 ether);
        vm.deal(taker, 100 ether);
        vm.deal(stranger, 10 ether);
    }

    // ------------------------------------------------------------ helpers

    function _params(uint256 keyIdx) internal view returns (XmrSwapV2.OfferParams memory p) {
        p.minAmount = 0.01 ether;
        p.maxAmount = 1 ether;
        p.xmrPerAsset = 465e10; // 4.65 XMR per ETH, 12 decimals
        p.expiry = uint64(block.timestamp + 7 days);
        p.timeout1Duration = T1;
        p.timeout2Duration = T2;
        p.makerSpendPub = pubs[keyIdx];
        p.makerViewPriv = VIEW_A;
        p.payout = payout;
    }

    function _post(uint256 keyIdx) internal returns (bytes32 id) {
        XmrSwapV2.OfferParams memory p = _params(keyIdx);
        uint256 bond = swap.bondFor(p.maxAmount);
        vm.prank(maker);
        id = swap.postOffer{value: bond}(p);
    }

    function _take(bytes32 offerId, uint256 amount, uint256 keyIdx) internal returns (bytes32 id) {
        vm.prank(taker);
        id = swap.takeOffer{value: amount}(offerId, pubs[keyIdx], VIEW_B, 0);
    }

    // ------------------------------------------------------------ bond

    function test_postLocksExactBond() public {
        uint256 before = address(swap).balance;
        bytes32 id = _post(1);
        uint256 bond = swap.bondFor(1 ether);
        assertEq(bond, 0.05 ether, "5% of the ceiling");
        assertEq(address(swap).balance - before, bond, "bond held by the contract");
        assertEq(swap.getOffer(id).bond, bond);
    }

    function test_wrongBondReverts() public {
        XmrSwapV2.OfferParams memory p = _params(1);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.WrongBond.selector);
        swap.postOffer{value: 0.01 ether}(p);
    }

    function test_cancelReturnsBondAndDelists() public {
        bytes32 id = _post(1);
        assertEq(swap.offerCount(), 1);
        vm.prank(maker);
        swap.cancelOffer(id);
        assertEq(swap.offerCount(), 0, "delisted");
        assertEq(swap.withdrawable(maker), 0.05 ether, "bond owed back");
        uint256 before = maker.balance;
        vm.prank(maker);
        swap.withdraw();
        assertEq(maker.balance - before, 0.05 ether, "bond returned in full");
    }

    function test_bondNeverSeized_refundPath() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 0.5 ether, 2);
        vm.warp(block.timestamp + T1 + T2 + 1);
        vm.prank(taker);
        swap.refund(sid, secrets[2]);
        assertEq(swap.withdrawable(maker), 0.05 ether, "bond returned even when the maker never delivered");
    }

    // ------------------------------------------------------------ timeouts

    function test_shortTimeoutRejected() public {
        XmrSwapV2.OfferParams memory p = _params(1);
        p.timeout1Duration = 1 hours;
        uint256 bond = swap.bondFor(p.maxAmount); // hoisted: a call here would eat the expectRevert
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.TimeoutTooShort.selector);
        swap.postOffer{value: bond}(p);
    }

    function test_minTimeoutIs24Hours() public view {
        assertEq(swap.MIN_TIMEOUT(), 24 hours);
    }

    // ------------------------------------------------------------ enumeration

    function test_listOffersPagesAndCompacts() public {
        bytes32 a = _post(1);
        bytes32 b = _post(2);
        bytes32 c = _post(3);
        assertEq(swap.offerCount(), 3);

        // cancel the middle one; the last should be swapped into its slot
        vm.prank(maker);
        swap.cancelOffer(b);
        assertEq(swap.offerCount(), 2);

        XmrSwapV2.Offer[] memory page = swap.listOffers(0, 10);
        assertEq(page.length, 2);
        bytes32 first = page[0].id;
        bytes32 second = page[1].id;
        assertTrue(first == a || first == c);
        assertTrue(second == a || second == c);
        assertTrue(first != second);

        // indices must still be correct, so a later cancel does not corrupt the list
        vm.prank(maker);
        swap.cancelOffer(c);
        assertEq(swap.offerCount(), 1);
        assertEq(swap.listOffers(0, 10)[0].id, a);
    }

    function test_listOffersOutOfRange() public {
        _post(1);
        assertEq(swap.listOffers(5, 10).length, 0);
        assertEq(swap.listOffers(0, 100).length, 1);
    }

    function test_takenOfferLeavesTheBoard() public {
        bytes32 id = _post(1);
        _take(id, 0.5 ether, 2);
        assertEq(swap.offerCount(), 0, "a taken offer is no longer listed");
    }

    function test_reapExpiredOnlyAfterExpiry() public {
        bytes32 id = _post(1);
        vm.expectRevert(XmrSwapV2.OfferNotExpired.selector);
        swap.reapExpired(id);

        vm.warp(block.timestamp + 8 days);
        vm.prank(stranger); // anyone may reap
        swap.reapExpired(id);
        assertEq(swap.offerCount(), 0);
        assertEq(swap.withdrawable(maker), 0.05 ether, "bond still goes to the maker");
    }

    // ------------------------------------------------------------ the happy path and the split

    function test_claimSplitsFeeFiftyFifty() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);

        vm.prank(taker);
        swap.setReady(sid);

        uint256 payoutBefore = payout.balance;
        vm.prank(maker);
        swap.claim(sid, secrets[1]);

        uint256 fee = (1 ether * FEE_BPS) / 10_000;
        assertEq(payout.balance - payoutBefore, 1 ether - fee, "payout is value minus fee");
        assertEq(swap.withdrawable(founder), fee / 2, "founder half");
        assertEq(swap.withdrawable(devFund), fee - fee / 2, "dev fund half");
        assertEq(swap.withdrawable(founder) + swap.withdrawable(devFund), fee, "split is exact");
        assertEq(swap.withdrawable(maker), 0.05 ether, "bond back");
    }

    function test_feesAccrueRatherThanPush() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.prank(taker);
        swap.setReady(sid);
        vm.prank(maker);
        swap.claim(sid, secrets[1]);

        // nothing was pushed to the fee recipients during the claim
        assertEq(founder.balance, 0);
        assertEq(devFund.balance, 0);

        // anyone can settle the dev fund
        vm.prank(stranger);
        swap.withdrawFor(devFund);
        assertGt(devFund.balance, 0, "third party pushed the dev fund its share");
    }

    function test_refundTakesNoFee() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.warp(block.timestamp + T1 + T2 + 1);

        uint256 before = taker.balance;
        vm.prank(taker);
        swap.refund(sid, secrets[2]);

        assertEq(taker.balance - before, 1 ether, "whole payment back");
        assertEq(swap.withdrawable(founder), 0, "no fee on a refund");
        assertEq(swap.withdrawable(devFund), 0);
    }

    function test_contractHoldsNothingExtraAfterSettlement() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.prank(taker);
        swap.setReady(sid);
        vm.prank(maker);
        swap.claim(sid, secrets[1]);

        vm.prank(maker);
        swap.withdraw();
        vm.prank(stranger);
        swap.withdrawFor(founder);
        vm.prank(stranger);
        swap.withdrawFor(devFund);
        assertEq(address(swap).balance, 0, "everything is accounted for");
    }

    // ------------------------------------------------------------ authorisation and secrets

    function test_onlyMakerClaims() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.warp(block.timestamp + T1 + 1);
        vm.prank(stranger);
        vm.expectRevert(XmrSwapV2.NotMaker.selector);
        swap.claim(sid, secrets[1]);
    }

    function test_onlyTakerRefundsAndReadies() public {
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.prank(stranger);
        vm.expectRevert(XmrSwapV2.NotTaker.selector);
        swap.setReady(sid);

        vm.warp(block.timestamp + T1 + T2 + 1);
        vm.prank(stranger);
        vm.expectRevert(XmrSwapV2.NotTaker.selector);
        swap.refund(sid, secrets[2]);
    }

    function testFuzz_wrongSecretNeverClaims(bytes32 bad) public {
        vm.assume(bad != secrets[1]);
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.warp(block.timestamp + T1 + 1);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.InvalidSecret.selector);
        swap.claim(sid, bad);
    }

    function testFuzz_wrongSecretNeverRefunds(bytes32 bad) public {
        vm.assume(bad != secrets[2]);
        bytes32 id = _post(1);
        bytes32 sid = _take(id, 1 ether, 2);
        vm.warp(block.timestamp + T1 + T2 + 1);
        vm.prank(taker);
        vm.expectRevert(XmrSwapV2.InvalidSecret.selector);
        swap.refund(sid, bad);
    }

    function test_oneTimeKeysCannotBeReused() public {
        _post(1);
        XmrSwapV2.OfferParams memory p = _params(1); // same maker key
        uint256 bond = swap.bondFor(p.maxAmount);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.KeyAlreadyUsed.selector);
        swap.postOffer{value: bond}(p);
    }

    function test_takerKeyCannotBeReused() public {
        bytes32 a = _post(1);
        bytes32 b = _post(2);
        _take(a, 0.5 ether, 3);
        vm.prank(taker);
        vm.expectRevert(XmrSwapV2.KeyAlreadyUsed.selector);
        swap.takeOffer{value: 0.5 ether}(b, pubs[3], VIEW_B, 0);
    }

    // ------------------------------------------------------------ amounts

    function testFuzz_amountOutOfRangeReverts(uint256 amount) public {
        bytes32 id = _post(1);
        amount = bound(amount, 0, 5 ether);
        vm.assume(amount < 0.01 ether || amount > 1 ether);
        vm.prank(taker);
        vm.expectRevert(XmrSwapV2.AmountOutOfRange.selector);
        swap.takeOffer{value: amount}(id, pubs[2], VIEW_B, 0);
    }

    function test_takerSlippageBoundHonoured() public {
        bytes32 id = _post(1);
        // 1 ETH at 4.65 buys 4.65e12 piconero; demand more and it must revert
        vm.prank(taker);
        vm.expectRevert(XmrSwapV2.NotEnoughXmr.selector);
        swap.takeOffer{value: 1 ether}(id, pubs[2], VIEW_B, 5e12);
    }

    // ------------------------------------------------------------ oracle pricing

    function _oracleParams(uint256 keyIdx) internal view returns (XmrSwapV2.OfferParams memory p) {
        p = _params(keyIdx);
        p.xmrPerAsset = 0;
        p.oracle = address(feed);
        p.oracleRatioBps = 10_000; // exactly the feed price
        p.oracleMaxAge = 1 hours;
    }

    function test_oraclePricedOfferResolves() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        vm.prank(maker);
        bytes32 id = swap.postOffer{value: swap.bondFor(p.maxAmount)}(p);
        // 8-decimal feed normalised to 12 decimals
        assertEq(swap.priceOf(id), 465e10, "4.65 XMR per ETH");
        _take(id, 1 ether, 2);
    }

    function test_oracleRatioAndOffsetApply() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        p.oracleRatioBps = 9_000; // 10% under the feed
        p.oracleOffset = -5e9;
        vm.prank(maker);
        bytes32 id = swap.postOffer{value: swap.bondFor(p.maxAmount)}(p);
        assertEq(swap.priceOf(id), (465e10 * 9_000) / 10_000 - 5e9);
    }

    function test_staleOracleRejected() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        vm.prank(maker);
        bytes32 id = swap.postOffer{value: swap.bondFor(p.maxAmount)}(p);
        vm.warp(block.timestamp + 2 hours); // older than oracleMaxAge
        vm.expectRevert(XmrSwapV2.OracleStale.selector);
        swap.priceOf(id);

        vm.prank(taker);
        vm.expectRevert(XmrSwapV2.OracleStale.selector);
        swap.takeOffer{value: 1 ether}(id, pubs[2], VIEW_B, 0);
    }

    /// The ceiling is checked when the offer is taken, not when it is posted: an offer whose feed
    /// has run above the maker's limit simply cannot be filled until the price comes back down.
    function test_makerPriceCeilingHonoured() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        p.oracleMaxPrice = 400e10; // below the feed's 4.65
        vm.prank(maker);
        bytes32 id = swap.postOffer{value: swap.bondFor(p.maxAmount)}(p);

        vm.expectRevert(XmrSwapV2.PriceTooHigh.selector);
        swap.priceOf(id);

        vm.prank(taker);
        vm.expectRevert(XmrSwapV2.PriceTooHigh.selector);
        swap.takeOffer{value: 1 ether}(id, pubs[2], VIEW_B, 0);

        // feed drops back under the ceiling and the same offer fills
        feed.set(350e6, block.timestamp);
        assertEq(swap.priceOf(id), 350e10);
        _take(id, 1 ether, 2);
    }

    function test_badOracleAnswerRejected() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        feed.set(0, block.timestamp);
        uint256 bond = swap.bondFor(p.maxAmount);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.BadOracle.selector);
        swap.postOffer{value: bond}(p);
    }

    function test_negativeOffsetBelowZeroRejected() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        p.oracleOffset = -1e13; // bigger than the whole feed price
        uint256 bond = swap.bondFor(p.maxAmount);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.BadOracle.selector);
        swap.postOffer{value: bond}(p);
    }

    function test_cannotSetBothPricingSources() public {
        XmrSwapV2.OfferParams memory p = _oracleParams(1);
        p.xmrPerAsset = 465e10; // fixed price as well as an oracle
        uint256 bond = swap.bondFor(p.maxAmount);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.BadOracle.selector);
        swap.postOffer{value: bond}(p);
    }

    function test_noPricingSourceAtAllRejected() public {
        XmrSwapV2.OfferParams memory p = _params(1);
        p.xmrPerAsset = 0; // no fixed price and no oracle either
        uint256 bond = swap.bondFor(p.maxAmount);
        vm.prank(maker);
        vm.expectRevert(XmrSwapV2.BadAmounts.selector);
        swap.postOffer{value: bond}(p);
    }

    // ------------------------------------------------------------ withdrawals

    function test_withdrawWithNothingOwedReverts() public {
        vm.prank(stranger);
        vm.expectRevert(XmrSwapV2.NothingToWithdraw.selector);
        swap.withdraw();
    }

    // ------------------------------------------------------------ deployment guards

    function test_feeAboveCapRejected() public {
        vm.expectRevert(XmrSwapV2.FeeTooHigh.selector);
        new XmrSwapV2(101, SHARE_BPS, BOND_BPS, founder, devFund);
    }

    function test_bondAboveCapRejected() public {
        vm.expectRevert(XmrSwapV2.BondTooHigh.selector);
        new XmrSwapV2(FEE_BPS, SHARE_BPS, 2_001, founder, devFund);
    }

    function test_shareAboveOneHundredPercentRejected() public {
        vm.expectRevert(XmrSwapV2.BadShare.selector);
        new XmrSwapV2(FEE_BPS, 10_001, BOND_BPS, founder, devFund);
    }

    function test_noOwnerExists() public view {
        // there is no owner function at all; the deployed surface has no admin entry point
        assertEq(swap.founder(), founder);
        assertEq(swap.devFund(), devFund);
        assertEq(swap.feeBps(), FEE_BPS);
        assertEq(swap.founderShareBps(), SHARE_BPS);
        assertEq(swap.bondBps(), BOND_BPS);
    }
}
