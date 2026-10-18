# Case Digestor Protocol

**Decentralized Legal Provenance & Peer-Attested Jurisprudence for the Philippine Legal Commons**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Chain](https://img.shields.io/badge/Cosmos_SDK-v0.50-6F3CE5.svg)](chain/)
[![Indexer](https://img.shields.io/badge/FastAPI-Python_3.11-009688.svg)](indexer-worker/)

Case Digestor Protocol is a decentralized legal infrastructure layer that transforms how Philippine Supreme Court jurisprudence is digested, verified, and preserved. It combines content-addressed IPFS storage with a Cosmos SDK appchain (`casedigestord`) to create tamper-evident, community-peer-reviewed case digests anchored on a federated blockchain — without requiring any user to hold cryptocurrency or pay gas fees.

Inspired by the [OpenGovChain](https://github.com/bettergovph/govchain) architecture.

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          Case Digestor Protocol                         │
│                                                                         │
│  ┌─────────────┐   ┌──────────────────┐   ┌──────────────────────────┐  │
│  │  web/       │   │ indexer-worker/   │   │  chain/ (casedigestord)  │  │
│  │  Next.js    │   │ FastAPI + IPFS    │   │  Cosmos SDK v0.50        │  │
│  │  App Router │   │ Kubo RPC Client   │   │  CometBFT Consensus     │  │
│  │  CosmJS     │──▶│ SHA-256 Verify    │──▶│                          │  │
│  │  Gasless TX │   │ CometBFT WS Sub   │◀──│  x/jurisprudence         │  │
│  │  WebAuthn   │   │ pgvector Index    │   │  x/attestations          │  │
│  └──────┬──────┘   └────────┬─────────┘   │  app/ante/ZeroGas        │  │
│         │                   │              │  x/authz + x/feegrant    │  │
│         │                   │              └────────────┬─────────────┘  │
│         │                   │                           │                │
│         ▼                   ▼                           ▼                │
│  ┌─────────────┐   ┌──────────────────┐   ┌──────────────────────────┐  │
│  │  Passkey /  │   │  Supabase        │   │  IPFS (Kubo)             │  │
│  │  WebAuthn   │   │  PostgreSQL      │   │  Content-Addressed       │  │
│  │  Signer     │   │  pgvector        │   │  CIDv1 Storage           │  │
│  └─────────────┘   └──────────────────┘   └──────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
```

### Data Flow (Phase 1 — Hybrid Anchoring)

```
  Supreme Court          AI Digest            Peer Review            On-Chain
  Decision (PDF)         Generation           & Attestation          Anchoring
  ─────────────          ──────────           ─────────────          ──────────
       │                      │                     │                     │
       ▼                      ▼                     ▼                     ▼
  ┌─────────┐          ┌───────────┐         ┌───────────┐         ┌───────────┐
  │ Scrape  │          │ Generate  │         │ Attorney  │         │ Anchor    │
  │ from SC │───PDF───▶│ Digest    │──IPFS──▶│ Review +  │──Tx───▶│ CIDv1 +   │
  │ website │          │ (AI)      │  Pin    │ Attest    │         │ SHA-256   │
  └─────────┘          └───────────┘         └───────────┘         └───────────┘
       │                      │                     │                     │
       ▼                      ▼                     ▼                     ▼
  Write to DB           AI_DRAFT             BAR_ATTESTED          Immutable
  (Supabase)            Stage                Stage                 Provenance
```

---

## The "Proof-of-Doctrine" Lifecycle

Every case digest progresses through a **4-stage peer-review lifecycle** before achieving immutable provenance status:

| Stage | Name | Description | Who Advances |
|-------|------|-------------|--------------|
| **1** | `AI_DRAFT` | Initial AI-generated digest is submitted and pinned to IPFS with a SHA-256 checksum anchored on-chain. | AI Coordinator / Contributor |
| **2** | `COMMUNITY_REVIEW` | Law students, paralegals, and legal researchers review the digest for accuracy, completeness, and citation integrity. | Community Curators (≥3 reviews) |
| **3** | `BAR_ATTESTED` | A verified member of the Philippine Bar (IBP-accredited attorney) provides a formal attestation with rubric scores for doctrine accuracy and citation integrity. | Licensed Attorney (ROLE_COUNSEL_REVIEWER) |
| **4** | `MODIFIED_OR_ABANDONED` | A subsequent En Banc or Division ruling modifies, distinguishes, or abandons the doctrine. A citation edge is created linking the old and new rulings. | Any verified contributor with proof CID |

### Citation Edges

When a newer Supreme Court decision affects an existing digest, a **Citation Edge** is created with one of these relations:

| Relation | Meaning |
|----------|---------|
| `AFFIRMS` | The newer ruling affirms and strengthens the existing doctrine |
| `DISTINGUISHES` | The newer ruling narrows or distinguishes the doctrine's applicability |
| `MODIFIES` | The newer ruling partially overrules or amends the doctrine |
| `ABANDONED` | The newer ruling completely overrules and abandons the doctrine |

---

## Repository Structure

```
opencasedigests/
├── chain/                          # Cosmos SDK appchain (casedigestord)
│   ├── proto/casedigestor/         #   Protobuf definitions
│   │   ├── jurisprudence/v1/       #     Supreme Court decision anchoring
│   │   └── attestations/v1/        #     Digest lifecycle & peer review
│   └── app/ante/                   #   ZeroGasDecorator (gasless transactions)
├── indexer-worker/                 # Python FastAPI microservice
│   ├── app/core/                   #   SHA-256 & canonical hashing
│   └── app/services/               #   IPFS client & CometBFT subscriber
├── contracts/                      # Rust CosmWasm smart contracts (Phase 2)
├── web/                            # Next.js frontend + CosmJS gasless signing
│   └── lib/chain/                  #   TypeScript authz/feegrant helpers
└── docs/                           # Public stakeholder documentation
```

---

## Key Design Principles

### Tokenless & Gasless by Default
No user — whether a law student, paralegal, or attorney — is ever required to hold cryptocurrency or pay gas fees. All state-changing transactions are wrapped inside `cosmos.authz.v1beta1.MsgExec` with the Case Digestor Master Authority as `fee.granter` via `x/feegrant`. The custom `ZeroGasAgencyDecorator` enforces a hard ceiling of **500,000 gas** per zero-fee transaction.

### No On-Chain Bloat
Full decision text, raw PDFs, vector embeddings, and personally identifiable information are **never** stored on-chain. The chain stores only:
- Bech32 addresses and role flags
- IPFS CIDv1 content identifiers
- 64-character hex SHA-256 checksums
- Digest stage transitions and rubric scores

### Tamper-Evident Provenance
Every digest and attestation is pinned to IPFS and its SHA-256 checksum is anchored on-chain. Any modification to the IPFS content will produce a different hash, immediately detectable by any verifier.

---

## 12-Month Phased Roadmap

### Phase 1A: Hybrid Anchoring Bootstrap (Months 1–3)
- [x] Repository scaffold and governance documentation
- [ ] Protobuf schema definitions for `x/jurisprudence` and `x/attestations`
- [ ] `ZeroGasAgencyDecorator` ante handler
- [ ] FastAPI indexer-worker with IPFS pinning and SHA-256 verification
- [ ] CometBFT WebSocket subscriber for event-driven indexing
- [ ] TypeScript gasless signing helper (`web/lib/chain/`)
- [ ] Testnet genesis with 3 institutional validators

### Phase 1B: PWA Integration & Community Tools (Months 4–6)
- [ ] Next.js App Router frontend with WebAuthn/Passkey signer
- [ ] AI digest generation pipeline (SC decision → structured JSON)
- [ ] Community review dashboard with rubric scoring UI
- [ ] Supabase `pgvector` semantic search over attested digests
- [ ] Attorney verification flow (IBP Roll number → on-chain role grant)
- [ ] Public attestation explorer and citation graph visualization

### Phase 2A: Federated Legal Commons (Months 7–9)
- [ ] Migrate from hybrid (DB-first) to chain-first writes
- [ ] CometBFT WebSocket indexer → PostgreSQL reactive pipeline
- [ ] CosmWasm research grant escrow contracts (`contracts/`)
- [ ] Institutional rubric smart contracts for customizable review criteria
- [ ] Cross-institutional validator federation (5+ law schools)

### Phase 2B: Production & Governance (Months 10–12)
- [ ] Mainnet launch with Cosmovisor automated upgrades
- [ ] Governance module for protocol parameter changes
- [ ] API gateway for third-party legal research integrations
- [ ] Comprehensive documentation and onboarding materials
- [ ] Security audit and penetration testing
- [ ] Public launch and institutional partnership announcements

---

## Getting Started

### Prerequisites
- **Go 1.22+** (for `chain/`)
- **Python 3.11+** (for `indexer-worker/`)
- **Node.js 20+** and **pnpm** (for `web/`)
- **IPFS Kubo** node (for content-addressed storage)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/cliffordx/opencasedigests.git
cd opencasedigests

# Build the chain binary
cd chain && go build -o build/casedigestord ./cmd/casedigestord/

# Set up the indexer worker
cd ../indexer-worker && python -m venv .venv && source .venv/bin/activate
pip install -e ".[dev]"

# Run tests
pytest tests/
```

See [`docs/VALIDATOR_GUIDE.md`](docs/VALIDATOR_GUIDE.md) for full validator node setup instructions.

---

## Contributing

We welcome contributions from law students, legal practitioners, legal engineers, and anyone passionate about open access to Philippine jurisprudence. Please read our contribution guidelines before submitting pull requests.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
