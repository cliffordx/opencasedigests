# Case Digestor Protocol — CosmWasm Smart Contracts (`contracts/`)

This directory houses the CosmWasm (`x/wasm`) smart contracts targeted for **Phase 2A (Federated Legal Commons)** of the Case Digestor Protocol.

---

## 1. Overview & Objectives

In Phase 2, smart contracts deployed on `casedigestord` will govern:
1. **Research Grant Escrows (`escrow-grants`):**
   - Enables academic institutions, foundations, and public donors to post bounties for high-priority legal research (e.g., digesting complex en banc antitrust, constitutional, or maritime law cases).
   - Milestone payouts are triggered automatically upon verified attorney attestations (`BAR_ATTESTED` state transition).
2. **Institutional Rubrics (`institutional-rubrics`):**
   - Configurable review criteria for universities (e.g., law school grading schemes, specialized citation formats, and minimum passing rubric thresholds).

---

## 2. Directory Layout (Phase 2A Workspace)

```
contracts/
├── Cargo.toml                  # Cargo workspace manifest
├── README.md                   # This specification and setup guide
├── contracts/
│   ├── grant-escrow/           # Research grant escrow and payout logic
│   └── institutional-rubrics/  # Customizable peer-review rubric contracts
└── packages/
    └── casedigestor-core/      # Shared messages, responses, and errors
```

---

## 3. Prerequisites & Development Environment

To build and test CosmWasm contracts locally:

```bash
# Ensure Rust and WebAssembly compilation targets are installed
rustup default stable
rustup target add wasm32-unknown-unknown

# Install CosmWasm CLI tools and linter
cargo install cargo-run-script
cargo install cosmwasm-check
```

### Compiling Optimized WASM Binaries
For deployment to `casedigestord` (`x/wasm`):

```bash
# Using Docker workspace optimizer
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/workspace-optimizer:0.16.0
```

---

## 4. Architectural Guardrails
- **Deterministic Execution:** No floating-point operations; use `Uint128` and `Decimal` types from `cosmwasm_std`.
- **Zero-Gas Compatibility:** Grant creation and dispute resolutions must be compatible with `ZeroGasAgencyDecorator` or `x/feegrant` relays.
- **No PII Stored in Contract State:** All student and attorney references must be Bech32 addresses.
