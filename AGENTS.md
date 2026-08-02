# Agent contributor guide

This repository is building a security boundary, not only a developer tool. Treat
changes to execution, policy, artifacts, networking, and dependency handling as
security-sensitive.

## Authorized defensive scope

Capsule repository work is defensive, authorized, and local-only. Work must stay
within this repository; controlled local test harnesses, processes, and fixtures;
owned disposable guests; and isolated test environments explicitly authorized
by their owner and named in the task. “Local-only” means confined to those
assets, not permission to touch unrelated or third-party infrastructure.

No repository task—including one using offensive-security terminology—authorizes
targeting third-party systems, credential theft, persistence on unrelated hosts,
uncontrolled scanning, malware delivery, or bypassing platform, product,
organization, or repository safety or security controls. Stop rather than infer
authority beyond the named scope.

Adversarial fixtures, fault injection, boundary-abuse tests, and hostile-guest
tests are defensive validation of deny-by-default behavior. Keep them controlled,
reproducible, minimal, and confined to the exact authorized scope; do not
repurpose them into real-world intrusion tooling.

When drafting or delegating work, name the defensive intent, control under test,
exact fixtures or targets, and authorized environment. Recommended framing:

> Defensively validate `<Capsule control>` using `<specific repository fixture,
> local harness, or owned disposable guest>` in `<named authorized test
> environment>`. Do not access any other system, identity, credential, or data,
> and preserve Capsule’s existing safeguards.

## Start here

Read these documents before changing behavior:

1. `docs/PROJECT.md`
2. `docs/ARCHITECTURE.md`
3. `docs/TECHNICAL_DESIGN.md`
4. `docs/security/THREAT_MODEL.md`
5. `docs/FEASIBILITY_SPIKES.md` for pre-freeze or platform work
6. the relevant ADRs in `docs/adr/`

The JSON Schemas in `schemas/` are canonical for the current scaffold but are
explicitly pre-freeze. Do not extend their mixed `Job` authority model as a
shortcut. Follow `docs/protocol/OBJECT_MODEL.md` and the feasibility gates before
replacing them; keep current examples, TypeScript types, behavior, and schemas
consistent until that coordinated replacement.

## AI Central context

AI Central steering and reviewed skills are linked locally into `.codex/` and
ignored by Git. This file remains authoritative for Capsule-specific security
requirements. Generic steering or skill guidance must not weaken the
deny-by-default boundary, expand task scope, or override the verification
requirements below.

See `.codex/AI_CENTRAL.md` for the installed revision, selection, provenance,
licenses, and refresh workflow.

## Orchestrated task delivery

Treat the brain/orchestrator task as the owner of integration and delivery.
Before dispatch, it must choose and state whether work is a user-visible task,
an internal sub-agent assignment, or a research/experiment assignment.

- A new user-visible task that may retain repository changes should normally use
  its own `codex/<topic>` branch and worktree. When its work is complete, the task
  must verify and commit the result, push the branch, open a pull request, and
  report the branch, commit, PR, verification, limitations, and deferred work to
  the orchestrator. Do not leave completed retained work only on a local branch.
- Read-only coordination, monitoring, and research tasks need no branch or PR.
  If such a task begins retaining repository changes, convert it to the normal
  branch-and-PR workflow unless the orchestrator explicitly groups it elsewhere
  before the changes are made.
- Sub-agents spawned inside an orchestrator task are parallel workers, not
  independent delivery units. They must work in the orchestrator task's current
  branch and worktree. They must not create their own branch, commit, push, or PR
  unless the orchestrator explicitly reclassifies the assignment as a separate
  user-visible task. The orchestrator owns conflict resolution, final
  verification, commits, and the single shared PR.
- A research or experiment assignment running in another task must always call
  back to its parent/orchestrator before completing. Its handoff must include the
  question tested, defensive and authorized scope, method, results, retained
  evidence or artifact paths, verification commands, confidence and limitations,
  unresolved questions, and any recommended decision. If it created a branch,
  commit, or PR, include those identifiers too.
- The orchestrator must collect every child handoff and confirm that all retained
  commits have a PR or an explicitly documented integration destination before
  declaring the group complete. Research conclusions that affect architecture or
  security claims must be recorded in the repository's canonical documentation
  or ADRs; chat history alone is not retained evidence.

### Pull request publishing

- Write every multiline pull request description to a real temporary or other
  local body file with real newline bytes. Create or update the pull request with
  `gh pr create --body-file <path>` or `gh pr edit --body-file <path>`; never pass
  JSON-escaped multiline text through inline `--body` arguments.
- After every create or edit, read the stored description back with
  `gh pr view --json body --jq .body` and inspect the rendered Markdown. The
  publishing handoff fails if the body is missing its expected real newlines or
  contains literal `\n` or `\r` escape sequences where Markdown line breaks were
  intended. Fix the body with `gh pr edit --body-file <path>` and verify it again
  before reporting the PR.
- Choose draft or ready state deliberately. Pass `--draft` when creating a draft,
  omit it only when the PR is intentionally ready, and verify `isDraft` on
  readback. Change an existing PR with `gh pr ready` or `gh pr ready --undo`
  rather than assuming that editing its body changes review state.

Safe draft example (the body contains no secrets and the temporary file is
removed on shell exit):

```bash
pr_body_file="$(mktemp "${TMPDIR:-/tmp}/capsule-pr-body.XXXXXX")"
trap 'rm -f "$pr_body_file"' EXIT

cat >"$pr_body_file" <<'EOF'
## Summary

- Document the pull request publishing invariant.

## Verification

- `make ci`
EOF

gh pr create --draft --base main --head "$(git branch --show-current)" \
  --title "docs: harden pull request publishing" \
  --body-file "$pr_body_file"

pr_body="$(gh pr view --json body --jq .body)"
printf '%s\n' "$pr_body"
if [[ "$pr_body" != *$'\n'* || "$pr_body" == *'\n'* || "$pr_body" == *'\r'* ]]; then
  printf '%s\n' 'pull request body newline verification failed' >&2
  exit 1
fi
gh pr view --json url,isDraft --jq '{url, isDraft}'
```

## Working rules

- Preserve deny-by-default capabilities.
- Only the Execution Supervisor may authorize and own creation, termination,
  destruction, or reconciliation of a hostile guest. A narrowly enrolled helper
  may perform a required platform operation only from a sealed Supervisor
  descriptor. Do not add a daemon-to-backend or daemon-to-helper path.
- Execute by Supervisor-issued registration ID only. Never accept replacement
  plan bytes, backend flags, images, mounts, or guest paths at execute time.
- The daemon must not possess Approval, installation-root, or Supervisor evidence
  private keys; issue user-content authority; retrieve user-only content; reset
  grants; or clear quarantine, repair, or trust-epoch state.
- The Approval Broker renders Supervisor-registered typed plan data, not
  daemon-supplied display text, and approval remains one-use and attempt-bound.
- Treat device identifiers as identifiers, not trust. Trust comes from explicit
  local key authorization, purpose binding, installation identity, and trust
  epoch. DIDs never grant authority.
- Do not add live network DID resolution, arbitrary DID methods/resolvers, remote
  JSON-LD contexts, or full TUF/network parsing to approval or execution paths.
- Do not treat an in-process JavaScript sandbox as the host security boundary.
- Do not add unrestricted filesystem, network, process, environment, or artifact
  access to make an example easier.
- Do not pass live host paths into a guest. File inputs must become immutable,
  content-addressed snapshots before execution.
- Do not silently clamp user-owned resource limits. Resolve defaults before
  approval, reject limits above the user ceiling, and enforce the approved
  values exactly or refuse the job.
- Keep runtime adapters separate from the isolation backend.
- Keep rich document, archive, spreadsheet, PDF, image, media, and preview parsing
  out of the daemon and Execution Supervisor. Use bounded Broker validators or a
  future disposable parser sandbox.
- Treat spike code as non-production. Product packages must not import it; retain
  reproducible fixtures/evidence and record the resulting decision before reuse.
- Do not add a new Supervisor responsibility or privileged helper without an ADR.
- Record consequential architecture decisions in an ADR.
- Never claim a backend, profile, control, integrity mode, or security tier is
  implemented, validated, secure, continuous, attested, or production-ready unless
  its exact mechanism and retained adversarial evidence support that claim.

## Verification

Run before handing off a change:

```sh
pnpm install
pnpm check
pnpm lint
pnpm test
pnpm verify:schemas
go test ./...
go vet ./...
go build ./...
```

Use Node.js 22 or newer, pnpm 10, and Go 1.23 or newer for the current scaffold.
Runtime and toolchain pins are provisional until the first implementation ADR
locks them down.
