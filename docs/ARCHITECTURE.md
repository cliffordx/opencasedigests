# Case Digestor Protocol — Technical Architecture Specification

> **Historical design — status corrected 20 October 2026.** The implementation does not yet provide the Cosmos modules, authenticated signing, JCS serialization, persistence or integrated indexing described below. See the [independent audit](audits/2026-10-20-independent-audit.md) for current evidence and [ADR-001 proposal](ARCHITECTURE_V2_PROPOSAL.md) for the recommended next design. Institutional names below are prospective examples, not partners. This document is preserved as research history, not a deployment specification.

**Document Version:** 1.0.0-draft  
**Last Updated:** October 2026  
**Status:** Phase 1A Implementation  
**Target Platform:** Cosmos SDK (Go 1.22+, CometBFT 0.38+), IPFS (Kubo), FastAPI, Supabase PostgreSQL (`pgvector`), Next.js App Router  

---

## 1. Executive Summary & Problem Space

Philippine legal informatics faces three systemic challenges:
1. **Centralized & Ephemeral Legal Repositories:** Official Supreme Court promulgations are scattered across proprietary databases, PDFs, and static websites with no verifiable cryptographic audit trail.
2. **AI Hallucinations in Legal Workflows:** Law students and practitioners increasingly use generative AI to synthesize case digests without formal provenance or verification mechanisms, risking doctrinal drift.
3. **Lack of Immutable Citation Graph:** Overruled, distinguished, or abandoned doctrines (e.g., shifts from *People v. Santos* to subsequent en banc doctrines) are manually curated in siloed textbook editions rather than indexed on an immutable, verifiable ledger.

**Case Digestor Protocol** solves this by creating a decentralized legal provenance layer and peer-attested jurisprudence registry. Modeled on the [OpenGovChain](https://github.com/bettergovph/govchain) architecture, Case Digestor bridges sovereign state legal promulgations with decentralized content addressing (IPFS) and deterministic state consensus (Cosmos SDK).

---

## 2. Dual-Phase System Architecture

The architecture is divided into two distinct operational phases to ensure zero disruption to the existing Case Digestor PWA while progressively decentralizing consensus.

```
+----------------------------------------------------------------------------------------------------+
|                                    PHASE 1: HYBRID ANCHORING (CURRENT)                            |
|                                                                                                    |
|  [Supreme Court PDF / PWA]                                                                         |
|            │                                                                                       |
|            ▼                                                                                       |
|  [FastAPI Indexer Worker] ──(RFC 8785 Canonical JSON)──► [Kubo IPFS Node]                          |
|            │                                                      │ (CIDv1)                        |
|            ▼ (Compute SHA-256)                                    ▼                                |
|  [Zero-Gas MsgExec Envelope] ──────────► [casedigestord Appchain Node]                             |
|                                                   │                                                |
|                                                   ▼ (CometBFT Event Emit)                          |
|                                          [CometBFT WS Subscriber] ──► [PostgreSQL (pgvector)]      |
+----------------------------------------------------------------------------------------------------+
|                                  PHASE 2: FEDERATED LEGAL COMMONS (FUTURE)                         |
|                                                                                                    |
|  [Attorney / Student WebAuthn Signer]                                                             |
|            │                                                                                       |
|            ▼ (Direct Gasless Tx via FeeGrant Sponsor)                                              |
|  [casedigestord Appchain & CosmWasm Escrows] ◄──► [Federated Law School Validators]                |
|            │                                                                                       |
|            ▼ (CometBFT Block Finality)                                                             |
|  [PWA Client App via Indexer Query Node]                                                           |
+----------------------------------------------------------------------------------------------------+
```

### Phase 1: Hybrid Anchoring (Bootstrap)
- **Primary Data Flow:** Write-to-DB -> Pin-to-IPFS -> Anchor-to-Chain (`casedigestord`).
- **User Impact:** Zero user friction or latency regression on the existing PWA.
- **Role of the Chain:** An append-only cryptographic notary securing official G.R. numbers, IPFS CIDv1 addresses, and 64-character hex SHA-256 checksums.

### Phase 2: Federated Legal Commons
- **Primary Data Flow:** Write-to-Chain -> CometBFT WebSocket Indexer -> Supabase PostgreSQL (`pgvector`).
- **Validators:** Run by law schools (UP College of Law, Ateneo Law, San Beda, UST), legal aid societies, and the Integrated Bar of the Philippines (IBP).
- **Smart Contracts:** CosmWasm (`x/wasm`) smart contracts governing legal research grants and institutional rubric configurations.

---

## 3. Cosmos SDK Module Architecture (`chain/`)

The application chain is compiled as `casedigestord` using Cosmos SDK v0.50+ and CometBFT v0.38+.

### 3.1 `x/jurisprudence` Module
Responsible for anchoring official Philippine Supreme Court decisions.
- **State Store Key:** `0x01 | gr_number_bytes`
- **Stored Object:** `SupremeCourtDecision`
  - `gr_number`: Supreme Court docket identifier (e.g., `G.R. No. 212345`).
  - `title`: Short title of the ruling (e.g., `People v. Dela Cruz`).
  - `promulgation_date`: Timestamp of official promulgation.
  - `division_or_enbanc`: `FIRST_DIVISION`, `SECOND_DIVISION`, `THIRD_DIVISION`, or `EN_BANC`.
  - `ponente`: Justice who penned the majority opinion.
  - `ipfs_cid`: CIDv1 string of the canonical document.
  - `sha256_checksum`: 64-character lowercase hex SHA-256 hash.
  - `anchored_by`: Bech32 address of the anchoring authority.
  - `block_height`: Block height at execution.

### 3.2 `x/attestations` Module
Implements the **Proof-of-Doctrine** lifecycle, rubric evaluations, and directed citation graph edges.
- **Digest Lifecycle:**
  1. `DIGEST_STAGE_AI_DRAFT`: Initial draft synthesized by AI worker or legal researcher.
  2. `DIGEST_STAGE_COMMUNITY_REVIEW`: Open for student and paralegal review; curators may submit revised revisions (`MsgSubmitCuratorRevision`).
  3. `DIGEST_STAGE_BAR_ATTESTED`: Verified attorney (IBP member) attests accuracy and citation integrity with scores $\ge 4/5$.
  4. `DIGEST_STAGE_MODIFIED_OR_ABANDONED`: Linked by citation edge when a subsequent SC en banc or division ruling alters or reverses the doctrine.
- **Citation Graph Edges:**
  Directed edges between digests with explicit relation types:
  - `CITATION_RELATION_AFFIRMS`
  - `CITATION_RELATION_DISTINGUISHES`
  - `CITATION_RELATION_MODIFIES`
  - `CITATION_RELATION_ABANDONED`

---

## 4. Tokenless & Gasless Protocol Architecture

Philippine law students and legal practitioners must **never** be required to hold cryptocurrency or pay gas fees. This is enforced through a four-part mechanism:

```
[User App: Next.js Client]
        │
        ▼ 1. Construct Inner Msg (e.g., MsgAttestDigest)
[User Signs Inner Msg with Ephemeral / WebAuthn Key]
        │
        ▼ 2. Envelope into cosmos.authz.v1beta1.MsgExec
[Relayer / Sponsoring Agency Node]
        │ (Signs outer tx, fee = 0uosmo or sponsored fee)
        ▼
[casedigestord AnteHandler Chain]
        │
        ├─► [ZeroGasAgencyDecorator]
        │       ├─ Check gas limit <= 500,000 (MaxZeroFeeGas)
        │       ├─ Unwrap MsgExec if present
        │       ├─ Verify ALL messages are in zeroFeeAllowedMessages whitelist
        │       └─ Mark tx as zero-fee exempt
        │
        ├─► [authz AnteHandler] (validates granter authorization)
        │
        └─► [DeductFeeDecorator] (bypassed for zero-fee exempt txs)
```

### Key Components:
1. **`ZeroGasAgencyDecorator` (`chain/app/ante/zero_gas_decorator.go`):**
   - Intercepts transactions before fee deduction.
   - Enforces a hard ceiling: `MaxZeroFeeGas = 500_000` to prevent gas exhaustion attacks.
   - Unwraps `cosmos.authz.v1beta1.MsgExec` envelopes recursively to inspect inner messages.
   - Rejects any transaction containing non-whitelisted message types if fee is zero.
2. **Sponsorship via `x/feegrant` & `x/authz`:**
   - Institutional nodes grant gas allowance to verified student and counsel accounts.
   - `MsgExec` allows institutional relayers to broadcast transactions on behalf of users while preserving cryptographic non-repudiation of the inner signature.

---

## 5. Storage, IPFS, & Cryptographic Integrity

### 5.1 No On-Chain Bloat or PII Rule
The blockchain is an immutable, high-cost replication state machine. Therefore:
- **Forbidden On-Chain:** Full Supreme Court texts, PDF binaries, vector embeddings, student names, attorney emails, and Roll of Attorneys numbers.
- **Allowed On-Chain:** Bech32 public addresses, role bitflags, IPFS CIDv1 strings, and 64-character hex SHA-256 hashes.

### 5.2 Canonical JSON (RFC 8785)
To ensure SHA-256 checksums are 100% deterministic across different platforms (Python, Go, TypeScript), all JSON payloads pinned to IPFS must be canonicalized according to RFC 8785:
- UTF-8 encoding with no byte-order mark (BOM).
- Keys sorted strictly by UTF-16 code units.
- No insignificant whitespace (no indent, no trailing spaces).
- Floating-point numbers formatted per ECMAScript specification.

### 5.3 IPFS Content Addressing
- Canonical digests and raw decisions are stored on IPFS nodes using Kubo RPC (`/api/v0/add`).
- Addresses must be formatted as **CIDv1 base32** (e.g., `bafybeic...`) to guarantee future-proof CID resolution.

---

## 6. Reactive Indexer Pipeline (`indexer-worker/`)

The indexer worker acts as the synchronization bridge between CometBFT block events, IPFS, and the search database:

1. **Event Listening:** Connects to `casedigestord` via CometBFT RPC WebSocket (`/websocket`).
2. **Subscription Queries:**
   - `tm.event='Tx' AND message.action='/casedigestor.jurisprudence.v1.MsgAnchorDecision'`
   - `tm.event='Tx' AND message.action='/casedigestor.attestations.v1.MsgAttestDigest'`
3. **Data Hydration:**
   - Extracts `ipfs_cid` and `sha256_checksum` from emitted events.
   - Fetches canonical document from IPFS Kubo node.
   - Recomputes SHA-256 hash using `indexer_worker.app.core.crypto.verify_content_integrity()`. If mismatched, flags security alert and halts processing for that item.
4. **Vector Storage:** Updates Supabase PostgreSQL (`pgvector`) with clean structured summaries and vector embeddings for semantic similarity search.

---

## 7. Determinism Guardrails for State Machine Code

All code within `chain/x/` must maintain strict consensus determinism:
- **No Float Arithmetic:** Use `sdk.Dec` or fixed-point integer math. Never use `float32` or `float64`.
- **No Non-Deterministic Map Iteration:** Standard Go map iteration order is randomized. All keeper iterations must use ordered slices or prefix iteration via Cosmos SDK KVStore iterator.
- **No Wall-Clock Time:** Never call `time.Now()`. All timestamps must derive from `ctx.BlockTime()`.
- **No Uncontrolled Goroutines:** Keepers must execute sequentially within the block execution lifecycle.
