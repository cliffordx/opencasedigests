# Case Digestor Protocol — Web Client SDK & PWA (`web/`)

> **Design sketch, 20 October 2026:** only a TypeScript helper exists. There is no Next.js application or passkey signer, and the helper's encoded payload does not interoperate with the current daemon. The example APIs below are illustrative and do not match the implemented exports. See the [audit](../docs/audits/2026-10-20-independent-audit.md); this is not a working integration guide.

This directory provides the TypeScript client libraries and Next.js App Router progressive web application (PWA) for the **Case Digestor Protocol**.

---

## 1. Features & Architecture

- **Tokenless & Gasless Transactions:** Implements `createGaslessDigestDraftTx`, `createGaslessAttestationTx`, and `createGaslessCitationLinkTx` using `x/authz` (`MsgExec`) and `x/feegrant` sponsorship.
- **Embedded Signing:** Zero cryptocurrency wallets required. Users sign with browser-local ephemeral keys or WebAuthn/Passkeys.
- **Strict Decoupling:** Client apps interact with IPFS for document storage and CometBFT RPC for on-chain state anchoring.

---

## 2. Directory Structure

```
web/
├── lib/
│   └── chain/
│       └── gasless-tx.ts    # CosmJS gasless transaction builder & broadcast helpers
├── package.json             # NPM package dependencies (@cosmjs/stargate, etc.)
├── tsconfig.json            # TypeScript compiler configuration
└── README.md                # Client architecture documentation
```

---

## 3. Usage Example

```typescript
import {
  createSigningClient,
  createGaslessAttestationTx,
  broadcastGaslessTx,
} from "./lib/chain/gasless-tx";

// 1. Initialize CosmJS signing client with user's offline signer
const client = await createSigningClient(rpcEndpoint, offlineSigner);

// 2. Build gasless attestation message (wrapped in MsgExec)
const gaslessMsg = createGaslessAttestationTx({
  attorneyAddress: "casedigestor1...",
  digestId: "GR-212345-v1",
  doctrineAccuracyScore: 5,
  citationIntegrityScore: 4,
  reviewNotesIpfsCid: "bafybeic...",
});

// 3. Broadcast without paying fees (sponsored by Master Authority)
const result = await broadcastGaslessTx(client, "casedigestor1...", [gaslessMsg]);
console.log("Attestation anchored at block height:", result.height);
```
