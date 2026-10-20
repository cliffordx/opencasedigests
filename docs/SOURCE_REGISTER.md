# Research source register — 20 October 2026

This register supports the [independent audit](audits/2026-10-20-independent-audit.md), [proposed architecture](ARCHITECTURE_V2_PROPOSAL.md), [roadmap](EXECUTIVE_ROADMAP.md), and [manifesto](MANIFESTO.md). Research cutoff: **20 October 2026, Asia/Manila**. These are public-source observations, not customer interviews, legal clearance, vendor performance tests, or partnership confirmations.

The audit uses local source snapshots of `opencasedigests` at `65d59127bed36366e2063caf71dd50c889249fed` and `case-digestor-pwa` at `984b9e2c31bc21bd426daa12c29dfaa8c37d2042`. The PWA's configured remote is `case-digestor`; its local checkout was inspected without refreshing the remote or inspecting production.

## Existing work and commercial alternatives

| ID | Primary source | What the evidence supports | Limit |
|---|---|---|---|
| M1 | [Supreme Court E-Library](https://elibrary.judiciary.gov.ph/) | An existing official Philippine legal research repository | Its existence does not establish complete coverage or permission for every downstream use |
| M2 | [Lawphil](https://lawphil.net/) and [copyright notice](https://lawphil.net/copyright.html) | Established free legal access; the notice distinguishes government works and applies a CC BY-NC 4.0 notice to its collection | Free access is not an unrestricted commercial dataset license |
| M3 | [CD Asia features](https://cdasia.com/features/) | Vendor describes cross references, Cited-in and citation-locator functionality | Marketing evidence, not an independently tested quality comparison |
| M4 | [Digest PH FAQ](https://www.digest.ph/faqs) | Vendor describes human/AI digests, shepardized cases, and community uploads rewarded with access points | Refutes broad novelty of community digest exchange; does not establish an open data license |
| M5 | [Anycase pricing](https://anycase.ai/pricing) | Observed posted Pro price: PHP999/month or PHP799/month equivalent billed yearly; Education: PHP599/month or PHP499/month equivalent billed yearly | Dynamic listed prices, not realized revenue or willingness to pay for this project; reconfirm before a pricing decision |
| M6 | [Free Law Project expert citator annotation call, 4 June 2025](https://free.law/2025/06/04/citator-annotation-volunteer/) | An open-source AI citator effort with expert-labelled citation treatments predates this proposal | US-focused; its “first” wording is the organization's claim, not our global novelty determination |
| M7 | [OpenGovChain repository](https://github.com/bettergovph/govchain) | Existing tokenless public-goods concept using Cosmos, IPFS and indexing; its README points to a separate [daemon repository](https://github.com/bettergovph/govchaind) | Architectural inspiration; not proof that this repo inherits its implementation, security, partners, or deployment |

**Novelty method:** compared official product/repository pages for public legal access, Philippine digests, community contributions, citation treatment, open legal datasets and blockchain provenance. This is a bounded public prior-art search. Private pilots, unpublished projects, patents, discontinued services and non-indexed academic work may be missing. No patentability or freedom-to-operate opinion is made.

## Philippine authority and publication constraints

| ID | Primary source | Supported point | Design consequence |
|---|---|---|---|
| L1 | [1987 Constitution, Article VIII §4(3)](https://lawphil.net/consti/cons1987.html) | Modification or reversal of doctrine laid down by the Supreme Court requires the Court sitting en banc | Do not mechanically label a Division citation as doctrinal reversal; distinguish disposition, application and treatment of a proposition |
| L2 | [RA 8293, §§175–176](https://lawphil.net/statutes/repacts/ra1997/ra_8293_1997.html) | Official legal texts and government works receive distinct treatment; §176 includes a profit-use approval provision and exceptions; third-party works are not stripped of rights by republication | Record source-specific rights; obtain advice on the exact commercial reuse, collection and commentary instead of “all legal content is copyright-free” |
| L3 | [Data Privacy Act, RA 10173](https://privacy.gov.ph/data-privacy-act-/) | Transparency, legitimate purpose and proportionality govern personal-data processing | Separate identity evidence, private work and public contribution records |
| L4 | [NPC Advisory Opinion 2020-032, 10 August 2020](https://privacy.gov.ph/wp-content/uploads/2020/11/Redacted-Advisory-Opinion-No.-2020-032.pdf) | Identifies blockchain controller-accountability and correction/erasure difficulties; recommends governance and privacy assessment in its registry context | The opinion is fact-specific, not approval for this service; define accountable operators and withdrawal handling before public immutable publication |
| L5 | [SC AI framework announcement, 19 March 2026](https://sc.judiciary.gov.ph/sc-adopts-framework-for-judiciarys-use-of-artificial-intelligence/) and [A.M. 25-11-28-SC text/index](https://elibrary.judiciary.gov.ph/thebookshelf/showdocs/11/101138) | Human judgment, accountability, fairness and transparency are relevant institutional requirements | Alignment is a design goal, not endorsement, procurement eligibility or permission to deploy in courts |
| L6 | [Rules on Electronic Evidence, A.M. 01-7-01-SC](https://lawphil.net/courts/supreme/am/am_01-7-01_sc_2001.html) | Electronic evidence has admissibility and authentication requirements | A hash or blockchain receipt is not automatic admissibility, notarization, or official certification; evaluate applicable later rules for any specific use |
| L7 | [SC: 5,572 new lawyers take oath, February 2026](https://sc.judiciary.gov.ph/5572-new-lawyers-take-their-oath-sign-roll-of-attorneys/) | A concrete recent entry cohort | Neither total active lawyers nor paying software customers; not a TAM estimate |

Some Judiciary pages returned HTTP 502 during direct retrieval. Indexed official text was available for L5; the E-Library and official administrative-matters index identify the resolution as **18 February 2026**, while one SC landing-page heading says 16 February. This audit does not resolve that metadata discrepancy or claim a complete clause-by-clause legal review. Obtain the authoritative instrument before implementation or a compliance claim.

## Engineering references

| ID | Primary source | Supported point |
|---|---|---|
| T1 | [RFC 8785: JSON Canonicalization Scheme](https://www.rfc-editor.org/rfc/rfc8785) | Deterministic JSON requires defined primitive serialization, UTF-16 property ordering and invalid-input handling; sorting Python dictionary keys alone is insufficient |
| T2 | [IPFS privacy and encryption](https://docs.ipfs.tech/concepts/privacy-and-encryption/) | Transport privacy does not make published content private; content visibility must be designed explicitly |
| T3 | [IPFS pinning](https://docs.ipfs.tech/how-to/pin-files/) | Persistence depends on retained/pinned copies; a CID is not a storage-duration guarantee |
| T4 | [CometBFT introduction, official documentation staging site](https://staging-cometbft-docs.netlify.app/staging/explanation/introduction/) | Safety assumes less than one-third Byzantine voting power; progress requires more than two-thirds; timing/network assumptions matter |
| T5 | [RFC 9162: Certificate Transparency v2](https://www.rfc-editor.org/rfc/rfc9162.html) | Established append-only Merkle-log, inclusion and consistency-proof patterns |

T5 is a reference pattern from certificate transparency, not a ready-made legal attestation protocol. Any adaptation still needs its own threat model, interoperability fixtures, keys, independent witnesses and operational review.

## Unknowns that research did not settle

- Current PWA revenue, churn, cohort retention, acquisition cost, inference costs and support burden.
- Actual customer preference for an auditable review record versus cheaper/general AI research.
- A complete Philippine legal-software market census or the uniqueness of this exact feature combination.
- Signed institutional commitments, a reviewer budget, and lawful access to a credential-verification service.
- An approved corpus license, privacy impact assessment, paid legal opinion or government grant award.
- Production deployment, consensus health, chain state, IPFS durability, or hosted PWA database authority.

All prices, headcounts, budgets and dates proposed in the roadmap are explicitly experiments or planning assumptions unless a source above supports them.
