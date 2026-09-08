// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {XmrSwap} from "../src/XmrSwap.sol";
import {KeyVectors} from "./KeyVectors.sol";

contract MockToken is ERC20 {
    constructor() ERC20("Mock USD", "MUSD") {
        _mint(msg.sender, 1_000_000e6);
    }
    function decimals() public pure override returns (uint8) { return 6; }
}

contract XmrSwapTest is Test {
    XmrSwap swap;
    MockToken token;

    address payable feeWallet = payable(makeAddr("feeWallet"));
    address payable maker = payable(makeAddr("maker"));
    address payable taker = payable(makeAddr("taker"));
    address payable makerPayout = payable(makeAddr("makerPayout"));
    address stranger = makeAddr("stranger");

    uint16 constant FEE_BPS = 15;
    uint64 constant T1 = 1 hours;
    uint64 constant T2 = 1 hours;

    bytes32[] secrets;
    bytes32[] pubs;

    // maker uses vector 1, taker uses vector 2
    bytes32 makerSecret; bytes32 makerPub;
    bytes32 takerSecret; bytes32 takerPub;
    bytes32 constant MAKER_VIEW = bytes32(uint256(0x1111));
    bytes32 constant TAKER_VIEW = bytes32(uint256(0x2222));

    function setUp() public {
        swap = new XmrSwap(FEE_BPS, feeWallet);
        token = new MockToken();
        secrets = KeyVectors.secrets();
        pubs = KeyVectors.pubs();
        makerSecret = secrets[1]; makerPub = pubs[1];
        takerSecret = secrets[2]; takerPub = pubs[2];
        vm.deal(taker, 100 ether);
        vm.deal(maker, 1 ether);
        token.transfer(taker, 100_000e6);
    }

    // ------------------------------------------------------------ helpers

    function _postEthOffer() internal returns (bytes32 offerId) {
        vm.prank(maker);
        offerId = swap.postOffer(
            address(0), 0.01 ether, 1 ether, 15e12 /* 15 XMR per ETH */, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, makerPayout
        );
    }

    function _postTokenOffer() internal returns (bytes32 offerId) {
        vm.prank(maker);
        offerId = swap.postOffer(
            address(token), 10e6, 1000e6, 6_666_666_666_666_666_666_666 /* piconero per 1e18 base units */, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, makerPayout
        );
    }

    function _takeEth(bytes32 offerId, uint256 amount) internal returns (bytes32 swapId) {
        vm.prank(taker);
        swapId = swap.takeOffer{value: amount}(offerId, amount, takerPub, TAKER_VIEW);
    }

    function _stage(bytes32 swapId) internal view returns (XmrSwap.Stage st) {
        (, , , , , , , , , st) = swap.swaps(swapId);
    }

    // ------------------------------------------------------------ constructor

    function test_constructor_rejectsFeeOverCap() public {
        vm.expectRevert(XmrSwap.FeeTooHigh.selector);
        new XmrSwap(101, feeWallet);
    }

    function test_constructor_rejectsZeroRecipientWithFee() public {
        vm.expectRevert(XmrSwap.ZeroAddress.selector);
        new XmrSwap(15, payable(address(0)));
    }

    function test_constructor_allowsZeroFee() public {
        XmrSwap free = new XmrSwap(0, payable(address(0)));
        assertEq(free.feeBps(), 0);
    }

    // ------------------------------------------------------------ offers

    function test_postOffer_storesAndEmits() public {
        vm.prank(maker);
        vm.expectEmit(false, true, true, true);
        emit XmrSwap.OfferPosted(bytes32(0), maker, address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days));
        bytes32 id = swap.postOffer(address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, makerPayout);
        (address m, address po, address a, uint128 mn, uint128 mx, uint256 px, , , , bytes32 sp, bytes32 vp, bool active) = swap.offers(id);
        assertEq(po, makerPayout);
        assertEq(m, maker); assertEq(a, address(0)); assertEq(mn, 0.01 ether); assertEq(mx, 1 ether);
        assertEq(px, 15e12); assertEq(sp, makerPub); assertEq(vp, MAKER_VIEW); assertTrue(active);
    }

    function test_postOffer_validation() public {
        uint64 exp = uint64(block.timestamp + 1 days);
        vm.startPrank(maker);
        vm.expectRevert(XmrSwap.ZeroKey.selector);
        swap.postOffer(address(0), 1, 2, 1, exp, T1, T2, bytes32(0), MAKER_VIEW, makerPayout);
        vm.expectRevert(XmrSwap.ZeroKey.selector);
        swap.postOffer(address(0), 1, 2, 1, exp, T1, T2, makerPub, bytes32(0), makerPayout);
        vm.expectRevert(XmrSwap.BadAmounts.selector);
        swap.postOffer(address(0), 0, 2, 1, exp, T1, T2, makerPub, MAKER_VIEW, makerPayout);
        vm.expectRevert(XmrSwap.BadAmounts.selector);
        swap.postOffer(address(0), 3, 2, 1, exp, T1, T2, makerPub, MAKER_VIEW, makerPayout);
        vm.expectRevert(XmrSwap.BadAmounts.selector);
        swap.postOffer(address(0), 1, 2, 0, exp, T1, T2, makerPub, MAKER_VIEW, makerPayout);
        vm.expectRevert(XmrSwap.TimeoutTooShort.selector);
        swap.postOffer(address(0), 1, 2, 1, exp, 59 minutes, T2, makerPub, MAKER_VIEW, makerPayout);
        vm.expectRevert(XmrSwap.TimeoutTooShort.selector);
        swap.postOffer(address(0), 1, 2, 1, exp, T1, 59 minutes, makerPub, MAKER_VIEW, makerPayout);
        vm.expectRevert(XmrSwap.OfferExpired.selector);
        swap.postOffer(address(0), 1, 2, 1, uint64(block.timestamp), T1, T2, makerPub, MAKER_VIEW, makerPayout);
        vm.stopPrank();
    }

    function test_postOffer_rejectsZeroPayout() public {
        vm.prank(maker);
        vm.expectRevert(XmrSwap.ZeroAddress.selector);
        swap.postOffer(address(0), 1, 2, 1, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, payable(address(0)));
    }

    function test_offerIds_unique() public {
        bytes32 a = _postEthOffer();
        vm.prank(maker);
        bytes32 b = swap.postOffer(address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days), T1, T2, pubs[4], MAKER_VIEW, makerPayout);
        assertTrue(a != b);
    }

    function test_postOffer_rejectsReusedMakerKey() public {
        _postEthOffer();
        vm.prank(maker);
        vm.expectRevert(XmrSwap.KeyAlreadyUsed.selector);
        swap.postOffer(address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, makerPayout);
        // a different maker cannot reuse it either
        vm.prank(stranger);
        vm.expectRevert(XmrSwap.KeyAlreadyUsed.selector);
        swap.postOffer(address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, makerPayout);
    }

    function test_takeOffer_rejectsReusedTakerKey() public {
        bytes32 id1 = _postEthOffer();
        _takeEth(id1, 0.5 ether);
        vm.prank(maker);
        bytes32 id2 = swap.postOffer(address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days), T1, T2, pubs[4], MAKER_VIEW, makerPayout);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.KeyAlreadyUsed.selector);
        swap.takeOffer{value: 0.5 ether}(id2, 0.5 ether, takerPub, TAKER_VIEW);
        // a taker key may not equal any maker key either
        vm.prank(taker);
        vm.expectRevert(XmrSwap.KeyAlreadyUsed.selector);
        swap.takeOffer{value: 0.5 ether}(id2, 0.5 ether, makerPub, TAKER_VIEW);
    }

    function test_cancelOffer() public {
        bytes32 id = _postEthOffer();
        vm.prank(stranger);
        vm.expectRevert(XmrSwap.NotMaker.selector);
        swap.cancelOffer(id);
        vm.prank(maker);
        swap.cancelOffer(id);
        (, , , , , , , , , , , bool active) = swap.offers(id);
        assertFalse(active);
        vm.prank(maker);
        vm.expectRevert(XmrSwap.OfferNotActive.selector);
        swap.cancelOffer(id);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.OfferNotActive.selector);
        swap.takeOffer{value: 0.5 ether}(id, 0.5 ether, takerPub, TAKER_VIEW);
    }

    // ------------------------------------------------------------ take

    function test_takeOffer_eth_locksAndEmits() public {
        bytes32 id = _postEthOffer();
        uint256 t0 = block.timestamp;
        vm.prank(taker);
        vm.expectEmit(false, true, true, true);
        emit XmrSwap.SwapCreated(bytes32(0), id, taker, 0.5 ether, 7.5e12, takerPub, TAKER_VIEW, uint64(t0 + T1), uint64(t0 + T1 + T2));
        bytes32 sid = swap.takeOffer{value: 0.5 ether}(id, 0.5 ether, takerPub, TAKER_VIEW);
        assertEq(address(swap).balance, 0.5 ether);
        (address tk, address mk, address po, address a, uint256 v, bytes32 mp, bytes32 tp, uint64 t1, uint64 t2, XmrSwap.Stage st) = swap.swaps(sid);
        assertEq(po, makerPayout);
        assertEq(tk, taker); assertEq(mk, maker); assertEq(a, address(0)); assertEq(v, 0.5 ether);
        assertEq(mp, makerPub); assertEq(tp, takerPub); assertEq(t1, t0 + T1); assertEq(t2, t0 + T1 + T2);
        assertEq(uint8(st), uint8(XmrSwap.Stage.PENDING));
        (, , , , , , , , , , , bool active) = swap.offers(id);
        assertFalse(active, "offer is single-use");
    }

    function test_takeOffer_eth_validation() public {
        bytes32 id = _postEthOffer();
        vm.startPrank(taker);
        vm.expectRevert(XmrSwap.WrongValue.selector);
        swap.takeOffer{value: 0.4 ether}(id, 0.5 ether, takerPub, TAKER_VIEW);
        vm.expectRevert(XmrSwap.AmountOutOfRange.selector);
        swap.takeOffer{value: 0.001 ether}(id, 0.001 ether, takerPub, TAKER_VIEW);
        vm.expectRevert(XmrSwap.AmountOutOfRange.selector);
        swap.takeOffer{value: 2 ether}(id, 2 ether, takerPub, TAKER_VIEW);
        vm.expectRevert(XmrSwap.ZeroKey.selector);
        swap.takeOffer{value: 0.5 ether}(id, 0.5 ether, bytes32(0), TAKER_VIEW);
        vm.stopPrank();
        vm.warp(block.timestamp + 1 days);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.OfferExpired.selector);
        swap.takeOffer{value: 0.5 ether}(id, 0.5 ether, takerPub, TAKER_VIEW);
    }

    function test_takeOffer_secondTakeFails() public {
        bytes32 id = _postEthOffer();
        _takeEth(id, 0.5 ether);
        vm.prank(stranger);
        vm.deal(stranger, 1 ether);
        vm.expectRevert(XmrSwap.OfferNotActive.selector);
        swap.takeOffer{value: 0.5 ether}(id, 0.5 ether, pubs[3], TAKER_VIEW);
    }

    function test_takeOffer_token() public {
        bytes32 id = _postTokenOffer();
        vm.startPrank(taker);
        token.approve(address(swap), 500e6);
        vm.expectRevert(XmrSwap.WrongValue.selector);
        swap.takeOffer{value: 1}(id, 500e6, takerPub, TAKER_VIEW);
        bytes32 sid = swap.takeOffer(id, 500e6, takerPub, TAKER_VIEW);
        vm.stopPrank();
        assertEq(token.balanceOf(address(swap)), 500e6);
        (, , , address a, uint256 v, , , , , ) = swap.swaps(sid);
        assertEq(a, address(token)); assertEq(v, 500e6);
    }

    // ------------------------------------------------------------ ready / claim

    function test_setReady_onlyTakerBeforeT1() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(stranger);
        vm.expectRevert(XmrSwap.NotTaker.selector);
        swap.setReady(sid);
        vm.prank(taker);
        swap.setReady(sid);
        assertEq(uint8(_stage(sid)), uint8(XmrSwap.Stage.READY));
        vm.prank(taker);
        vm.expectRevert(XmrSwap.AlreadyReady.selector);
        swap.setReady(sid);
    }

    function test_setReady_rejectedAfterT1() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.warp(block.timestamp + T1);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.TooLateToClaim.selector);
        swap.setReady(sid);
    }

    function test_claim_tooEarlyWithoutReady() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(maker);
        vm.expectRevert(XmrSwap.TooEarlyToClaim.selector);
        swap.claim(sid, makerSecret);
    }

    function test_claim_afterReady_paysMakerAndFee() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(taker); swap.setReady(sid);
        uint256 fee = 0.5 ether * uint256(FEE_BPS) / 10_000; // 0.00075 ETH
        uint256 makerBefore = makerPayout.balance; uint256 feeBefore = feeWallet.balance;
        vm.prank(maker);
        vm.expectEmit(true, false, false, true);
        emit XmrSwap.Claimed(sid, makerSecret, 0.5 ether - fee, fee);
        swap.claim(sid, makerSecret);
        assertEq(makerPayout.balance - makerBefore, 0.5 ether - fee);
        assertEq(feeWallet.balance - feeBefore, fee);
        assertEq(fee, 0.00075 ether);
        assertEq(address(swap).balance, 0);
        assertEq(uint8(_stage(sid)), uint8(XmrSwap.Stage.COMPLETED));
        vm.prank(maker);
        vm.expectRevert(XmrSwap.SwapCompleted.selector);
        swap.claim(sid, makerSecret);
    }

    function test_claim_afterT1_withoutReady() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.warp(block.timestamp + T1);
        vm.prank(maker);
        swap.claim(sid, makerSecret);
        assertEq(uint8(_stage(sid)), uint8(XmrSwap.Stage.COMPLETED));
    }

    function test_claim_afterT2_rejected() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.warp(block.timestamp + T1 + T2);
        vm.prank(maker);
        vm.expectRevert(XmrSwap.TooLateToClaim.selector);
        swap.claim(sid, makerSecret);
    }

    function test_claim_wrongSecretRejected() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(taker); swap.setReady(sid);
        vm.startPrank(maker);
        vm.expectRevert(XmrSwap.InvalidSecret.selector);
        swap.claim(sid, takerSecret);
        vm.expectRevert(XmrSwap.InvalidSecret.selector);
        swap.claim(sid, bytes32(0));
        bytes32 flipped = makerSecret ^ bytes32(uint256(1) << 200);
        vm.expectRevert(XmrSwap.InvalidSecret.selector);
        swap.claim(sid, flipped);
        vm.stopPrank();
    }

    function test_claim_onlyMaker() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(taker); swap.setReady(sid);
        vm.prank(stranger);
        vm.expectRevert(XmrSwap.NotMaker.selector);
        swap.claim(sid, makerSecret);
    }

    function test_claim_token_paysInToken() public {
        bytes32 id = _postTokenOffer();
        vm.startPrank(taker);
        token.approve(address(swap), 500e6);
        bytes32 sid = swap.takeOffer(id, 500e6, takerPub, TAKER_VIEW);
        swap.setReady(sid);
        vm.stopPrank();
        vm.prank(maker);
        swap.claim(sid, makerSecret);
        uint256 fee = 500e6 * uint256(FEE_BPS) / 10_000; // 0.75 MUSD
        assertEq(token.balanceOf(makerPayout), 500e6 - fee);
        assertEq(token.balanceOf(feeWallet), fee);
        assertEq(token.balanceOf(address(swap)), 0);
    }

    function test_claim_zeroFeeContract() public {
        XmrSwap free = new XmrSwap(0, payable(address(0)));
        vm.prank(maker);
        bytes32 id = free.postOffer(address(0), 0.01 ether, 1 ether, 15e12, uint64(block.timestamp + 1 days), T1, T2, makerPub, MAKER_VIEW, makerPayout);
        vm.prank(taker);
        bytes32 sid = free.takeOffer{value: 0.5 ether}(id, 0.5 ether, takerPub, TAKER_VIEW);
        vm.prank(taker); free.setReady(sid);
        uint256 before = makerPayout.balance;
        vm.prank(maker); free.claim(sid, makerSecret);
        assertEq(makerPayout.balance - before, 0.5 ether);
    }

    // ------------------------------------------------------------ refund

    function test_refund_beforeT1_notReady() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        uint256 before = taker.balance;
        vm.prank(taker);
        vm.expectEmit(true, false, false, true);
        emit XmrSwap.Refunded(sid, takerSecret);
        swap.refund(sid, takerSecret);
        assertEq(taker.balance - before, 0.5 ether, "full refund, no fee");
        assertEq(uint8(_stage(sid)), uint8(XmrSwap.Stage.COMPLETED));
    }

    function test_refund_blockedAfterReadyUntilT2() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(taker); swap.setReady(sid);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.NotTimeToRefund.selector);
        swap.refund(sid, takerSecret);
        vm.warp(block.timestamp + T1 + T2 - 1);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.NotTimeToRefund.selector);
        swap.refund(sid, takerSecret);
        vm.warp(block.timestamp + 1);
        vm.prank(taker);
        swap.refund(sid, takerSecret);
    }

    function test_refund_blockedBetweenT1AndT2() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.warp(block.timestamp + T1);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.NotTimeToRefund.selector);
        swap.refund(sid, takerSecret);
    }

    function test_refund_wrongSecretOrCaller() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(taker);
        vm.expectRevert(XmrSwap.InvalidSecret.selector);
        swap.refund(sid, makerSecret);
        vm.prank(stranger);
        vm.expectRevert(XmrSwap.NotTaker.selector);
        swap.refund(sid, takerSecret);
    }

    function test_refund_token() public {
        bytes32 id = _postTokenOffer();
        vm.startPrank(taker);
        token.approve(address(swap), 500e6);
        bytes32 sid = swap.takeOffer(id, 500e6, takerPub, TAKER_VIEW);
        uint256 before = token.balanceOf(taker);
        swap.refund(sid, takerSecret);
        vm.stopPrank();
        assertEq(token.balanceOf(taker) - before, 500e6);
    }

    function test_claimAfterRefundFails() public {
        bytes32 sid = _takeEth(_postEthOffer(), 0.5 ether);
        vm.prank(taker); swap.refund(sid, takerSecret);
        vm.warp(block.timestamp + T1);
        vm.prank(maker);
        vm.expectRevert(XmrSwap.SwapCompleted.selector);
        swap.claim(sid, makerSecret);
    }

    // ------------------------------------------------------------ key check vs Monero vectors

    function test_secretMatches_allVectors() public view {
        uint256 n = secrets.length;
        assertEq(n, KeyVectors.N);
        for (uint256 i = 0; i < n; i++) {
            assertTrue(swap.secretMatches(secrets[i], pubs[i]), "vector should match");
            assertFalse(swap.secretMatches(secrets[i], pubs[(i + 1) % n]), "other pub must not match");
        }
    }

    function test_secretMatches_basePoint() public view {
        // scalar 1 -> ed25519 base point, standard compressed encoding
        assertEq(pubs[0], bytes32(0x5866666666666666666666666666666666666666666666666666666666666666));
        assertTrue(swap.secretMatches(secrets[0], pubs[0]));
    }

    function test_secretMatches_rejectsZero() public view {
        assertFalse(swap.secretMatches(bytes32(0), pubs[0]));
    }

    function test_secretMatches_singleBitFlips() public view {
        for (uint256 b = 0; b < 256; b += 37) {
            bytes32 flipped = secrets[1] ^ bytes32(uint256(1) << b);
            assertFalse(swap.secretMatches(flipped, pubs[1]));
        }
    }

    // ------------------------------------------------------------ fuzz

    function testFuzz_claimFeeSplitExact(uint96 amount) public {
        amount = uint96(bound(amount, 0.01 ether, 1 ether));
        bytes32 sid = _takeEth(_postEthOffer(), amount);
        vm.prank(taker); swap.setReady(sid);
        uint256 mb = makerPayout.balance; uint256 fb = feeWallet.balance;
        vm.prank(maker); swap.claim(sid, makerSecret);
        uint256 fee = uint256(amount) * uint256(FEE_BPS) / 10_000;
        assertEq(makerPayout.balance - mb, amount - fee);
        assertEq(feeWallet.balance - fb, fee);
        assertEq(address(swap).balance, 0);
    }

    function testFuzz_randomSecretNeverMatches(bytes32 r) public view {
        vm.assume(r != secrets[1]);
        assertFalse(swap.secretMatches(r, pubs[1]));
    }
}
