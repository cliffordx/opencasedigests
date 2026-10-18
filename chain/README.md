# Case Digestor Protocol — Cosmos SDK Appchain (`chain/`)

This directory houses the Cosmos SDK application chain (`casedigestord`), custom AnteHandlers, and Protobuf schema definitions for the **Case Digestor Protocol**.

---

## 1. Architecture & Core Modules

The chain provides the cryptographic provenance and peer-attestation registry for Philippine Supreme Court decisions.

- **`proto/casedigestor/`**: Canonical Protobuf schemas:
  - `jurisprudence/v1/`: Decision anchoring (`SupremeCourtDecision`, `MsgAnchorDecision`).
  - `attestations/v1/`: Case digest lifecycle (`DigestRecord`), peer attestations (`PeerAttestation`), and citation edges (`CitationEdge`).
- **`app/ante/`**: Custom AnteHandlers:
  - `zero_gas_decorator.go`: Enforces zero-fee execution for authorized legal practitioner and curator messages wrapped in `cosmos.authz.v1beta1.MsgExec` with a hard gas ceiling (`MaxZeroFeeGas = 500,000`).

---

## 2. Deterministic Consensus Guardrails
All keeper logic and state-transition code must adhere to strict determinism:
1. **No Float Arithmetic:** All mathematical operations must use `math.Int` or `sdk.Dec`. Never use `float32` or `float64`.
2. **No Standard Go Map Iteration:** Keepers must use ordered KVStore iterators or sorted slices.
3. **No Wall-Clock Time:** Never call `time.Now()`. Derive all timestamps from `ctx.BlockTime()`.
4. **No On-Chain PII or Bloat:** Never store attorney names, emails, Roll of Attorneys numbers, or decision text on-chain. Only store Bech32 addresses, role bitflags, IPFS CIDv1 strings, and 64-character hex SHA-256 checksums.

---

## 3. Tooling & Development

```bash
# Verify Protobuf schemas with Buf
buf lint

# Generate Go and gRPC stubs
buf generate

# Compile chain binary (Phase 1B)
go build -o build/casedigestord ./cmd/casedigestord
```
