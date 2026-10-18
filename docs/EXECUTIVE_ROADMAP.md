# Case Digestor Protocol — Executive Roadmap & Stakeholder Strategy

**Document Version:** 1.0.0  
**Target Horizon:** October 2026 – October 2027 (12-Month Phased Rollout)  
**Stakeholder Audience:** Law School Deans, Judicial IT Directors, Integrated Bar of the Philippines (IBP) Chapters, Legal Informatics Researchers, and Infrastructure Validators  

---

## 1. Vision & Strategic Alignment

The **Case Digestor Protocol** establishes an open, decentralized legal provenance infrastructure for Philippine jurisprudence. By coupling decentralized storage (IPFS) with Byzantine Fault Tolerant consensus (Cosmos SDK / CometBFT), the protocol guarantees that legal doctrines, citations, and summaries are permanently tamper-evident and freely accessible.

### Alignment with the Philippine Judicial Innovations
The protocol directly supports the principles outlined in the **Supreme Court of the Philippines Strategic Plan for Judicial Innovations 2022–2027 (SPJI)**:
- **Efficiency:** Drastically reduces redundant legal research and citation errors through verified provenance.
- **Access:** Ensures equal access to verified case digests for indigent litigants and provincial law faculties.
- **Transparency:** Establishes an unalterable public audit trail for jurisprudence citations and subsequent doctrinal modifications.

---

## 2. 12-Month Phased Execution Roadmap

```
2026 Q4 (Months 1-3)        2027 Q1 (Months 4-6)        2027 Q2 (Months 7-9)        2027 Q3 (Months 10-12)
[ Phase 1A: Bootstrap ] ──► [ Phase 1B: PWA Pilot ] ──► [ Phase 2A: Federation ] ──► [ Phase 2B: Mainnet ]
 - Hybrid Anchoring Core     - WebAuthn Gasless Client   - CosmWasm Escrows          - Cosmovisor Upgrades
 - Proto Schemas & Ante      - 3 Law School Pilots       - 5+ Institutional Nodes    - Public Legal Commons
 - FastAPI/IPFS Worker       - IBP Verified Rubrics      - Chain-First Ingestion     - SPJI Data Gateway
```

---

### Phase 1A: Hybrid Anchoring Bootstrap (Months 1–3)
**Primary Objective:** Deliver a rock-solid, zero-friction cryptographic anchoring backend for the existing Case Digestor PWA without altering user-facing UX.

- [x] **Repository & Protocol Foundations:** Standardized monorepo boundaries (`chain/`, `indexer-worker/`, `contracts/`, `web/`, `docs/`).
- [x] **Deterministic Protobuf Schemas:** Finalize `x/jurisprudence` and `x/attestations` state structures, transactions, and event payloads.
- [x] **Zero-Gas Ante Decorator:** Enforce gasless user execution via `cosmos.authz.v1beta1.MsgExec` unwrapping and a 500k gas ceiling.
- [x] **FastAPI & Kubo IPFS Indexer:** Implement RFC 8785 canonical JSON serialization, SHA-256 verification, and CometBFT WebSocket listener.
- [x] **Gasless Signing SDK:** Create TypeScript client helpers for sponsored signing in `web/lib/chain/gasless-tx.ts`.
- [x] **3-Node Testnet Cluster:** Deploy institutional testnet with 3 seed validators in a private testing consortium.

**Key Milestone Indicator (KMI):** 1,000 historical Supreme Court en banc decisions anchored on testnet with verifiable IPFS CIDv1 addresses and SHA-256 checksums.

---

### Phase 1B: PWA Integration & Law School Pilots (Months 4–6)
**Primary Objective:** Deploy gasless client signing to law students and verify attorney review workflows in university legal aid clinics.

- [ ] **WebAuthn / Passkey Signer:** Integrate browser-native Passkeys to generate and manage Cosmos secp256k1 keypairs without seed phrases.
- [ ] **Attorney Verification Integration:** Secure credential verification checking against the Supreme Court Roll of Attorneys (yielding zero-PII on-chain role attestation).
- [ ] **University Clinical Pilot:** Deploy the hybrid platform to 3 Philippine law school legal aid clinics (e.g., UP, Ateneo, San Beda).
- [ ] **Automated AI Ingestion Pipeline:** Convert raw SC promulgations into structured canonical JSON digests with human-in-the-loop review.
- [ ] **Citation Edge Explorer:** Public graphical interface showing doctrinal links (affirmed, distinguished, modified, or abandoned).

**Key Milestone Indicator (KMI):** 250 verified attorneys active; 5,000 case digests progressing from `AI_DRAFT` to `BAR_ATTESTED`.

---

### Phase 2A: Federated Legal Commons (Months 7–9)
**Primary Objective:** Transition from hybrid database writes to decentralized chain-first consensus, backed by institutional validators and smart contracts.

- [ ] **CosmWasm Research Grants Escrow (`contracts/`):** Smart contracts allowing academic institutions and foundations to sponsor digest bounties.
- [ ] **Institutional Rubrics:** Configurable on-chain evaluation criteria for university-specific doctrinal review.
- [ ] **Validator Federation Expansion:** Onboard at least 5 institutional validators hosted by university computer science faculties and law libraries.
- [ ] **Chain-First Direct Ingestion:** High-throughput CometBFT state replication directly driving Supabase `pgvector` indexing.

**Key Milestone Indicator (KMI):** Decentralized block production by 5+ independent legal and academic institutions with zero downtime.

---

### Phase 2B: Mainnet Genesis & Public Commons (Months 10–12)
**Primary Objective:** Public mainnet launch, protocol parameter DAO governance, and integration with national legal research portals.

- [ ] **Cosmovisor Automated In-Place Upgrades:** Zero-downtime consensus binary migrations for network updates.
- [ ] **Public Read-Only API Gateway:** High-performance gRPC and REST endpoints for external legal research software and university portals.
- [ ] **Formal Security Audit:** Comprehensive audit of Cosmos SDK keeper logic, AnteDecorators, and CosmWasm smart contracts.
- [ ] **Governance Module Activation:** Democratic governance council consisting of law faculty delegates, bar representatives, and validator operators.

**Key Milestone Indicator (KMI):** Complete Philippine Supreme Court jurisprudence (1901–present) anchored with tamper-evident provenance.

---

## 3. Institutional Validator Framework

Institutions joining as network validators operate under a non-commercial, public-goods charter:

| Participant Category | Role | Responsibilities |
| :--- | :--- | :--- |
| **Law Schools & Universities** | Core Consensus Validator | Run dedicated Sentry + Validator nodes; sponsor law student legal aid digests. |
| **IBP Chapters** | Attestation Authority | Verify practicing counsel identities; review and attest critical doctrinal shifts. |
| **Legal Aid Foundations** | Grant Escrow Sponsor | Fund digest research grants via CosmWasm escrow contracts for landmark public interest cases. |
| **Independent Node Operators** | Sentry & Indexer Nodes | Provide high-throughput RPC, WebSocket, and IPFS pinning gateways to the public. |

---

## 4. Governance & Doctrinal Dispute Resolution

To avoid arbitrary state changes or doctrinal censorship:
1. **No Forking of Supreme Court Precedent:** The protocol anchors official promulgated decisions verbatim as issued by the Court.
2. **Plurality of Attestation:** If two attorneys disagree on the doctrinal holding of a decision, both attestations and rubric scores remain on-chain; difference is resolved by subsequent citations and consensus rubric weighting.
3. **Decentralized Parameter Adjustments:** Minimum rubric scores, gas ceilings (`MaxZeroFeeGas`), and authorized anchoring authorities are managed exclusively through transparent on-chain governance proposals.
