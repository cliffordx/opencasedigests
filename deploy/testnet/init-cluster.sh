#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

echo "=== Case Digestor Protocol: Initializing 3-Node Testnet Cluster ==="

# 1. Build host binary if not present
if [ ! -f "${REPO_ROOT}/chain/build/casedigestord" ]; then
    echo "Building casedigestord binary..."
    (cd "${REPO_ROOT}/chain" && go build -o build/casedigestord ./cmd/casedigestord)
fi

# 2. Generate clean testnet configurations
echo "Generating consortium validator configs (UP Law, Ateneo Law, San Beda Law)..."
rm -rf "${SCRIPT_DIR}/data"
"${REPO_ROOT}/chain/build/casedigestord" testnet \
    --output-dir "${SCRIPT_DIR}/data" \
    --docker=true \
    --chain-id "casedigestor-testnet-1"

# 3. Build docker images and start cluster
echo "Starting testnet cluster via Docker Compose..."
cd "${SCRIPT_DIR}"
docker compose build
docker compose up -d

echo "=== Testnet Cluster Started Successfully ==="
echo "Validator 0 (UP Law):       RPC: http://localhost:26657 | P2P: 26656"
echo "Validator 1 (Ateneo Law):   RPC: http://localhost:26667 | P2P: 26666"
echo "Validator 2 (San Beda Law): RPC: http://localhost:26677 | P2P: 26676"
echo "IPFS Gateway:               RPC: http://localhost:5001  | Gateway: http://localhost:8080"
echo ""
echo "Run '${SCRIPT_DIR}/verify-cluster.sh' to check network health."
