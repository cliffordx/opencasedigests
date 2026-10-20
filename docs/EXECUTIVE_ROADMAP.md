# Open Case Digests — evidence-gated roadmap

**Version 2 proposal · 20 October 2026 · Planning horizon: October 2026–October 2029**

This replaces the original fixed 12-month mainnet schedule. Dates are planning windows from an assumed 20 October 2026 start, not commitments. A gate that fails moves the dependent work; it does not lower the requirement. Estimates assume the staffing below. A founder-only effort will take longer or deliver less coverage.

The [independent audit](audits/2026-10-20-independent-audit.md) supplies the evidence and investment judgment. The [manifesto](MANIFESTO.md) defines the public purpose; [ADR-001](ARCHITECTURE_V2_PROPOSAL.md) proposes the architecture. These documents do not authorize spending, partnerships, production changes or a change to binding repository instructions.

## 1. The outcome to work toward

By the end of the first validated phase, a reader should be able to inspect a small public Philippine legal corpus, follow each reviewed proposition to a source passage, see exactly what was reviewed, identify corrections/disputes, and independently verify the published version. Case Digestor should demonstrate that this evidence helps a real study or professional task enough to support its ongoing cost.

The long-term outcome is a maintained, interoperable commons and a sustainable service business. Neither national completeness nor a blockchain mainnet is a prerequisite for success.

### Three separate scorecards

1. **Commons:** lawful reusable releases, accessible source evidence, useful corrections, independent reuse and continuity.
2. **Product:** reduced verification effort, recurring use, retained customers and willingness to pay for a specific workflow.
3. **Operations:** repeatable review quality, affordable maintenance, secure publication, recovery and accountable governance.

A successful score on one does not establish the others. Grant-funded public value can justify a small commons without proving a profitable SaaS; a profitable PWA does not justify calling its data open.

## 2. Program ownership and capacity

| Responsibility | Initial owner/effort assumption | Accountable output |
|---|---|---|
| Product and financing | Founder, 0.5–1 FTE | Scope, spending cap, PWA baseline, explicit go/hold/narrow decisions |
| Legal/editorial standard | Qualified Philippine legal lead, initially 0.25–0.5 FTE | Source/treatment rubric, adjudication, review sampling and legal claims |
| Software and integration | One experienced engineer, initially founder or 1 FTE | Isolated public-release prototype, verifier, tests and recovery |
| Review panel | 3–5 practitioners contributing scheduled paid/volunteer hours | Independent reviews, measured time, disagreements and calibration |
| Research curation | 1–2 part-time curators | Acquisition records, source locators, draft preparation and issue tracking |
| Privacy/rights/security | Named specialists on a scoped engagement | Rights inventory, privacy assessment, threat review and release findings |
| Institutional relationships | Founder with prospective faculty/library counterpart | Written pilot scope, budget owner and actual usage commitment |

An AI agent can assist retrieval, coding, comparison, drafting and test generation. It cannot supply legal admission, independent human editorial judgment, institutional consent or customer demand.

Do not call named universities, the IBP, BetterGov or the Supreme Court partners without agreement. Start with a faculty champion and a small practitioner group; validator recruitment follows editorial usefulness.

## 3. First 90 days: prove the smallest useful commons

### Days 0–14 · 20 October–2 November 2026 · Reset claims and define authority

**Owners:** founder + legal lead + engineer.

- Adopt or revise the proposed manifesto and architecture decision explicitly.
- Label the existing chain, browser helper and deployment kit as prototypes; preserve audit evidence.
- Open a remediation backlog for A01–A13. Decide which code remains experimental and which is necessary for the first public-release pilot.
- Choose one subject slice from actual PWA demand. Prefer 25–50 issuances with sufficient source access and several challenging treatment examples; do not choose solely for ease.
- Create a source-rights register, a public/private field inventory and a legal-review scope. Identify third-party annotations, contributor rights and restricted content.
- Define reviewed issuance identity, versioned digest/proposition identifiers and public-release authority. Reject docket-only identity.
- Define a provisional public content license and contribution agreement for counsel review; keep the existing MIT code license separate. No retroactive licensing of existing user work.
- Obtain a source snapshot of current PWA authority controls. Preserve annual gate-off, owner/CAS behavior, privacy, rewards and billing invariants.

**Gate G0:** named accountable owners; bounded corpus; explicit field/rights policy; no unsupported institution or readiness claims in the current entry documentation. Missing legal authority holds publication, but does not block synthetic/offline engineering.

### Days 15–30 · 3–18 November 2026 · Establish demand and a baseline

**Owners:** product lead + legal lead.

- Conduct 15–20 structured interviews: current PWA users, practicing lawyers, faculty/library staff and legal-aid users. Recruit outside the founder's closest supporters.
- Ask about the last real source-verification problem, current workaround, time/cost, errors and purchasing authority. Avoid asking only whether people like the idea of blockchain.
- Record current PWA usage/retention and cost baselines from an appropriately scoped owner-approved analytics report. No automatic export of personal data into the public repo.
- Build 100–200 bounded research tasks with answer rubrics: correct issuance, source support, opinion type, later treatment and uncertainty. Separate development cases from locked evaluation cases.
- Record current PWA/source-only results and a documented general-model retrieval baseline. If comparing paid competitors, use authorized accounts and record exact versions/date/scope.
- Calibrate reviewers on 10–15 examples, including disagreements. Define severity: fabricated/wrong authority, unsupported holding, wrong treatment, omission and presentation error.

Allocate the pilot's editorial-hour budget explicitly among corpus review, benchmark authoring, independent blind evaluation, calibration and incidents before G1. The illustrative 120–240 hours is a total planning envelope, not 120–240 hours for each activity; reduce corpus/task counts or increase funding if measured work does not fit.

**Gate G1:** at least one recurring job with evidence of pain; a buyer or sponsor hypothesis; evaluable tasks; reviewers who can explain their disagreements. If enthusiasm is abstract, narrow the problem before implementation grows.

### Days 31–60 · 19 November–18 December 2026 · Build and review the pilot

**Owners:** engineer + curators + legal lead.

- Implement an isolated approved-public-release schema and immutable revisions. Use newly approved sources; do not bulk export PWA personal digests.
- Establish exact-byte storage, versioned canonicalization, signing and a standalone verifier. Prefer an ordinary versioned artifact store as the initial baseline.
- Build a minimal review workflow with evidence passages, exact-version binding, credential checks, conflicts, independent reviews and adjudication. A one-click submission follows completed review.
- Develop a durable outbox and replay-safe delivery in a non-production environment. Use synthetic private canaries and no real client matters.
- Implement separate source, review, publication, integrity and later-treatment statuses. No aggregate “verified law” seal.
- Review the 25–50 issuances. Log real minutes and cost per accepted digest/proposition, not just number of clicks.
- Simulate changed source bytes, later resolution, separate opinion, misleading headnote, duplicate reviewer, forged signer, wrong revision, withdrawal and unavailable source.
- Have a second implementation or independent reviewer verify release fixtures and hash/signature failures.

**Gate G2:** every released candidate has a resolved issuance, source provenance, rights decision and version-bound review; all publication/privacy/security blocking tests pass. A single unresolved high-severity trust defect holds external reliance.

### Days 61–90 · 19 December 2026–17 January 2027 · Validate usefulness and cost

**Owners:** product + editorial + independent evaluator.

- Run a controlled pilot with roughly 20–40 users across the chosen segment(s), avoiding broad public claims.
- Compare baseline and evidence-backed workflows on matched tasks; randomize order where feasible. Record correctness, time, false confidence and the ability to locate the actual source.
- Seek 5–10 real paid pilot customers **or** a paid institutional pilot with a named budget owner. These are proposed decision thresholds, not promised results or a representative market sample.
- Predeclare eligible feature exposure and the existing-PWA comparison. Count attributable evidence-workflow use and incremental payment, renewal or retention improvement; existing subscribers who would have paid anyway do not establish a commons moat. Measure maintenance cost separately from the PWA's baseline costs.
- Ask an independent party to reconstruct at least one release and its corrections without PWA access.
- Test all queue failure points, lost responses, duplicate events, key rotation, withdrawal and restoration. Measure user-path overhead rather than asserting zero latency.
- Publish a pilot report with the denominator, exclusions, disagreements, costs and failures. Do not disclose private participants or source material without authority.

**Gate G3:** evidence of repeated useful use attributable to the new workflow; acceptable quality under the benchmark; a credible reviewer supply and funding model; no unresolved public-release blocker. Small pilot samples establish a direction, not a statistically proven retention lift. If willingness to pay is weak but public benefit is demonstrated, choose a smaller grant/sponsor-funded commons instead of disguising it as a proven SaaS moat.

## 4. Months 4–6 · January–April 2027 · One product workflow

**Dependency:** G3 passed or an explicitly funded narrower public-good plan.

**Goal:** connect the useful public evidence to one existing PWA task.

- Ship a feature-gated read path for reviewed public releases. Preserve source/review/coverage distinctions in the UI.
- Reuse existing PWA components and source-authority contracts. Define the exact public-release API; no broad PWA database reader.
- If production integration is pursued, complete its separate owner-run review, migration and release process. Source/local tests are not hosted authorization proof.
- Maintain a 100–200-digest scoped corpus only if the pilot's measured throughput supports it. Coverage count is secondary to maintained evidence quality.
- Add correction alerts, research-pack export or study feedback according to the validated job. Do not build all three simultaneously without capacity.
- Establish reviewer onboarding, scheduled compensation, rechecks, conflict reporting and backup adjudication.
- Secure one academic/clinic pilot with a concrete faculty workload, student consent/alternative path and renewal date. Names in a deck are not commitments.
- Publish schema, license inventory, verifier instructions, release snapshots and limitations. Ensure low-bandwidth and keyboard/screen-reader access.

**Gate G4:** 8 weeks of repeat use; at least one renewal/repeat-payment observation; positive contribution margin for the tested paid job or explicitly funded subsidy; sustainable review queue; correction and restoration drills pass.

**Hold/narrow:** if the service only sells while heavily subsidized, expose the subsidy and reduce scope. Do not add validators to fix an editorial or customer-retention problem.

## 5. Months 7–12 · April–October 2027 · Maintenance and distribution

**Dependency:** G4 and funded editorial capacity.

- Expand within the proven subject/domain before adding unrelated areas. Work toward 500 maintained digests only if the budget and hours support it.
- Add passage-level citation-treatment assertions with a clear reviewed universe and last-checked date. Preserve unresolved interpretations.
- Release a public benchmark subset and retain a separate rights-cleared holdout. Rotate compromised holdout items; never train and test on the same reviewed answers without disclosure.
- Test whether reviewed evidence improves retrieval, citation checking and Mock Bar evaluation. Do not label the system RLHF without an actual learning/evaluation pipeline.
- Introduce independent mirrors and signed checkpoints. Evaluate a witnessed transparency log using inclusion/consistency proofs and deliberate equivocation tests.
- Improve ingestion/source-change detection, source locators, metadata quality, monitoring, retrieval availability and complaint handling.
- Develop one institutional offering and one practitioner offering with separate costs and renewal measures. Avoid confusing free academic access with a committed paid deployment.
- Run external security review of the real public release path. Remediate findings before expanding exposure.

**Gate G5:** repeatable quality and maintenance economics across two review cycles; at least one independent consumer verifies/reuses releases; institutional or individual renewal evidence; a funded next-year plan.

**End-of-year decision:** continue as a commons plus SaaS, continue as a smaller sponsored commons, or stop expansion and preserve exports. A mainnet is not due merely because twelve months elapsed.

## 6. Months 13–18 · October 2027–April 2028 · Decide whether federation solves a real problem

**Dependency:** G5. This is an optional research lane, not the next required product release.

Interview independent institutions about a precise shared-control requirement: who must be able to publish, who must detect censorship/equivocation, and why signed witnessed releases are inadequate.

Compare:

1. Signed versioned releases and independent mirrors.
2. A witnessed append-only log.
3. A fully integrated Cosmos anchoring application.

Evaluate total operator cost, governance, incident recovery, key compromise, privacy handling, independent verification, and user benefit. Keep the PWA isolated in every option.

**Chain-specific work only if justified:** generated codecs; SDK module/ante wiring; credentials and authorization; immutable versioned state; deterministic full-state commitment; actual gas/size/rate limits; signed replay protection; persistent commit/recovery; event contract; checkpoint/backfill; reproducible builds and supply-chain review.

**Gate G6:** documented shared-writer need; written commitments from at least four genuinely independent appropriately weighted operators; named independent key custodians; governance charter; 12-month operating funds; a comparison demonstrating why simpler witnessing is inadequate. Four is a fault-tolerance planning minimum under equal voting power, not proof of decentralization.

If G6 fails, continue the simpler architecture and product roadmap. Do not repeatedly restart blockchain work to avoid accepting a negative result.

## 7. Months 19–24 · April–October 2028 · Conditional federation or stronger services

### If G6 passes

- Run a staging consortium with separate organizations, machines, credentials, administrators and documented operational duties.
- Obtain an independent security review covering authorization, evidence commitments, governance, economics of spam, dependencies and deployment exposure.
- Demonstrate one-validator outage, network partition, restart/replay, corrupted storage, checkpoint mismatch, lost signer, compromised reviewer, replacement validator and upgrade rollback procedures.
- Rebuild the public index from durable history and verify identical release records. Catch deliberate state divergence and equivocation.
- Retain private PWA data, billing and owner saves outside chain authority.
- Run for at least 90 days with measured availability and incident logs before a production anchoring decision. No “zero downtime” guarantee.

### If G6 does not pass

- Invest in editorial freshness, institutional renewal, accessibility, private workspaces, APIs and reliable signed exports.
- Add independent log witnesses or mirrors where useful without pretending they are validators.
- Measure the value of change alerts and citation-checking workflows against their costs.

**Gate G7:** production anchoring requires documented owner/operator approval after security and recovery acceptance. Service expansion requires renewal and positive economics. Both paths preserve public portability.

## 8. Months 25–36 · October 2028–October 2029 · Durable institution and bounded expansion

- Establish a legal/organizational arrangement appropriate to actual funding: company stewardship with independent editorial governance, a separate public-interest entity, or a consortium. Obtain advice; do not assume a foundation or DAO is necessary.
- Publish governance roles, conflicts, annual operating costs, public access commitments and the process for changing editorial policy.
- Make successor operation practical: periodic exports, reproducible verifier, key succession, lawful private-data return/deletion and a funded continuity plan.
- Expand subjects based on demand and available reviewer expertise. Publish coverage maps instead of claiming complete national coverage.
- Consider legislation/administrative issuances only with separate identity, update and rights models.
- Investigate regional interoperability only after a real partner and a concrete exchange need emerge. IBC is not synonymous with semantic interoperability.
- Treat notarial registries, property diligence, title-chain verification and financial escrow as separate business/legal proposals. No automatic expansion from a working digest commons.

**Gate G8:** at least two annual cycles of accountable maintenance or a smaller sustainable scope; repeat customers/sponsors; independent reuse; functioning correction/appeal process; continuity not dependent on one person.

### Years 4–5, if justified

A national legal commons with multiple maintainers, stable public standards and services used across institutions is a worthy ambition. Cross-border legal graphs and registries are options, not promises. Re-run market, rights, authority and technical diligence for each new domain.

## 9. Quality and product acceptance measures

These are **proposed pilot thresholds**, to be calibrated and fixed before evaluation. They are not observed results or guarantees of legal safety.

| Measure | Proposed gate | How to measure |
|---|---|---|
| Release completeness | 100% of pilot public releases have required source/version/rights/review fields | Independent inventory with missing/duplicate/wrong-version adversarial fixtures |
| Source binding | No wrong-issuance release in the evaluated corpus | Two-person identity review, including docket collisions and later resolutions |
| Critical errors | Zero fabricated authority or materially false holding/treatment in release sample | Independent blinded review; a discovery triggers hold and root-cause work |
| Citation support | At least 98% supported claims in a predeclared sample | Exact passage review; report n and confidence interval, not just percentage |
| Treatment detection | Report precision and recall on a bounded known-label benchmark | Include negative/ambiguous examples; do not generalize recall to the whole legal corpus |
| Verification time | Target 25% median reduction against baseline without reduced accuracy | Matched tasks; include full source checking and assistance time |
| Reviewer independence | Two eligible independent reviews for pilot “reviewed” designation | No self-review; conflict check; duplicate identity and stale-version tests |
| Repeat use | Target at least 40% of activated pilot users returning to the evidence workflow in week four | Define exposure/activation, compare existing-PWA baseline, exclude one-time demos and unrelated usage |
| Payment | 5–10 actual paying users or a paid institution, followed by repeat payment/renewal attributable to the tested offering | Collected revenue and workflow use; distinguish incremental value from pre-existing subscription revenue |
| Maintenance | High-severity corrections triaged within 1 working day; public update target 3 working days | Named duty owner and incident log; complex unresolved questions stay visibly disputed |
| Public availability | Pilot target 99.5% artifact retrieval availability over 30 days | External checks; count source unavailable and mirror failure distinctly |
| Outbox recovery | Zero lost releases or duplicate economic effects in injected failures | Crash/retry at every state boundary; deterministic reconciliation |
| PWA overhead | Provisional p95 added synchronous overhead ≤50ms under agreed load | Compare same environment/load, include DB outbox work; no external pin/chain wait |
| Economics | Positive paid-job contribution or a documented funded subsidy | Include recurring review, support, inference and maintenance costs |

Zero observed errors is not zero risk. With no errors in 100 independent trials, the rough one-sided 95% upper bound remains about 3%; correlated legal propositions make simple binomial interpretations weaker. Report sample composition, reviewer disagreement and uncertainty. Do not turn a pilot score into “98% legally correct” marketing.

## 10. Test strategy and release evidence

### Unit and schema

Cross-language canonical bytes; malformed/duplicate JSON keys; bounded numbers and Unicode; source/digest/review revision binding; signature domain separation; score bounds; credential expiry/revocation; legal status independence; license/field allowlist.

### Integration

Actual public-release transaction to durable outbox; worker to storage; real emitted event to index; failed transaction ignored; replay/backfill; source update; public quarantine; access isolation; PWA no-credit/no-reward/no-owner-mutation invariants. For any Cosmos lane, exercise real signed transactions and enforced grants rather than handwritten JSON substitutes.

### End-to-end

A user locates the source, reads review scope, reports a correction and obtains the corrected version. A reviewer signs the intended version; another device/retry cannot silently change it. A third party verifies the export without logging into the PWA. Test mobile, keyboard, screen reader and low-bandwidth behavior when UI is implemented.

### Security and operations

Threat model; private-data canaries; hostile strings/CIDs; quotas; upstream fetch limits; admin API isolation; key custody; scoped credentials; dependency inventory; monitoring; restore and deletion/withdrawal drills. Production evidence is a distinct gate and must identify environment, versions, owner and time.

### Evidence bundle per gate

Store: commit IDs, test commands/results, fixtures and expected outputs, independent review findings, benchmark protocol/results, corpus denominator/coverage, costs, named accountable decision maker, open risks and the resulting go/hold/narrow decision. Publish only public-safe portions.

## 11. Budget, capacity and opportunity cost

Illustrative 90-day pilot model in PHP; **assumed rates for planning, not salary estimates, supplier quotations or spending authorization**:

| Item | Assumption | Scenario cost |
|---|---|---:|
| Legal review/curation | 120–240 hours × PHP600 | 72,000–144,000 |
| Engineering support beyond founder | 120 hours × PHP1,000 | 120,000 |
| Privacy/rights/security advice | 40 hours × PHP2,000 | 80,000 |
| Hosting/tools/contingent compute | 3 months × PHP10,000 | 30,000 |
| Subtotal | Before contingency | 302,000–374,000 |
| Contingency | 20% | 60,400–74,800 |
| Total scenario | Excludes founder opportunity cost, taxes and unpriced changes | **362,400–448,800** |

Founder time is not free economically. At 20 hours/week for 13 weeks, reserve **260 hours** and record which PWA improvements those hours displace. If the budget is smaller, reduce corpus/reviewer scope and use existing infrastructure; do not remove identity, privacy or review gates.

For ongoing planning, estimate accepted and maintained items from actual throughput:

```text
maintenance hours = released items × measured monthly recheck/correction rate × hours per recheck
new-item hours/month = max(0, productive hours/month - maintenance hours/month - calibration and incident hours/month)
accepted new items/month = new-item hours/month / measured hours per accepted item
```

Maintain separate budgets for the commons, PWA feature development and optional protocol research. Proposed first-year allocation after validation: at least half the combined effort on editorial quality and useful product workflows; cap speculative chain work at 10–15% until G6 evidence exists. These allocations are recommendations to expose tradeoffs, not financial forecasts.

Do not assume grants will cover operations. Research specific DOST/DTI/DICT or civic funding calls only when applying; eligibility, award and disbursement each need evidence. A tokenless design is not an award criterion by itself.

## 12. Roadmap risk register

| Risk | Owner | Trigger | Required response |
|---|---|---|---|
| Rights uncertainty | Legal lead | Unclear source/contributor reuse authority | Hold affected publication; use approved alternative |
| Private data leak | Privacy + engineering | Canary or real private field reaches export | Stop release worker; incident response and scoped remediation |
| Editorial harm | Editorial lead | Material wrong holding/treatment | Quarantine claim, notify affected users as appropriate, independent review |
| Reviewer scarcity | Editorial lead | Queue exceeds two planned review cycles | Reduce coverage, fund review, stop new ingestion |
| Weak demand | Product owner | Repeat use/payment gates fail | Narrow job or choose sponsored public-good model |
| Cost overrun | Founder | Forecast exceeds funded runway/cap | Stop optional expansion; preserve core operation |
| Governance capture | Independent editorial participants | Sponsor demands favorable review or founder suppresses dispute | Conflict disclosure, independent appeal and transparent decision |
| Infrastructure distraction | Founder | Chain consumes effort while G3/G4 fail | Freeze protocol expansion |
| Nominal federation | Operators | Shared keys/admin/host undermine independence | Correct claims; redesign custody before federation gate |
| Maintainer departure | Operations owner | Critical function has one trained person | Document/run restore and handoff, fund continuity |

## 13. First implementation backlog after the strategy decision

The audit intentionally did not implement these changes. Order dependencies as follows:

1. Public-release/source/review identity specification and rights policy.
2. Regression fixtures for unsafe publication, wrong issuance and stale review binding.
3. Canonical payload and verifier with independent test vectors.
4. Small review workflow and benchmark corpus.
5. Durable publication outbox, scoped reader and replay-safe artifact delivery.
6. PWA read-only evidence display for approved releases, behind a gate.
7. Correction/withdrawal, monitoring, public export and independent reuse.
8. Paid pilot and maintenance economics.
9. Only if justified: witnessed log, then isolated corrected Cosmos anchoring.

Within the existing prototype, A01–A08 block trust claims. Do not cosmetically fix the badge or event name and call the protocol secure. A chain implementation needs the complete authorization, state commitment, persistence and integration path proved together.

The founder's next investment decision is therefore concrete: commit to the bounded public-release/editorial pilot and its measured gates, or reduce the scope. There is no need to promise years of infrastructure in order to begin useful work.
