# Development machine setup and recovery

Date: 2026-08-05

Status: canonical onboarding runbook. Part A (toolchain setup) is usable today by any contributor.
Part B (Apple identity/credential setup) sequences the existing canonical Apple docs into an
execution checklist; it inherits their `BLOCKED` items and authorizes nothing by itself. Part C
(machine-loss recovery) is unrehearsed guidance — it has not been exercised end to end on a real
replacement machine, so treat it as a drill script to validate, not proven evidence. This document
changes no product behavior, signing state, credential, key, or portal resource.

## Purpose

Capsule's setup knowledge is otherwise scattered across [Development](DEVELOPMENT.md) (language
toolchains), [Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md)
(Apple identity/credential inventory and policy), and
[Apple Development provisioning plan](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md) (App ID/profile
proposal). Each is correct and remains canonical for its subject. None of them is sequenced as "do
this, then this" for a fresh machine, a new contributor, or a lost/replaced Mac. This document is
that sequence. It does not restate policy or rationale already covered by the linked docs — follow
the links for the *why*.

Two audiences:

- **A new contributor** who only needs Part A (repository toolchain). Most contributions need
  nothing from Part B.
- **Dylan, on a new or wiped Mac**, or anyone auditing continuity risk, who needs Part A, Part B,
  and Part C.

## Part A — Repository toolchain (every contributor)

Do this on any Mac before touching the repository.

1. Install Xcode and accept its license. `xcode-select --install` alone is not sufficient for the
   Apple-identity work in Part B, which needs the full Xcode app (Settings > Accounts, Certificates
   pane). No Xcode version is currently pinned anywhere in this repository or its CI — record the
   exact version you install (`xcodebuild -version`) in your own notes; there is no product Xcode
   project yet to pin a version against (see [Open gaps](#open-gaps) below).
2. Install Go, matching `go.mod` (currently 1.23 or newer).
3. Install Node.js, matching `.node-version` (currently 22.22.1).
4. Enable pnpm via Corepack, matching `package.json`'s `packageManager` pin (currently 10.28.2):
   ```sh
   corepack enable
   ```
5. Install Bun (currently pinned to 1.3.14 in [Development](DEVELOPMENT.md#prerequisites)) only if
   you are working on runtime-profile experiments. It is not an admitted workload profile and
   ordinary contributions do not need it.
6. Clone the repository, then install JS/TS dependencies:
   ```sh
   pnpm install --frozen-lockfile
   ```
7. Create the ignored local AI Central symlinks if you use that steering content:
   ```sh
   pnpm codex:links -- --dry-run   # preview
   pnpm codex:links                # apply
   ```
   Set `AI_CENTRAL_HOME` first if the AI Central repository is not checked out at `../ai-central`.
8. Verify the toolchain is correctly pinned:
   ```sh
   make ci
   ```
   A clean `make ci` run is the acceptance bar for Part A. If it fails, fix the toolchain before
   writing code — do not chase test failures caused by a wrong Go/Node/pnpm version.
9. Rust is needed only if you work on the Source Validator artifacts under
   `artifacts/mjs-source-validator-*`. There is currently no `Cargo.toml` or `rust-toolchain.toml`
   committed to this repository — those artifacts are disposable spikes retained in the separate
   [`Shrimpworks/capsule-experiments`](https://github.com/Shrimpworks/capsule-experiments) archive,
   per the [Rust engineering standards](RUST_ENGINEERING_STANDARDS.md#toolchain-and-dependencies)
   and the [feasibility spike workflow](DEVELOPMENT.md#feasibility-spike-workflow). Pin the exact
   toolchain (`rust-toolchain.toml`) inside that experiment when you start one; there is nothing to
   install in this repository itself yet.
10. Read [CONTRIBUTING.md](../CONTRIBUTING.md)'s reading list before making a behavioral change.

Nothing in Part A touches an Apple account, certificate, or Keychain identity.

## Part B — Apple identity and credential setup

Only needed for work that requires a signed, installed, or provisioned build (installed-test rows
in the [provisioning plan](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md), Source Validator R3, or any
Developer ID/notarization task). This section is a checklist that sequences existing canonical
content — every step links to the doc that defines it in full; do not skip to a step without
reading its linked section first.

1. Confirm Apple Developer Program membership and Team ID via
   [§1 "Confirm membership and team in the portal"](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#1-confirm-membership-and-team-in-the-portal).
   Current selected team: `3DDR84M4JS` (Individual). Do not trust a certificate's parenthesized
   display suffix as the Team ID — see the
   [Team ID reconciliation](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#team-id-reconciled-exact-role-inputs-still-stop-credentialed-work).
2. Freeze bundle topology and identifiers before creating anything in the portal — follow
   [§2 "Freeze identifiers and capabilities"](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#2-freeze-identifiers-and-capabilities)
   together with the
   [component/bundle-ID proposal](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md#component-inventory-and-bundle-identifier-proposal).
   Option A vs. Option B (single main app vs. separate installer) must be decided here, not later —
   it changes how many App IDs you register.
3. Obtain or verify the Apple Development signing identity —
   [§3 "Obtain the Apple Development identity"](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#3-obtain-the-apple-development-identity).
   Verify subject OU and signed `TeamIdentifier` both equal the selected Team ID before using any
   identity found already present in Keychain Access; do not select by common name.
4. Register the test Mac and create per-App-ID development profiles —
   [§4 "Register the test Mac and create development profiles"](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#4-register-the-test-mac-and-create-development-profiles).
   Inspect every downloaded profile's metadata with the documented read-only commands before use.
5. Configure Xcode targets and author entitlements per role — Part A's Xcode install plus
   [§5 "Configure Xcode targets"](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#5-configure-xcode-targets)
   and the [minimum/prohibited entitlements tables](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md#minimum-entitlements-per-role).
6. Only when distribution work actually begins: Developer ID identity
   ([§6](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#6-obtain-developer-id-only-when-distribution-begins))
   and notarization credentials
   ([§7](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#7-choose-notarization-authentication)). Do not
   pull these steps forward — development-signed local testing needs neither.
7. Before starting a credentialed task, satisfy every item in the
   [replacement-input checklist](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#exact-decision-and-replacement-inputs)
   and the
   [verification checklist](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#verification-checklist-for-a-credentialed-task).
   Record results in the redacted inventory (see [below](#redacted-credential-inventory)), never as
   secret bytes in Git.

Never paste a private key, `.p12`, `.p8`, password, or token into a command, terminal transcript,
issue, PR, or shell history — see
[Setup without exposing secrets](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#setup-without-exposing-secrets).

## Part C — Machine-loss / reconstruction drill

Scenario: the Mac holding the Apple Development identity, Developer ID identity, and any local
Capsule operational keys is lost, wiped, or stolen. This section separates what is recoverable from
what must be recreated, and orders the recreation. It restates the
[rotation/revocation/recovery rules](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#rotation-revocation-expiry-recovery-update-and-rollback)
as an execution sequence; that section remains the authority on *why* each rule exists.

### Step 0 — Contain

1. Assume every private key that lived only in that Mac's login Keychain is gone, not merely
   inaccessible. Apple Development and Developer ID private keys are normally exportable but are
   only actually recoverable if a prior encrypted export exists in
   [approved offline custody](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#storage-and-access-policy).
   Nonexportable Secure Enclave-backed Capsule operational keys (installation-root, Approval,
   Supervisor evidence) cannot be recovered under any circumstance — plan for reissue, not restore.
2. If loss is by theft or compromise rather than ordinary hardware failure, treat it as the
   [suspected-compromise path](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#developer-id-certificate)
   for any Developer ID identity that existed on the machine: stop releases, revoke with Apple, and
   rotate notarization credentials independently — do not simply reissue and continue.
3. Check whether an encrypted backup of the old login Keychain exists under the documented custody
   procedure. If none exists (expected under the current per-machine setup — see
   [Open gaps](#open-gaps)), skip straight to reissue.

### Step 1 — Rebuild the toolchain

Repeat [Part A](#part-a-repository-toolchain-every-contributor) in full on the replacement Mac.
Nothing here is Apple-identity-specific.

### Step 2 — Reissue Apple Development identity and profiles

1. Sign in to the Apple Developer portal on the new Mac and re-verify Team `3DDR84M4JS` membership
   ([Part B step 1](#part-b-apple-identity-and-credential-setup)).
2. Create a **new** Apple Development certificate on the new Mac
   ([Part B step 3](#part-b-apple-identity-and-credential-setup)). This is a fresh private key —
   there is no way to recreate the old one's bytes. Verify subject OU and `TeamIdentifier` before
   using it.
3. Regenerate every development profile that referenced the old certificate's fingerprint
   ([rotation guidance](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#apple-development-certificate)).
   Register the new Mac's device UDID if profile-based device registration is in use.
4. Revoke the old certificate only after confirming no other authorized host still depends on it.

### Step 3 — Reissue Developer ID and notarization credentials (only if distribution work exists)

1. Create a new Developer ID Application certificate on the new Mac if the lost machine held the
   distribution identity; do not reuse or attempt to export the old one.
2. Re-run `xcrun notarytool store-credentials` interactively on the new Mac to recreate the saved
   Keychain profile; a previous saved profile does not survive the machine loss
   ([§7](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#7-choose-notarization-authentication)).
3. If a team App Store Connect API key was in use, revoke the old key in App Store Connect and
   generate a replacement; its `.p8` is only downloadable once, so the original cannot be recovered
   ([notarization credential rotation](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#notarization-credentials)).

### Step 4 — Reissue Capsule operational keys, if any were ever installed

Capsule's own installation-root, Approval, and Supervisor evidence keys are not implemented in
product code today (see the
[credentials-versus-protocol-keys table](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#apple-credentials-versus-capsule-protocol-keys)),
so as of this writing there is nothing installed to lose. When that changes, this step must be
updated to point at the exact authorized forward-epoch trust-transition ceremony instead of a
restore — Capsule protocol keys are never recovered from backup by design.

### Step 5 — Update the redacted inventory

Record every reissued certificate/profile in the redacted inventory
([template](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#redacted-public-metadata-inventory-template))
with status `replaced` for every retired row and `active` for every new one. An inventory update is
not optional — it is the only artifact that lets a future reconstruction (or this same drill,
rehearsed) start from known state instead of rediscovering it.

## Redacted credential inventory

The [Apple certificates doc](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#redacted-public-metadata-inventory-template)
defines the exact table shape (public metadata only — fingerprints, serials, UUIDs, expirations;
never private keys, `.p12`/`.p8` contents, or passwords) but does not name where the filled table
lives. **That location is not yet decided.** Until it is, the inventory does not exist anywhere
outside ad hoc chat/terminal history, which is not retained evidence. Pick one and record the
decision here:

- [ ] Password manager (e.g., a 1Password vault/note dedicated to Capsule Apple credentials)
- [ ] The private `Shrimpworks/capsule-experiments` archive already used for spike/evidence custody
- [ ] A new dedicated private repository
- [ ] Something else

This is the single highest-priority open item in this document — every other step in Part B and
Part C assumes the inventory has somewhere to go.

## Open gaps

Tracked here until closed, so this audit is not lost to chat history:

1. **Inventory location undecided** — see [above](#redacted-credential-inventory).
2. **No Xcode/Swift toolchain pin** anywhere in the repository (no `.xcode-version`, no CI job that
   builds native/Swift code — CI currently runs only `ubuntu-latest` jobs). There is also no product
   Xcode project yet to pin a version against; add the pin when one exists.
3. **No committed Rust toolchain** — `rust-toolchain.toml` is documented as required
   ([Rust engineering standards](RUST_ENGINEERING_STANDARDS.md#toolchain-and-dependencies)) but
   nothing in this repository currently pins one, because no Rust source is committed here yet
   (spikes only, in `capsule-experiments`).
4. **No rehearsed machine-loss drill.** Part C above is written but has never been executed against
   a real replacement Mac. Treat it as a script to test, not proven recovery evidence, until someone
   runs it (ideally deliberately, not during an actual loss) and records the result.
5. **App IDs/profiles proposed, not created.** The [provisioning plan](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md)
   is still a proposal — none of its explicit App IDs exist in the Apple portal yet. Part B above is
   therefore unexercised as a full sequence; the first person to run it end to end should correct
   this document with what actually happened, not just what was planned.
6. **Individual-membership single point of failure.** Only Dylan, as Account Holder, can create
   Developer ID certificates or manage Apple Developer team resources
   ([role limits](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md#account-membership-team-and-roles)). If
   continuity across multiple maintainers becomes a requirement, that is an account-type decision
   (individual → organization membership) outside this document's scope, not something Part C's
   drill can work around.
