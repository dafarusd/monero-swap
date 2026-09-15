#!/usr/bin/env bash
# Deploys XmrSwapV2 and verifies the source on the chain's explorer.
#   ./deploy-v2.sh <rpc-url> <key-file> <fee-bps> <founder-share-bps> <bond-bps> <founder> <dev-fund> [etherscan-api-key]
# The key file holds the deployer's hex private key; it is read by forge, never printed.
#
# Fee, split, bond and both recipients are fixed at deployment and can never be changed.
# Check every argument twice before running this against mainnet.
set -euo pipefail
RPC="${1:?rpc url}"
KEYFILE="${2:?key file}"
FEE_BPS="${3:?fee bps}"
SHARE_BPS="${4:?founder share bps}"
BOND_BPS="${5:?bond bps}"
FOUNDER="${6:?founder address}"
DEVFUND="${7:?dev fund address}"
APIKEY="${8:-}"

[ "$FEE_BPS" -le 100 ] || { echo "fee cap is 100 bps (1%)"; exit 1; }
[ "$SHARE_BPS" -le 10000 ] || { echo "founder share cap is 10000 bps"; exit 1; }
[ "$BOND_BPS" -le 2000 ] || { echo "bond cap is 2000 bps (20%)"; exit 1; }
[[ "$FOUNDER" =~ ^0x[0-9a-fA-F]{40}$ ]] || { echo "founder must be a 0x address"; exit 1; }
[[ "$DEVFUND" =~ ^0x[0-9a-fA-F]{40}$ ]] || { echo "dev fund must be a 0x address"; exit 1; }
[ "$FOUNDER" != "$DEVFUND" ] || { echo "founder and dev fund must differ"; exit 1; }

cd "$(dirname "$0")"
KEY="0x$(tr -d '[:space:]' < "$KEYFILE" | sed 's/^0x//')"

forge build >/dev/null
CHAIN=$(cast chain-id --rpc-url "$RPC")
FROM=$(cast wallet address --private-key "$KEY")
echo "chain    $CHAIN"
echo "deployer $FROM  ($(cast balance "$FROM" --rpc-url "$RPC" --ether) ETH)"
echo "fee      ${FEE_BPS} bps, split ${SHARE_BPS}/$((10000 - SHARE_BPS))"
echo "founder  $FOUNDER"
echo "devFund  $DEVFUND"
echo "bond     ${BOND_BPS} bps of an offer's ceiling"

ARGS=$(cast abi-encode 'constructor(uint16,uint16,uint16,address,address)' \
  "$FEE_BPS" "$SHARE_BPS" "$BOND_BPS" "$FOUNDER" "$DEVFUND")

for attempt in 1 2 3; do
  OUT=$(forge create src/XmrSwapV2.sol:XmrSwapV2 --rpc-url "$RPC" --private-key "$KEY" --broadcast \
    --constructor-args "$FEE_BPS" "$SHARE_BPS" "$BOND_BPS" "$FOUNDER" "$DEVFUND" 2>&1)
  ADDR=$(echo "$OUT" | grep 'Deployed to' | awk '{print $3}')
  TX=$(echo "$OUT" | grep 'Transaction hash' | awk '{print $3}')
  [ -n "$ADDR" ] && break
  echo "forge create did not report an address (attempt $attempt):"
  echo "$OUT" | grep -vE '^\s*[0-9]+ │|^\s*│|╭|╰|━|internal-function|help:' | tail -8
  sleep 5
done
[ -n "${ADDR:-}" ] || exit 1

echo
echo "deployed: $ADDR"
echo "tx:       ${TX:-unknown}"
echo

# Read the live contract back. If any of these disagree with the arguments above, stop and
# do not publish the address — the constructor values are permanent.
echo "--- on-chain state, read back from $ADDR"
echo "feeBps          $(cast call "$ADDR" 'feeBps()(uint16)' --rpc-url "$RPC")"
echo "founderShareBps $(cast call "$ADDR" 'founderShareBps()(uint16)' --rpc-url "$RPC")"
echo "bondBps         $(cast call "$ADDR" 'bondBps()(uint16)' --rpc-url "$RPC")"
echo "founder         $(cast call "$ADDR" 'founder()(address)' --rpc-url "$RPC")"
echo "devFund         $(cast call "$ADDR" 'devFund()(address)' --rpc-url "$RPC")"
echo "MIN_TIMEOUT     $(cast call "$ADDR" 'MIN_TIMEOUT()(uint64)' --rpc-url "$RPC") seconds"
echo "offerCount      $(cast call "$ADDR" 'offerCount()(uint256)' --rpc-url "$RPC")"

if [ -n "$APIKEY" ]; then
  echo
  forge verify-contract "$ADDR" src/XmrSwapV2.sol:XmrSwapV2 --chain-id "$CHAIN" \
    --etherscan-api-key "$APIKEY" --constructor-args "$ARGS" --watch \
    || echo "verification failed; retry with: forge verify-contract $ADDR src/XmrSwapV2.sol:XmrSwapV2 --chain-id $CHAIN --constructor-args $ARGS"
fi
