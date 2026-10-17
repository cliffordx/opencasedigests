# Institutional Validator Guide: Case Digestor Protocol

Welcome to the **Case Digestor Network**. This guide outlines how law schools, legal aid clinics, bar chapters, and research institutions can set up, secure, and operate a federated validation node for the Philippine Legal Commons.

## Hardware & System Requirements

Because the network operates as a **tokenless public good** (relying on Proof-of-Authority / Federated BFT consensus rather than heavy economic mining), validator nodes have minimal hardware requirements:

* **OS:** Ubuntu 22.04 LTS (recommended) or macOS / WSL2
* **CPU:** 4+ Dedicated vCPUs
* **RAM:** 8 GB to 16 GB minimum
* **Storage:** 100 GB NVMe SSD (to accommodate blockchain state and local IPFS block storage)
* **Ports Open:** `26656` (CometBFT P2P peer communication) and `26657` (CometBFT RPC listener, restricted to localhost/internal indexers).

---

## Step-by-Step Installation

### 1. Install Dependencies & Go Runtime

```bash
# Update system and install build-essential, git, and curl
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl git build-essential jq

# Install Go 1.22+ (Required for Cosmos SDK v0.50+)
export GO_VERSION="1.22.5"
wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.bashrc
source ~/.bashrc

```

### 2. Clone and Build the Chain Binary (`casedigestor-d`)

```bash
git clone https://github.com/bettergovph/govchain.git casedigestor-chain
cd casedigestor-chain/chain

# Build the node binary
go build -o build/casedigestord ./cmd/casedigestord/
sudo cp build/casedigestord /usr/local/bin/
casedigestord version

```

### 3. Initialize Node & Configure Genesis

```bash
# Initialize validator configuration (replace 'InstitutionName-Validator')
casedigestord init "CebuLaw-Validator" --chain-id casedigestor-testnet-1

# Fetch the official shared genesis.json from the network coordinator
curl -o ~/.casedigestord/config/genesis.json https://raw.githubusercontent.com/bettergovph/govchain/main/config/genesis-testnet.json

```

### 4. Run via Systemd (Production Service)

Create a systemd service file at `/etc/systemd/system/casedigestord.service`:

```ini
[Unit]
Description=Case Digestor Validator Node
After=network.online.target

[Service]
User=$USER
ExecStart=/usr/local/bin/casedigestord start
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target

```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable casedigestord
sudo systemctl start casedigestord
sudo journalctl -u casedigestord -f -n 50

```
