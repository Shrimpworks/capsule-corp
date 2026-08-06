# Issue label taxonomy

Labels classify independent dimensions. Prefer one work type, one narrow component, and only the
workflow, platform, or status labels that add useful information.

## Work type and concern

| Label | Use |
| --- | --- |
| `bug` | Reproducible non-security incorrect behavior. |
| `proposal` | A capability or direction awaiting a maintainer decision. |
| `enhancement` | An accepted improvement ready for implementation. |
| `implementation` | A bounded implementation slice with acceptance criteria. |
| `research` | A falsifiable research question or controlled experiment. |
| `architecture` | An ADR, trust-boundary, topology, or authority decision. |
| `documentation` | Documentation or contributor guidance. |
| `security-boundary` | Public design, hardening, or verification work that affects a security boundary. Never use it for a vulnerability disclosure. |
| `evidence` | Retained conformance, adversarial, reproduction, or audit evidence. |

## Workflow

| Label | Use |
| --- | --- |
| `needs:triage` | Needs initial scope, ownership, and categorization. |
| `needs:decision` | Needs an explicit maintainer or ADR decision. |
| `needs:evidence` | Needs exact retained evidence before its claim can advance. |
| `help wanted` | Maintainers welcome contributor help. |
| `good first issue` | Small, well-scoped entry task with maintainer guidance. |
| `question` | Needs clarification or maintainer guidance. |
| `duplicate` | Duplicates an existing issue or pull request. |
| `invalid` | Malformed, unsupported, or outside repository scope. |
| `wontfix` | Closed without planned work. This is not canonical `NO_GO` unless the exact candidate is explicitly rejected. |

## Component

- `component:daemon`
- `component:broker`
- `component:supervisor`
- `component:protocol-ipc`
- `component:storage`
- `component:source-validation`
- `component:runtime-backend`
- `component:guest`
- `component:install-update`
- `component:docs-governance`
- `component:tooling`

## Platform

- `platform:macos`
- `platform:linux`
- `platform:cross-platform`

## Canonical work status

Status labels mirror [`docs/STATUS_LANGUAGE.md`](../docs/STATUS_LANGUAGE.md). They describe only
the exact labeled work item, never ADR lifecycle, control-evidence maturity, or product admission.

| Label | Meaning |
| --- | --- |
| `status:passed` | Exact scoped item met every declared acceptance condition. |
| `status:in-progress-trending-good` | Active work; latest evidence is closing risk or blockers. |
| `status:in-progress-trending-bad` | Active work; latest evidence found new difficulty, while the path remains open. |
| `status:blocked` | Intended path awaits a named dependency, decision, artifact, credential, contract, or owner. |
| `status:no-go` | Exact candidate is abandoned and has a documented replacement or fallback. |

Close completed issues when practical. Status labels are useful when an issue is retained as a
dashboard or evidence record; remove stale workflow and status labels deliberately.

Suspected vulnerabilities do not receive a public issue label. Use
[private vulnerability reporting](https://github.com/Shrimpworks/capsule-corp/security/advisories/new).
