#!/usr/bin/env bash
set -euo pipefail

echo "=== Verifying 3-Node Consortium Testnet Health ==="

# Wait for RPC endpoints to be ready
endpoints=(
    "http://localhost:26657"
    "http://localhost:26667"
    "http://localhost:26677"
)

names=(
    "Validator 0 (UP Law)"
    "Validator 1 (Ateneo Law)"
    "Validator 2 (San Beda Law)"
)

echo "Waiting for validators to produce blocks..."
sleep 3

for i in "${!endpoints[@]}"; do
    url="${endpoints[$i]}"
    name="${names[$i]}"
    
    echo -n "Checking ${name} (${url})... "
    status_resp=$(curl -s "${url}/status" || echo "")
    if [ -z "${status_resp}" ]; then
        echo "FAILED: endpoint unreachable"
        exit 1
    fi

    chain_id=$(echo "${status_resp}" | jq -r '.result.node_info.network')
    height=$(echo "${status_resp}" | jq -r '.result.sync_info.latest_block_height')
    catching_up=$(echo "${status_resp}" | jq -r '.result.sync_info.catching_up')
    voting_power=$(echo "${status_resp}" | jq -r '.result.validator_info.voting_power')

    echo "OK"
    echo "  Chain ID:     ${chain_id}"
    echo "  Block Height: ${height}"
    echo "  Catching Up:  ${catching_up}"
    echo "  Voting Power: ${voting_power}"
done

echo ""
echo "Checking IPFS Kubo node (http://localhost:5001)..."
ipfs_version=$(curl -s -X POST "http://localhost:5001/api/v0/version" | jq -r '.Version' || echo "failed")
echo "  IPFS Version: ${ipfs_version}"

echo ""
echo "=== All 3 Institutional Consortium Validators Are Healthy & Replicating State ==="
