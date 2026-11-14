# Canonical Go patch prerequisite

Status: PASSED for exact Go1.26.11 local validation. Fresh independent review1/3 returned Ready with no actionable findings;
hosted acceptance remains a final publication step. Full command results and limits are retained in
[canonical verification](canonical-verification.md) and [validation ledger](validation.json).

The reviewed #387 head had green CI from 2026-10-15. Its normal merge on
2026-11-13 America/Toronto triggered a new advisory database evaluation:
[CI43533822264](https://github.com/Shrimpworks/capsule-corp/actions/runs/43533822264)
failed on GO-2026-6859 (net/http, fixed1.26.11) and GO-2026-6809
(crypto/x509, fixed1.26.10) in the pinned Go1.26.9 standard library.
This is new dated baseline evidence, not an archive parser regression.

## Decision and dependency checklist

Use Go1.26.11, the smallest same-minor patch closing both observed scanner findings.
[Official release history](https://go.dev/doc/devel/release#go1.26.11) records the
2026-11-10 release, including net/http security fixes. The
[x509 advisory](https://pkg.go.dev/vuln/GO-2026-6809) describes its affected range.
Keep the exact scanner1.6.0, current advisory evaluation, all CI checks and branch
protections. Do not select the host's Go1.27.2 or rewrite historical evidence.

Applicable ecosystem rows: GO-3 (standard library) and TEST-1 (quality/security tools).
Trust: existing compiler and standard-library supply chain, with no new module,
product authority, algorithm, privilege or network path. Go/BSD provenance and
existing dependency graph remain applicable. Upgrade/removal owner: Capsule maintainer.
Reproduction uses exact downloaded toolchain, not the host default. Existing full
checks and current vulnerability scan are the acceptance tests; no new application
behavior test is warranted for a version-only patch. If replacement/removal becomes
necessary, retain the security floor and requalify the alternate toolchain; do not
restore a vulnerable version to make a check pass.

Active changed requirements: go.mod, AGENTS.md, README.md, docs/DEVELOPMENT.md,
and docs/DEV_MACHINE_SETUP_AND_RECOVERY.md. CI already reads go.mod. Existing
canonical test expectations derive that pin, requiring no hardcoded test update.
The frozen C5b19 Gate2 Go1.25.13 pin is not advanced; it is retained solely for
this offline inert conformance work, not certified as a current secure runtime.

## Toolchain provenance

Downloaded through Go's normal verified toolchain mechanism:
`GOTOOLCHAIN=go1.26.11 go version`.

- Module: `golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64`
- Cached module zip checksum: `h1:cGxzEKi4mhhQFd6Cz/2NSL4Y9nfIBCirzU1yn2rnOY8=`
- LICENSE SHA256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`
- bin/go SHA256: `20e07ffb4f08faaf125e413db024b6e988d3ec2799dcb3568a720d14be0ad133`

Archive toolchain, frozen host/ABI packet, product/runtime admission and child-run
authorization remain unchanged. Canonical scanner success is dated source analysis,
not absence of all possible vulnerabilities or product admission.
