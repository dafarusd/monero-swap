#!/usr/bin/env bash
# Deploys XmrSwap and verifies the source on the chain's explorer.
#   ./deploy.sh <rpc-url> <key-file> <fee-bps> <fee-recipient> [etherscan-api-key]
# The key file holds the deployer's hex private key; it is read by forge, never printed.
set -euo pipefail
RPC="${1:?rpc url}"; KEYFILE="${2:?key file}"; FEE_BPS="${3:?fee bps}"; FEE_TO="${4:?fee recipient}"; APIKEY="${5:-}"
[ "$FEE_BPS" -le 100 ] || { echo "fee cap is 100 bps (1%)"; exit 1; }
[[ "$FEE_TO" =~ ^0x[0-9a-fA-F]{40}$ ]] || { echo "fee recipient must be a 0x address"; exit 1; }
cd "$(dirname "$0")"
forge build >/dev/null
CHAIN=$(cast chain-id --rpc-url "$RPC")
FROM=$(cast wallet address --private-key "0x$(tr -d '[:space:]' < "$KEYFILE" | sed 's/^0x//')")
echo "chain $CHAIN, deployer $FROM, balance $(cast balance "$FROM" --rpc-url "$RPC" --ether) ETH"
echo "fee ${FEE_BPS} bps -> $FEE_TO"
OUT=$(forge create src/XmrSwap.sol:XmrSwap --rpc-url "$RPC" --private-key "0x$(tr -d '[:space:]' < "$KEYFILE" | sed 's/^0x//')" --broadcast --constructor-args "$FEE_BPS" "$FEE_TO" 2>&1)
ADDR=$(echo "$OUT" | grep 'Deployed to' | awk '{print $3}')
[ -n "$ADDR" ] || { echo "$OUT" | tail -5; exit 1; }
echo "deployed: $ADDR"
echo "feeBps=$(cast call "$ADDR" 'feeBps()(uint16)' --rpc-url "$RPC")  feeRecipient=$(cast call "$ADDR" 'feeRecipient()(address)' --rpc-url "$RPC")"
if [ -n "$APIKEY" ]; then
  forge verify-contract "$ADDR" src/XmrSwap.sol:XmrSwap --chain-id "$CHAIN" --etherscan-api-key "$APIKEY" \
    --constructor-args "$(cast abi-encode 'constructor(uint16,address)' "$FEE_BPS" "$FEE_TO")" --watch || echo "verification failed; retry with: forge verify-contract $ADDR src/XmrSwap.sol:XmrSwap --chain-id $CHAIN"
fi
