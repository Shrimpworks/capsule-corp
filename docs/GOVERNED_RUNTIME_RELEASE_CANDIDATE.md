# Governed `deno_core` Linux/arm64 release candidate

Date: 2026-08-04

Work item: unsigned runtime release-candidate consumption contract

Status: `PASSED`

Parent runtime selection: `IN_PROGRESS — TRENDING_GOOD`

Release publication: `BLOCKED`

Runtime/profile admission: `BLOCKED`; `RUNTIME-001` remains `unsupported`

## Result and boundary

The closed machine-readable candidate is
[`candidate-manifest.json`](../schemas/conformance/governed-deno-core-release-candidate/candidate-manifest.json).
Its self-digest is SHA-256
`78cf2e99e58a4e79413f22889dd19f794ac7cdce3e4ec5c167d6c2051d19afaa` under the
manifest's exact zeroed-field rule. C2A/C2B may name that candidate ID and digest as a passive
input. The reference grants no execution, guest, backend, signature, publication, selection, or
admission authority.

The manifest closes the exact consumption identity of the merged fork-native construction. It
does not copy the one-time experiment harness or large artifacts into this repository. Its public
evidence references are pinned to `Shrimpworks/capsule-experiments` merge
`fa03d7043b4f0653081d6c5733d597f49f6efd1c`, tree
`f80775335232ff4750f62998e5cc4d8e120ce90e`, exact relative paths, and SHA-256 digests.

## Exact governed sources and subjects

The governed source identities are:

- Deno upstream `14eea3160ae5834476aa3b9d317b8d41d991b982`, governed head
  `9adb0b68b55bca81644827f1e7749a3acb091bed`, and merge
  `ea18b9dc21ff8ebd19347be7095f47937ee14ec2`;
- `rusty_v8` upstream `d305e6afa7736f6e298c30ae6646f7709ee9382b`, governed head
  `80e863ddb942a4aa2b384e794fc23e35b9d2bb15`, and merge
  `cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`; and
- the separate C2 backend boundary at libkrun upstream
  `728df8125077d0db44265f6e997c72b81b65c015`, governed head
  `8a2c91943793668f31a1cf7af431933be935bb58`, merge
  `cf0333cdba478cc34a8570a65b38412da7fd3ecc`, and unchanged five-patch aggregate
  `d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`. Libkrun is not a
  component of this runtime candidate; these exact identities prevent C2 from silently treating a
  different backend line as already bound.

The eight C1 subjects are:

| Subject | Size | SHA-256 |
| --- | ---: | --- |
| raw `rusty_v8` archive | 184,789,162 | `e964d6b1b3689e91f8cf488d8a9f05764a03434b2e2e8347be5067300d39a7de` |
| gzip `rusty_v8` archive | 37,674,703 | `1ae209c9e4ba5803d010d2c79ee4cc0af0126c5a7ebcca211c7e41deaede4cd2` |
| `deno_core` binary | 68,497,624 | `56d3acefd2cc2f5136a0b8143c47131e49a58fbf66382dfd3e84f715ce8e2898` |
| startup snapshot | 699,988 | `4e8965217d5a6675a880326eee6f690bbeec7e7cb243decf2f3e9f453a871a2c` |
| two-file runtime bundle | 20,983,891 | `0cc08f93e82fcfe68b033e8807975a3bd67068a817da811a87a73aedc3f23937` |
| 22-entry root manifest | 1,807 | `100832dbb37737f29341bc5404df6d4405b8d6b706f274028892801fa88e7de8` |
| root tar | 71,895,040 | `9c46b45c4d220aedcc47c9ee53e875bc71d31d0b881b51740aaa9b882b5741e6` |
| root gzip | 22,192,615 | `e847651b35cd425dd8f6fe3bd45d693aff0af244e3a7bd30c629fa125cac62e8` |

The manifest also fixes Deno and `rusty_v8` source archives, Cargo locks and source closure,
builder images, compiler commits, LLVM, Clang, GN, Ninja, Rust target libraries, the 20-gitlink V8
source lock, all ten dynamic-root construction/evidence inputs, loader/library/source
correspondence, licenses/notices, CycloneDX and SPDX SBOMs, generated binding, build metadata, and
both unsigned provenance statements. A missing, extra, substituted, stale, movable, unknown,
reordered, over-cap, or mode-incompatible component refuses.

The runtime registry and final link contain exactly, and in this order:

1. `op_get_ext_import_meta_proto`
2. `op_get_extras_binding_object`
3. `op_set_captured_bootstrap`

The final-link statement digest is
`18793b423311e8e6a906cf4572d343ee33406c2c228343c42b6d85665c134429`.
No extension or module loader is present.

## Closed envelope and verifier

The consumption envelope has exactly 32 regular files: one manifest, eight runtime/root subjects,
13 supply-chain materials, and ten dynamic-root construction/evidence inputs. Envelope modes are
exactly `0644`. The inner root has exactly 22 entries: 11 directories, ten regular files, and one
symlink; only modes `0644` and `0755` occur. The manifest declares per-role byte caps and a 1 GiB
envelope cap. Array order is the declared role order; the inner-root manifest uses strictly
increasing unsigned UTF-8 path-byte order. No duplicate or unlisted path is permitted.

Run the repository-local verification with:

```sh
node scripts/verify-governed-runtime-release-candidate.mjs
node --test scripts/verify-governed-runtime-release-candidate.test.mjs
```

To recheck the retained public evidence and locally cloned fork metadata without network access:

```sh
node scripts/verify-governed-runtime-release-candidate.mjs \
  --evidence-root /path/to/capsule-experiments \
  --deno-repo /path/to/deno \
  --rusty-v8-repo /path/to/rusty_v8 \
  --libkrun-repo /path/to/libkrun
```

The bounded
[`mutation-corpus.json`](../schemas/conformance/governed-deno-core-release-candidate/mutation-corpus.json)
proves refusals for manifest-byte drift, unknown fields, movable and stale refs, ancestry
substitution, missing/extra/substituted components, cap/mode/order changes, evidence or C1 drift,
and invented signature/publication/admission state. This verifier reads public evidence or exact
local fixtures only. It does not download artifacts, parse large archives, execute a runtime,
start a guest, or contact a signing or release service.

## Governed release process

Publication remains blocked until a separately authorized release task completes all of these
steps on the exact candidate or creates a new candidate digest:

1. **Source and patch review.** Review every governed Deno commit from the upstream anchor, both
   patch/restoration oracles, every `rusty_v8` governed delta, the 20-gitlink source lock, the V8
   four-patch chain, effective generated build metadata, Cargo closure, dynamic-root package/source
   mapping, license/notice closure, SBOMs, and unsigned provenance. Unknown or unexplained bytes
   stop the release.
2. **Two-person approval target.** One runtime maintainer prepares the candidate; a second
   independent runtime/supply-chain reviewer approves the exact manifest digest, source trees,
   reproduction result, advisory posture, and limitations. The preparer cannot supply both
   approvals. Until the repository has an enforced two-person release rule and two recorded
   approvals, publication stays `BLOCKED`.
3. **CI reproduction.** Start from empty outputs/caches, fetch only digest-declared materials,
   disable the build/test/evidence network, use the pinned builders and toolchains, expose one
   snapshot CPU pinned to set 0, reproduce all eight subjects, then rerun fork verifiers, three-op
   registry/final link, source/license/SBOM/provenance closure, root/file-open/descriptor/syscall/
   restoration mutations, the offline verifier, and this repository's complete verification.
   Record CI identities and logs as new signed provenance inputs; the existing same-host and
   comparison-run evidence is not independent-builder evidence.
4. **Advisory ownership and SLA target.** Capsule runtime maintainers own Deno/`deno_core` and
   `rusty_v8`/V8 advisory monitoring; the release engineering owner owns builder, Cargo, Debian
   root, SBOM, provenance, and publication infrastructure. Target initial triage is one business
   day for critical/high advisories and five business days for other applicable advisories.
   Critical/high applicability targets a revoke-or-mitigate decision within 72 hours. These are
   release admission targets, not evidence that an operating service currently meets them.
5. **Signature and distribution.** Replace the manifest's absent artifact/provenance signature
   placeholders only in a new reviewed candidate, bind signer purpose and key authorization, add
   transparency/publication locations and revocation metadata, and verify readback. Linux runtime
   subjects are not Apple-notarized; final macOS distribution notarization remains a separate
   installed-profile placeholder. This task does not access keys or publish anything.
6. **Publication and readback.** Publish only the exact closed 32-file envelope, reject extras,
   download through the intended consumer path, and reverify size, mode, ordering, SHA-256,
   signatures, provenance, and manifest self-digest before marking publication `PASSED`.

## Revocation, supersession, and rollback

- Revocation identifies the candidate ID and manifest digest, reason, effective time, advisory or
  incident reference, and replacement if any. Revoked candidates refuse new selection and
  execution; revocation never rewrites the old manifest.
- Supersession creates a new candidate ID and self-digest with an explicit `supersedes` edge. A
  new release does not silently make old bytes acceptable.
- Rollback is allowed only to an independently unrevoked, unexpired, previously admitted manifest
  whose complete profile and trust epoch remain valid. “Last known bytes” or a movable tag is not
  rollback authority. If no such manifest exists, runtime selection and execution remain blocked.
- Partial publication, missing signatures, stale metadata, mixed candidate files, or readback
  mismatch revokes the incomplete publication attempt and retains the prior admission state; it
  never produces ordinary success.

## Mandatory reruns and remaining boundary

Any fork commit/tree/merge, patch/restoration oracle, V8/gitlink, toolchain, builder, Cargo lock or
source, generated metadata/binding, archive, runtime binary/snapshot, root package/source/member,
SBOM, notice, provenance, C1, loader/op/global/flag, signature, or publication-state change creates
a new manifest digest. It must repeat source review, two-person approval, clean CI reconstruction,
all fork and runtime/root mutations, offline closure verification, advisory review, and publication
readback. A backend/kernel/init/launcher/descriptor/resource change does not alter this runtime
candidate by itself, but invalidates C2 and every later composed-profile result and requires that
full composition corpus to rerun.

Construction of the exact unsigned subjects is `PASSED`. Governed release publication is
`BLOCKED` on independent review, enforced two-person approval, independent CI reproduction,
signatures, publication, and readback. Runtime selection remains `IN_PROGRESS — TRENDING_GOOD`.
Admission remains `BLOCKED` on release closure, C2 external-isolation composition, exact
libkrun/kernel/init/launcher/descriptors/resources, remaining P0 controls, installed distribution,
and a separate admission decision. `RUNTIME-001` remains `unsupported`.
