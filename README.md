# Open Case Digests / Case Digestor Protocol

**A proposed Philippine legal commons with source-linked digests, accountable review and verifiable publication history.**

[![License: MIT](https://img.shields.io/badge/Code_license-MIT-blue.svg)](LICENSE)

## Status — 20 October 2026

This repository is an **experimental prototype and research project**. It is **NOT** a production Cosmos application, **NOT** an independently operated institutional consortium, or a source of certified legal conclusions.

The mission is worth testing: make it easier to trace a Philippine legal proposition to its source, inspect who reviewed which version, understand later treatment and see corrections. A hash verifies bytes; it does not establish legal correctness or judicial authority.

## Read the current strategy

- [Proposed founding manifesto](docs/MANIFESTO.md)
- [Evidence-gated 36-month roadmap](docs/EXECUTIVE_ROADMAP.md)
- [Architecture proposal: public releases before federation](docs/ARCHITECTURE_V2_PROPOSAL.md)
- [Research sources, prior art and evidence limits](docs/SOURCE_REGISTER.md)

The recommendation is a bounded editorial and product pilot, followed by reliable public releases and independent verification. Cosmos/IPFS anchoring remains a conditional experiment.

These recommendations do not approve a production migration or change binding repository instructions. The [original architecture](docs/ARCHITECTURE.md) and [validator guide](docs/VALIDATOR_GUIDE.md) remain historical design material with status notices.

## Repository contents

| Directory | What exists |
|---|---|
| `chain/` | CometBFT ABCI JSON prototype; separate Cosmos ante decorator; protobuf specifications; narrow tests. No implemented `chain/x/` modules |
| `indexer-worker/` | SHA-256/Python JSON helpers, Kubo client, subscriber scaffold and FastAPI health endpoint. No complete anchoring/database pipeline |
| `web/` | CosmJS helper and TypeScript configuration. No working Next.js application or passkey signer |
| `contracts/` | Placeholder workspace for future research |
| `deploy/testnet/` | Local three-node scripts/Compose; not evidence of institutional participation or independent custody |
| `docs/` | Audit, research, proposed charter, architecture and roadmap |

Inspired by [OpenGovChain](https://github.com/bettergovph/govchain). That relationship is architectural inspiration, not a claimed fork lineage, security inheritance or partnership.

## Proposed public promise

For each approved public release, make its source, exact version, review scope, rights basis, correction history and known limitations inspectable. Make verification reproducible. Preserve private work. Build commercial services by making that knowledge useful.

No end user should need cryptocurrency or pay gas. Operating costs and accountable editorial work still need funding. Public availability and legal interpretation must remain distinguishable from cryptographic integrity.

## Development evidence

At the audited source revision: eight Go tests, 21 Python tests and TypeScript typechecking passed. Go vet and shell syntax checks passed. Python lint reported three test-style issues; web build/lint failed because Next is absent; three Go files had formatting differences. No live multi-node, IPFS, browser, database or production integration was established.

## Licensing and participation

The existing [MIT license](LICENSE) applies to repository code as provided. It does not grant rights to third-party legal databases, annotations, contributed digests or private user material. Public content licensing and contributor terms remain decisions to resolve before a corpus release.

Contributions should focus first on source identity, version-bound review, independently testable verification, accessible reading and correction workflows. Named universities and professional bodies in historical examples are not confirmed partners.
