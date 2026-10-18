# Institutional Validator Guide: Case Digestor Protocol

Welcome to the **Case Digestor Network**. This guide outlines how law schools, legal aid clinics, bar chapters, and research institutions can set up, secure, and operate a federated validation node for the Philippine Legal Commons.

---

## Hardware & System Requirements

Because the network operates as a **tokenless public good** (relying on Proof-of-Authority / Federated BFT consensus rather than heavy economic mining), validator nodes have minimal hardware requirements:

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **OS** | Ubuntu 22.04 LTS | Ubuntu 24.04 LTS |
| **CPU** | 4 dedicated vCPUs | 8 vCPUs |
| **RAM** | 8 GB | 16 GB |
| **Storage** | 100 GB NVMe SSD | 250 GB NVMe SSD |
| **Network** | 100 Mbps symmetric | 1 Gbps symmetric |

### Required Ports

| Port | Protocol | Purpose | Access |
|------|----------|---------|--------|
| `26656` | TCP | CometBFT P2P peer communication | Public (all validators) |
| `26657` | TCP | CometBFT RPC | Restricted (localhost / internal indexers only) |
| `1317` | TCP | Cosmos SDK REST API | Optional (localhost only) |
| `9090` | TCP | gRPC | Optional (localhost only) |

---

## Step-by-Step Installation

### 1. Install Dependencies & Go Runtime

```bash
# Update system and install build-essential, git, and curl
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl git build-essential jq wget

# Install Go 1.22+ (Required for Cosmos SDK v0.50+)
export GO_VERSION="1.22.5"
wget "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
rm "go${GO_VERSION}.linux-amd64.tar.gz"

# Add Go to PATH
cat >> ~/.bashrc << 'EOF'
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
EOF
source ~/.bashrc

# Verify installation
go version
```

### 2. Clone and Build the Chain Binary (`casedigestord`)

```bash
git clone https://github.com/cliffordx/opencasedigests.git casedigestor-chain
cd casedigestor-chain/chain

# Build the node binary
go build -o build/casedigestord ./cmd/casedigestord/
sudo cp build/casedigestord /usr/local/bin/

# Verify binary
casedigestord version
```

### 3. Initialize Node & Configure Genesis

```bash
# Initialize validator configuration
# Replace 'InstitutionName-Validator' with your institution's name
casedigestord init "CebuLaw-Validator" --chain-id casedigestor-testnet-1

# Fetch the official shared genesis.json from the network coordinator
curl -o ~/.casedigestord/config/genesis.json \
  https://raw.githubusercontent.com/cliffordx/opencasedigests/main/config/genesis-testnet.json
```

### 4. Configure Seed Peers & Networking

Edit `~/.casedigestord/config/config.toml`:

```toml
# ---- P2P Configuration ----
[p2p]
# Seed nodes provided by the network coordinator
# Replace with actual seed node addresses from the coordination channel
seeds = ""

# Persistent peers (other institutional validators)
persistent_peers = ""

# Maximum number of inbound + outbound peers
max_num_inbound_peers = 40
max_num_outbound_peers = 10

# ---- RPC Configuration ----
[rpc]
# Restrict RPC to localhost for security
laddr = "tcp://127.0.0.1:26657"
```

> **Note:** Seed peer addresses will be distributed through the institutional coordination channel (Telegram / Signal group). Contact the network coordinator to receive your peer ID assignment.

### 5. Tune `app.toml` for Performance

Edit `~/.casedigestord/config/app.toml`:

```toml
# ---- State Pruning ----
# Keep last 100 states + every 500th state for archival queries
pruning = "custom"
pruning-keep-recent = "100"
pruning-keep-every = "500"
pruning-interval = "10"

# ---- gRPC ----
[grpc]
enable = true
address = "localhost:9090"

# ---- REST API ----
[api]
enable = true
address = "tcp://localhost:1317"
```

---

## Running the Validator

### 6. Run via Systemd (Production Service)

Create a systemd service file at `/etc/systemd/system/casedigestord.service`:

```ini
[Unit]
Description=Case Digestor Validator Node
Documentation=https://github.com/cliffordx/opencasedigests/blob/main/docs/VALIDATOR_GUIDE.md
After=network-online.target
Wants=network-online.target

[Service]
User=casedigestor
Group=casedigestor
Type=simple
ExecStart=/usr/local/bin/casedigestord start
Restart=always
RestartSec=3
LimitNOFILE=65535

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=/home/casedigestor/.casedigestord

# Environment
Environment="HOME=/home/casedigestor"

[Install]
WantedBy=multi-user.target
```

```bash
# Create a dedicated service user
sudo useradd -m -s /bin/bash casedigestor
sudo cp /usr/local/bin/casedigestord /usr/local/bin/

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable casedigestord
sudo systemctl start casedigestord

# Monitor logs
sudo journalctl -u casedigestord -f -n 50
```

---

## Validator Key Management

### 7. Back Up Your Validator Keys

Your validator's identity is defined by its private key. **Losing this key means losing your validator identity permanently.**

```bash
# Critical files to back up securely
~/.casedigestord/config/priv_validator_key.json   # Validator signing key
~/.casedigestord/config/node_key.json              # P2P identity key

# Create an encrypted backup
tar czf validator-keys-backup.tar.gz \
  ~/.casedigestord/config/priv_validator_key.json \
  ~/.casedigestord/config/node_key.json

# Encrypt with GPG (recommended)
gpg --symmetric --cipher-algo AES256 validator-keys-backup.tar.gz

# Store the encrypted backup in a secure offline location
# Then remove the unencrypted archive
rm validator-keys-backup.tar.gz
```

> **⚠️ Security Warning:** Never share `priv_validator_key.json` or store it in a public repository. If compromised, an attacker can sign blocks on behalf of your institution.

---

## Automated Upgrades with Cosmovisor

### 8. Install Cosmovisor

Cosmovisor automates binary upgrades when governance proposals are passed on-chain:

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

# Initialize Cosmovisor directory structure
export DAEMON_NAME=casedigestord
export DAEMON_HOME=$HOME/.casedigestord

cosmovisor init /usr/local/bin/casedigestord
```

Update the systemd service to use Cosmovisor:

```ini
[Service]
ExecStart=/home/casedigestor/go/bin/cosmovisor run start
Environment="DAEMON_NAME=casedigestord"
Environment="DAEMON_HOME=/home/casedigestor/.casedigestord"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="UNSAFE_SKIP_BACKUP=false"
```

### Preparing for an Upgrade

When a chain upgrade proposal passes:

```bash
# Download the new binary for the upgrade height
mkdir -p ~/.casedigestord/cosmovisor/upgrades/<upgrade-name>/bin
cp /path/to/new/casedigestord ~/.casedigestord/cosmovisor/upgrades/<upgrade-name>/bin/

# Cosmovisor will automatically swap binaries at the upgrade height
```

---

## Monitoring & Observability

### 9. Health Checks

```bash
# Check sync status
casedigestord status 2>&1 | jq '.SyncInfo'

# Check validator voting power
casedigestord status 2>&1 | jq '.ValidatorInfo'

# Check connected peers
curl -s http://localhost:26657/net_info | jq '.result.n_peers'
```

### 10. Prometheus Metrics

Enable Prometheus metrics in `config.toml`:

```toml
[instrumentation]
prometheus = true
prometheus_listen_addr = "localhost:26660"
```

Key metrics to monitor:

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `cometbft_consensus_height` | Current block height | Falling behind peers |
| `cometbft_consensus_rounds` | Consensus rounds per height | > 1 consistently |
| `cometbft_consensus_validators` | Active validator count | Unexpected changes |
| `cometbft_p2p_peers` | Connected peer count | < 2 |
| `cometbft_consensus_missing_validators` | Validators missing from votes | Your node appearing |

---

## Security Hardening

### 11. Firewall Configuration

```bash
# Allow only essential ports
sudo ufw default deny incoming
sudo ufw default allow outgoing

# CometBFT P2P (required — public)
sudo ufw allow 26656/tcp

# SSH (restrict to your IP)
sudo ufw allow from <YOUR_IP>/32 to any port 22

# Enable firewall
sudo ufw enable
```

### 12. Sentry Node Architecture (Recommended)

For high-availability institutional deployments, run your validator behind sentry nodes:

```
                    ┌─────────┐
              ┌────▶│ Sentry  │◀───── Public Internet
              │     │ Node A  │
              │     └─────────┘
              │
┌───────────┐ │     ┌─────────┐
│ Validator │◀┼────▶│ Sentry  │◀───── Public Internet
│ (Private) │ │     │ Node B  │
└───────────┘ │     └─────────┘
              │
              │     ┌─────────┐
              └────▶│ Sentry  │◀───── Public Internet
                    │ Node C  │
                    └─────────┘
```

The validator node connects **only** to sentry nodes (never directly to the public P2P network), reducing attack surface.

---

## Troubleshooting

| Symptom | Likely Cause | Resolution |
|---------|-------------|------------|
| Node stuck at height 0 | Wrong genesis or no peers | Verify `genesis.json` hash, check seed peers |
| `connection refused` on 26657 | RPC bound to wrong address | Check `config.toml` `[rpc] laddr` |
| Validator not signing blocks | Key mismatch or node not caught up | Verify `priv_validator_key.json`, wait for sync |
| High memory usage | Unconfigured pruning | Set pruning in `app.toml` |
| Peer count drops to 0 | Firewall blocking 26656 | Check `ufw` / security group rules |

---

## Support & Coordination

- **Network Coordinator:** Contact via the institutional coordination channel
- **Documentation:** [`docs/`](https://github.com/cliffordx/opencasedigests/tree/main/docs)
- **Issues:** [GitHub Issues](https://github.com/cliffordx/opencasedigests/issues)
