// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/Test.sol";
import {XmrSwap} from "../src/XmrSwap.sol";
import {KeyVectors} from "./KeyVectors.sol";

/// @dev Property tests for the things a stateful invariant run cannot pin down: that a wrong
///      secret never unlocks money, that only the two parties can settle, that one-time keys
///      are truly one-time, and that the fee split is exact for any amount.
contract XmrSwapFuzzTest is Test {
    XmrSwap swap;

    address payable feeWallet = payable(makeAddr("fuzzFee"));
    address payable maker = payable(makeAddr("fuzzMaker"));
    address payable taker = payable(makeAddr("fuzzTaker"));
    address payable makerPayout = payable(makeAddr("fuzzPayout"));

    uint16 constant FEE_BPS = 15;
    uint64 constant T1 = 1 hours;
    uint64 constant T2 = 1 hours;
    bytes32 constant VIEW_A = bytes32(uint256(0x1111));
    bytes32 constant VIEW_B = bytes32(uint256(0x2222));

    bytes32[] secrets;
    bytes32[] pubs;
    bytes32 makerSecret;
    bytes32 makerPub;
    bytes32 takerSecret;
    bytes32 takerPub;

    function setUp() public {
        swap = new XmrSwap(FEE_BPS, feeWallet);
        secrets = KeyVectors.secrets();
        pubs = KeyVectors.pubs();
        makerSecret = secrets[1];
        makerPub = pubs[1];
        takerSecret = secrets[2];
        takerPub = pubs[2];
        vm.deal(taker, 10_000 ether);
    }

    function _postOffer(uint128 minA, uint128 maxA, bytes32 mPub) internal returns (bytes32) {
        vm.prank(maker);
        return swap.postOffer(
            address(0), minA, maxA, 15e12, uint64(block.timestamp + 7 days), T1, T2, mPub, VIEW_A, makerPayout
        );
    }

    function _openSwap(uint256 amount) internal returns (bytes32 swapId) {
        bytes32 offerId = _postOffer(1, uint128(10_000 ether), makerPub);
        vm.prank(taker);
        swapId = swap.takeOffer{value: amount}(offerId, amount, takerPub, VIEW_B);
    }

    // ------------------------------------------------------------ secrets

    /// No secret other than the committed one may ever release the maker's payment.
    function testFuzz_claimWithWrongSecretReverts(bytes32 badSecret) public {
        vm.assume(badSecret != makerSecret);
        bytes32 id = _openSwap(1 ether);
        vm.warp(block.timestamp + T1 + 1); // past t1, so timing is not what blocks it
        vm.prank(maker);
        vm.expectRevert(XmrSwap.InvalidSecret.selector);
        swap.claim(id, badSecret);
    }

    /// Same on the refund side: only the taker's committed secret gets the money back.
    function testFuzz_refundWithWrongSecretReverts(bytes32 badSecret) public {
        vm.assume(badSecret != takerSecret);
        bytes32 id = _openSwap(1 ether);
        vm.warp(block.timestamp + T1 + T2 + 1); // past t2, so timing is not what blocks it
        vm.prank(taker);
        vm.expectRevert(XmrSwap.InvalidSecret.selector);
        swap.refund(id, badSecret);
    }

    // ------------------------------------------------------------ authorisation

    /// Nobody but the maker can claim, and nobody but the taker can refund, secret or not.
    function testFuzz_strangerCannotSettle(address who) public {
        vm.assume(who != maker && who != taker && who != address(0));
        bytes32 id = _openSwap(1 ether);
        vm.warp(block.timestamp + T1 + 1);

        vm.prank(who);
        vm.expectRevert(XmrSwap.NotMaker.selector);
        swap.claim(id, makerSecret);

        vm.warp(block.timestamp + T2 + 1);
        vm.prank(who);
        vm.expectRevert(XmrSwap.NotTaker.selector);
        swap.refund(id, takerSecret);
    }

    /// The taker cannot set ready on someone else's swap.
    function testFuzz_strangerCannotSetReady(address who) public {
        vm.assume(who != taker && who != address(0));
        bytes32 id = _openSwap(1 ether);
        vm.prank(who);
        vm.expectRevert(XmrSwap.NotTaker.selector);
        swap.setReady(id);
    }

    // ------------------------------------------------------------ one-time keys

    /// A spend key is single use. Reusing one on a second offer must be rejected, otherwise a
    /// revealed secret from swap one would unlock the Monero behind swap two.
    function testFuzz_makerKeyCannotBeReused(uint256 vecSeed) public {
        uint256 i = bound(vecSeed, 3, KeyVectors.N - 1);
        bytes32 pub = pubs[i];
        _postOffer(1, uint128(1 ether), pub);
        vm.expectRevert(XmrSwap.KeyAlreadyUsed.selector);
        _postOffer(1, uint128(1 ether), pub);
    }

    /// Same for the taker's key across two different offers.
    function testFuzz_takerKeyCannotBeReused(uint256 vecSeed) public {
        uint256 i = bound(vecSeed, 3, KeyVectors.N - 1);
        bytes32 offerA = _postOffer(1, uint128(1 ether), pubs[1]);
        bytes32 offerB = _postOffer(1, uint128(1 ether), pubs[2]);

        vm.prank(taker);
        swap.takeOffer{value: 1 ether}(offerA, 1 ether, pubs[i], VIEW_B);

        vm.prank(taker);
        vm.expectRevert(XmrSwap.KeyAlreadyUsed.selector);
        swap.takeOffer{value: 1 ether}(offerB, 1 ether, pubs[i], VIEW_B);
    }

    // ------------------------------------------------------------ money

    /// For any amount, payout + fee equals the locked value exactly, and the fee is the stated
    /// basis points with no rounding leak.
    function testFuzz_claimSplitsExactly(uint256 amountSeed) public {
        uint256 amount = bound(amountSeed, 1, 5_000 ether);
        bytes32 id = _openSwap(amount);
        vm.warp(block.timestamp + T1 + 1);

        uint256 payoutBefore = makerPayout.balance;
        uint256 feeBefore = feeWallet.balance;

        vm.prank(maker);
        swap.claim(id, makerSecret);

        uint256 gotPayout = makerPayout.balance - payoutBefore;
        uint256 gotFee = feeWallet.balance - feeBefore;

        assertEq(gotPayout + gotFee, amount, "payout + fee != locked");
        assertEq(gotFee, (amount * FEE_BPS) / 10_000, "fee is not exactly feeBps");
        assertLe(gotFee * 10_000, amount * swap.MAX_FEE_BPS(), "fee above cap");
        assertEq(address(swap).balance, 0, "contract kept a remainder");
    }

    /// A refund hands back every wei, with no fee taken.
    function testFuzz_refundReturnsEverything(uint256 amountSeed) public {
        uint256 amount = bound(amountSeed, 1, 5_000 ether);
        bytes32 id = _openSwap(amount);
        vm.warp(block.timestamp + T1 + T2 + 1);

        uint256 before = taker.balance;
        uint256 feeBefore = feeWallet.balance;

        vm.prank(taker);
        swap.refund(id, takerSecret);

        assertEq(taker.balance - before, amount, "refund shorted the taker");
        assertEq(feeWallet.balance, feeBefore, "a fee was taken on a refund");
        assertEq(address(swap).balance, 0, "contract kept a remainder");
    }

    /// Paying outside the offer's range is always rejected.
    function testFuzz_amountOutOfRangeReverts(uint256 amountSeed) public {
        uint128 minA = 1 ether;
        uint128 maxA = 2 ether;
        bytes32 offerId = _postOffer(minA, maxA, makerPub);
        uint256 amount = bound(amountSeed, 0, 10 ether);
        vm.assume(amount < minA || amount > maxA);

        vm.prank(taker);
        vm.expectRevert(XmrSwap.AmountOutOfRange.selector);
        swap.takeOffer{value: amount}(offerId, amount, takerPub, VIEW_B);
    }

    /// Declaring one amount and sending a different quantity of ETH is always rejected.
    function testFuzz_mismatchedValueReverts(uint256 sentSeed) public {
        bytes32 offerId = _postOffer(1, uint128(10 ether), makerPub);
        uint256 declared = 1 ether;
        uint256 sent = bound(sentSeed, 0, 10 ether);
        vm.assume(sent != declared);

        vm.prank(taker);
        vm.expectRevert(XmrSwap.WrongValue.selector);
        swap.takeOffer{value: sent}(offerId, declared, takerPub, VIEW_B);
    }
}
