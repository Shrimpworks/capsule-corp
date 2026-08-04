# Authenticated local IPC S1 consistency stop

Date: 2026-08-03

Status: retained historical TypeScript stop; first-release dependency superseded by accepted
ADR-0034. The architecture owner originally selected ADR-0030's versioned atomic-cutover path, and
no S1 message, bridge, reply, refusal, or field-authority fixture was generated. That stop still
governs any later TypeScript cutover, but it no longer blocks the `.mjs` first-release S1 path.

## Defensive scope

This review defensively checked only passive repository contracts and retained local fixtures. It
created no XPC listener, endpoint, key, identity, content, runtime, backend, process, guest, or
deployment and accessed no external system, credential, identity, or data.

## Question and decision

Can S1 retain one exact 562-byte complete role-binding record while also including the three
distinct TypeScript approved-byte source roles merged with ADR-0030?

No. The architecture owner selected ADR-0030's versioned atomic cutover. S1 fixture construction
stays stopped until the dependency chain below produces a finalized plan, complete role model, and
newly versioned registration method and binding record. Encoding a fixture now would either omit
approved-byte authority, conflate nominal roles, or reinterpret a v0 method as v1. Each option is
an authority change, not a fixture-only implementation choice.

Accepted ADR-0034 later chose a different first-release premise: one byte-exact pass-through
`main.mjs` member has one authoring/executable source role and therefore fits plan v0 without the
three TypeScript roles. The current 562-byte complete role projection can now be retained for that
first-release contract, while exact source-manifest and source bytes are separately validated
atomic registration inputs. The revised passive S1 work is defined in ADR-0034 and the
authenticated local IPC plan. Nothing here turns the 626-byte TypeScript arithmetic into a layout.

## Exact conflict

ADR-0029 defines the 562-byte `RegisterPlanV0` projection as:

| Field group | Bytes |
| --- | ---: |
| Record version | 1 |
| Installation ID | 16 |
| Epoch, source manifest, inline input, and runtime bundle digests | 128 |
| Review count | 1 |
| Eight fixed review-digest slots | 256 |
| Profile registry, backend validation, backend configuration, trust snapshot, and policy digests | 160 |
| **Total** | **562** |

The merged approved-byte contract and ADR-0030 require a future plan to own three distinct source
roles:

1. original-authoring source-manifest digest;
2. executable-JavaScript source-manifest digest; and
3. ordered transformation-record-set digest.

The current v0 projection has only one source-manifest digest. Replacing that one role with all
three adds two 32-byte fields: `562 + 64 = 626` bytes. This is only size arithmetic, not a proposed
wire layout or known answer. Field order, record version, method version, caps, and classifications
remain undecided here.

ADR-0030 explicitly makes the current 562-byte record v0-only and requires the approved-byte
migration to version the registration method and complete binding record together. The cutover
plan likewise forbids freezing or reinterpreting 562 bytes for the three new plan roles. The
current TypeScript and Go registration/lifecycle bindings still contain the single v0
`sourceManifestDigest`; the approved-byte roles remain separate passive Slice A fixtures.

`RegisterPlanV0` and its 562-byte projection therefore remain historical/current-plan-v0 design
only. They will not be frozen into an S1 corpus. The 626-byte arithmetic above is not an approved
layout, cap, record, or known answer.

## Retained baselines and counts

This stop preserves the existing generated corpora unchanged:

- main conformance corpus: 82 rules, 262 cases, and 368 fixture files, excluding its manifest and
  manifest schema;
- TypeScript approved-byte corpus: 9 known-answer files and 14 refusal mutations; and
- S1 artifacts added by this task: 0 fixture files, 0 cases, and 0 byte known answers.

Relevant retained known answers are:

| Artifact | Bytes | SHA-256 |
| --- | ---: | --- |
| ExecutionPlan v0 | 530 | `627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c` |
| PlanRegistration v0 | 165 | `f3569d37ad6d787c2cdd575ef9ec6c369bbe495157c43110fc9e9d610a277614` |
| Candidate approval envelope | 375 | `fb0a9e7c983f6f3986260dce857edf6b18cba99ee386f9532300dbdc31a5a3bd` |
| Original-authoring manifest | 171 | `1010ae00c786a6266348173c7760e0190be4cc280be1f71c8549f09727e4b183` |
| Executable-JavaScript manifest | 174 | `295138062d0785785373b8c468fee75f77a28131d0974f30f69c4050425e9814` |
| Transformation-record set | 714 | `5738283a5accdbd8b736af81982dc46068172ec502f5c43e4113fe7de10c76eb` |

There is no retained 562-byte S1 record and no 626-byte successor known answer. Producing either
would violate the selected atomic-cutover path.

## Retained TypeScript integration path

S1 may resume only after this dependency order completes:

1. accept the transformation-owner and immutable source-store topology;
2. finalize `ExecutionPlan` v1 and its complete nominal role model;
3. merge and integrate the separately developed canonical field-authority manifest;
4. define a newly versioned registration method and binding record with reviewed field order,
   exact caps, cap-plus-one behavior, and cross-language known answers; and
5. explicitly review the other registration fetch, approval submission, and attempt-request
   method projections and version each one whose typed shape changes.

This sequence now applies only if the conditional TypeScript feature resumes. First-release S1
instead follows ADR-0034's M1/M2 single-member `.mjs` contract and keeps one active plan-v0 shape.

There is no dual active v0/v1 acceptance, optional transformation role, generic fallback, or field
inference. This task did not duplicate the field-authority manifest, infer its field ownership, or
create a parallel manifest. S1 classification must consume that canonical manifest only after it
merges and after the newly versioned method defines the fields to classify.

## Deferred boundaries

- S1 common-header, four-request, four-success-reply, fixed-refusal, maxima/cap-plus-one,
  copy-ownership, response-loss, idempotency, and state/effect oracle fixtures remain uncreated.
- S2 remains blocked on the shared newly versioned S1 bytes and finalized complete-role contract.
  It must not define the missing record or field authority inside a Go facade.
- S3 remains blocked on the same shared fixture contract. It must not use native message parsing to
  create a de facto field or version decision.

The existing S2 prohibition on product IPC and the existing S3 no-product local harness boundary
remain unchanged.

## Confidence and limitations

Confidence is high because the size result follows directly from fixed-width fields and ADR-0030
states the versioning rule explicitly. This record does not decide the eventual method name,
record layout, field authority, exact cap, known answer, or other method-version outcomes; those
belong to the selected dependency chain.
