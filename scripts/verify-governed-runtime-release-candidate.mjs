import { execFileSync } from "node:child_process";
import { createHash } from "node:crypto";
import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

const manifestUrl = new URL(
  "../schemas/conformance/governed-deno-core-release-candidate/candidate-manifest.json",
  import.meta.url,
);
const c1Url = new URL(
  "../schemas/conformance/c1-governed-deno-core/controlled-development-profile.json",
  import.meta.url,
);

export const EXPECTED_SELF_DIGEST =
  "78cf2e99e58a4e79413f22889dd19f794ac7cdce3e4ec5c167d6c2051d19afaa";
const ZERO_DIGEST = "0".repeat(64);

const expectedTopLevelKeys = [
  "objectType",
  "schemaVersion",
  "candidateId",
  "selfDigest",
  "status",
  "target",
  "governedSources",
  "construction",
  "runtimeSubjects",
  "runtimeAuthority",
  "dynamicRoot",
  "supplyChainMaterials",
  "evidence",
  "compatibility",
  "closurePolicy",
  "releaseState",
  "limitations",
  "invalidationTriggers",
];

const expectedSubjects = Object.freeze({
  "rusty-v8-static-archive-raw": [
    184_789_162,
    "e964d6b1b3689e91f8cf488d8a9f05764a03434b2e2e8347be5067300d39a7de",
  ],
  "rusty-v8-static-archive-gzip": [
    37_674_703,
    "1ae209c9e4ba5803d010d2c79ee4cc0af0126c5a7ebcca211c7e41deaede4cd2",
  ],
  "deno-core-binary": [
    68_497_624,
    "56d3acefd2cc2f5136a0b8143c47131e49a58fbf66382dfd3e84f715ce8e2898",
  ],
  "startup-snapshot": [699_988, "4e8965217d5a6675a880326eee6f690bbeec7e7cb243decf2f3e9f453a871a2c"],
  "two-file-runtime-bundle": [
    20_983_891,
    "0cc08f93e82fcfe68b033e8807975a3bd67068a817da811a87a73aedc3f23937",
  ],
  "dynamic-root-manifest": [
    1_807,
    "100832dbb37737f29341bc5404df6d4405b8d6b706f274028892801fa88e7de8",
  ],
  "dynamic-root-tar": [
    71_895_040,
    "9c46b45c4d220aedcc47c9ee53e875bc71d31d0b881b51740aaa9b882b5741e6",
  ],
  "dynamic-root-gzip": [
    22_192_615,
    "e847651b35cd425dd8f6fe3bd45d693aff0af244e3a7bd30c629fa125cac62e8",
  ],
});

export class CandidateRefusal extends Error {
  constructor(code, message) {
    super(`${code}: ${message}`);
    this.name = "CandidateRefusal";
    this.code = code;
  }
}

function refuse(condition, code, message) {
  if (!condition) throw new CandidateRefusal(code, message);
}

function sha256(bytes) {
  return createHash("sha256").update(bytes).digest("hex");
}

function unique(values) {
  return new Set(values).size === values.length;
}

function exactJson(value) {
  return JSON.stringify(value);
}

export function verifyManifestBytes(bytes) {
  refuse(bytes.length <= 32_768, "RC_COMPONENT_SIZE_MISMATCH", "manifest exceeds 32 KiB cap");
  const text = bytes.toString("utf8");
  let manifest;
  try {
    manifest = JSON.parse(text);
  } catch (error) {
    throw new CandidateRefusal("RC_MANIFEST_DIGEST_MISMATCH", `invalid JSON: ${error.message}`);
  }

  const digest = manifest?.selfDigest?.sha256;
  const needle = `    "sha256": "${digest}"\n  },\n  "status"`;
  refuse(
    typeof digest === "string" && text.split(needle).length === 2,
    "RC_MANIFEST_DIGEST_MISMATCH",
    "self-digest field is missing, duplicated, or not in the exact encoding",
  );
  const normalized = Buffer.from(
    text.replace(needle, `    "sha256": "${ZERO_DIGEST}"\n  },\n  "status"`),
  );
  refuse(sha256(normalized) === digest, "RC_MANIFEST_DIGEST_MISMATCH", "self-digest mismatch");
  refuse(digest === EXPECTED_SELF_DIGEST, "RC_MANIFEST_DIGEST_MISMATCH", "unknown manifest bytes");
  verifyParsedManifest(manifest);
  const envelopeBytes = [
    ...manifest.runtimeSubjects,
    ...manifest.supplyChainMaterials,
    ...manifest.dynamicRoot.sourceCorrespondence,
  ].reduce((total, entry) => total + entry.size, bytes.length);
  refuse(
    envelopeBytes <= manifest.closurePolicy.byteCaps.candidateEnvelopeTotal,
    "RC_COMPONENT_SIZE_MISMATCH",
    "candidate envelope exceeds aggregate cap",
  );
  return manifest;
}

export function verifyParsedManifest(manifest) {
  refuse(
    exactJson(Object.keys(manifest)) === exactJson(expectedTopLevelKeys),
    "RC_MANIFEST_UNKNOWN_FIELD",
    "top-level fields are missing, extra, or reordered",
  );
  refuse(
    manifest.objectType === "capsule.governed-deno-core.release-candidate-consumption-manifest" &&
      manifest.schemaVersion === 1 &&
      manifest.candidateId === "governed-deno-core-linux-arm64-rc-2026-08-04-1",
    "RC_COMPONENT_SUBSTITUTED",
    "candidate identity mismatch",
  );

  const sourceText = JSON.stringify(manifest.governedSources);
  refuse(
    !/(refs\/heads\/|\/tree\/(main|master)(?:\b|\/)|@(main|master|latest)\b)/.test(sourceText),
    "RC_MOVABLE_REFERENCE",
    "movable source reference present",
  );
  const deno = manifest.governedSources.deno;
  const rustyV8 = manifest.governedSources.rustyV8;
  refuse(
    rustyV8.governedHead === "80e863ddb942a4aa2b384e794fc23e35b9d2bb15" &&
      !rustyV8.rejectedTerminalHeads.includes(rustyV8.governedHead),
    "RC_STALE_SOURCE_HEAD",
    "stale or unknown rusty_v8 head",
  );
  refuse(
    exactJson(deno.mergeParents) ===
      exactJson([
        "14eea3160ae5834476aa3b9d317b8d41d991b982",
        "9adb0b68b55bca81644827f1e7749a3acb091bed",
      ]) &&
      deno.governedTree === "72edd0f7b5f83b918945860653714e344c8a303f" &&
      exactJson(rustyV8.mergeParents) ===
        exactJson([
          "eddede228a9214c4dfb6a85aeca22abc0679100d",
          "80e863ddb942a4aa2b384e794fc23e35b9d2bb15",
        ]) &&
      rustyV8.governedTree === "d8950a7a1ee907761720b23d24eaa9b63aa33b10",
    "RC_SOURCE_ANCESTRY_MISMATCH",
    "merge parents or governed tree mismatch",
  );

  refuse(manifest.runtimeSubjects.length >= 8, "RC_COMPONENT_MISSING", "runtime subject missing");
  refuse(manifest.runtimeSubjects.length <= 8, "RC_COMPONENT_EXTRA", "runtime subject extra");
  refuse(
    manifest.supplyChainMaterials.length === 13 &&
      manifest.dynamicRoot.sourceCorrespondence.length === 10,
    manifest.supplyChainMaterials.length < 13 ||
      manifest.dynamicRoot.sourceCorrespondence.length < 10
      ? "RC_COMPONENT_MISSING"
      : "RC_COMPONENT_EXTRA",
    "supply-chain closure count mismatch",
  );

  const subjects = Object.fromEntries(
    manifest.runtimeSubjects.map((entry) => [entry.role, [entry.size, entry.sha256]]),
  );
  refuse(
    exactJson(subjects) === exactJson(expectedSubjects),
    "RC_COMPONENT_SUBSTITUTED",
    "runtime subject identity mismatch",
  );

  const caps = manifest.closurePolicy.byteCaps;
  const capByRole = {
    "rusty-v8-static-archive-raw": caps.rustyV8RawArchive,
    "rusty-v8-static-archive-gzip": caps.rustyV8GzipArchive,
    "deno-core-binary": caps.denoBinary,
    "startup-snapshot": caps.denoSnapshot,
    "two-file-runtime-bundle": caps.denoTwoFileBundle,
    "dynamic-root-manifest": caps.dynamicRootManifest,
    "dynamic-root-tar": caps.dynamicRootTar,
    "dynamic-root-gzip": caps.dynamicRootGzip,
  };
  refuse(
    manifest.runtimeSubjects.every((entry) => entry.size <= capByRole[entry.role]) &&
      manifest.supplyChainMaterials.every(
        (entry) => entry.size <= caps.individualSupplyChainMaterial,
      ) &&
      manifest.dynamicRoot.sourceCorrespondence.every(
        (entry) => entry.size <= caps.individualDynamicRootInput,
      ),
    "RC_COMPONENT_SIZE_MISMATCH",
    "component exceeds declared cap",
  );

  refuse(
    exactJson(manifest.closurePolicy.allowedEnvelopeModes) === exactJson(["0644"]) &&
      exactJson(manifest.closurePolicy.innerRootAllowedModes) === exactJson(["0644", "0755"]),
    "RC_COMPONENT_MODE_MISMATCH",
    "mode vocabulary mismatch",
  );
  refuse(
    exactJson(manifest.runtimeAuthority.builtinOps) ===
      exactJson([
        "op_get_ext_import_meta_proto",
        "op_get_extras_binding_object",
        "op_set_captured_bootstrap",
      ]) &&
      manifest.runtimeAuthority.finalLinkStatementSha256 ===
        "18793b423311e8e6a906cf4572d343ee33406c2c228343c42b6d85665c134429" &&
      manifest.runtimeAuthority.extensionRegistry.length === 0 &&
      manifest.runtimeAuthority.moduleLoader === "none",
    "RC_COMPONENT_ORDER_MISMATCH",
    "three-op registry/final-link statement mismatch",
  );

  const allPaths = [
    "candidate-manifest.json",
    ...manifest.runtimeSubjects.map((entry) => entry.path),
    ...manifest.supplyChainMaterials.map((entry) => entry.path),
    ...manifest.dynamicRoot.sourceCorrespondence.map((entry) => entry.path),
  ];
  refuse(
    allPaths.length === 32 && unique(allPaths),
    allPaths.length < 32 ? "RC_COMPONENT_MISSING" : "RC_COMPONENT_EXTRA",
    "candidate envelope is not the exact unique 32-file closure",
  );

  refuse(
    manifest.evidence.mergeCommit === "fa03d7043b4f0653081d6c5733d597f49f6efd1c" &&
      manifest.evidence.mergeTree === "f80775335232ff4750f62998e5cc4d8e120ce90e" &&
      manifest.evidence.files.find(
        (entry) => entry.path === "evidence/2026-08-04/ref-verification.json",
      )?.sha256 === "7cb086ce168e56bbd9b9f7caa415fb8ba5f55abf0ad103c2b020cfa015badb89",
    "RC_EVIDENCE_MISMATCH",
    "retained evidence identity mismatch",
  );
  refuse(
    manifest.compatibility.c1FixtureSize === 9_289 &&
      manifest.compatibility.c1FixtureSha256 ===
        "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e" &&
      manifest.compatibility.executionEligible === false &&
      manifest.compatibility.admissionEligible === false,
    "RC_C1_MISMATCH",
    "C1 compatibility mismatch",
  );
  const libkrun = manifest.compatibility.c2Dependencies.libkrun;
  refuse(
    libkrun.upstreamCommit === "728df8125077d0db44265f6e997c72b81b65c015" &&
      libkrun.governedHead === "8a2c91943793668f31a1cf7af431933be935bb58" &&
      libkrun.governedTree === "ffa4131ddcc6ec66edd623381dae94189ccd3fee" &&
      libkrun.mergeCommit === "cf0333cdba478cc34a8570a65b38412da7fd3ecc" &&
      exactJson(libkrun.mergeParents) ===
        exactJson([
          "4ea8d1de861ed1c0636fc800b6da8fb71a086aa5",
          "8a2c91943793668f31a1cf7af431933be935bb58",
        ]) &&
      libkrun.includedInThisCandidate === false,
    "RC_SOURCE_ANCESTRY_MISMATCH",
    "libkrun C2 boundary identity mismatch",
  );
  refuse(
    manifest.releaseState.artifactSignature === "absent-required-before-publication" &&
      manifest.releaseState.provenanceSignature === "absent-required-before-publication",
    "RC_SIGNATURE_STATE_MISMATCH",
    "signature placeholder state mismatch",
  );
  refuse(
    manifest.releaseState.publication === "blocked-not-published",
    "RC_PUBLICATION_STATE_MISMATCH",
    "publication state mismatch",
  );
  refuse(
    manifest.releaseState.runtimeSelected === false &&
      manifest.releaseState.runtimeAdmitted === false &&
      manifest.status.runtime001 === "unsupported",
    "RC_ADMISSION_STATE_MISMATCH",
    "selection/admission state mismatch",
  );
  return manifest;
}

export async function verifyC1(manifest, path = fileURLToPath(c1Url)) {
  const bytes = await readFile(path);
  refuse(
    bytes.length === manifest.compatibility.c1FixtureSize,
    "RC_C1_MISMATCH",
    "C1 size mismatch",
  );
  refuse(
    sha256(bytes) === manifest.compatibility.c1FixtureSha256,
    "RC_C1_MISMATCH",
    "C1 digest mismatch",
  );
  const c1 = JSON.parse(bytes);
  const c1Digests = Object.values(c1.artifacts).map((entry) => entry.sha256);
  const subjectDigests = manifest.runtimeSubjects.map((entry) => entry.sha256);
  refuse(
    exactJson(c1Digests) === exactJson(subjectDigests),
    "RC_C1_MISMATCH",
    "C1 subjects differ",
  );
}

export async function verifyEvidence(manifest, evidenceRoot) {
  const base = join(evidenceRoot, manifest.evidence.basePath);
  for (const entry of manifest.evidence.files) {
    let bytes;
    try {
      bytes = await readFile(join(base, entry.path));
    } catch {
      throw new CandidateRefusal("RC_COMPONENT_MISSING", `missing evidence ${entry.path}`);
    }
    refuse(
      sha256(bytes) === entry.sha256,
      "RC_EVIDENCE_MISMATCH",
      `evidence differs: ${entry.path}`,
    );
  }
}

function git(repo, args) {
  return execFileSync("git", ["-C", repo, ...args], { encoding: "utf8" }).trim();
}

export function verifyGitSources(manifest, denoRepo, rustyV8Repo, libkrunRepo) {
  for (const [name, source, repo] of [
    ["deno", manifest.governedSources.deno, denoRepo],
    ["rusty_v8", manifest.governedSources.rustyV8, rustyV8Repo],
    ["libkrun", manifest.compatibility.c2Dependencies.libkrun, libkrunRepo],
  ]) {
    refuse(
      git(repo, ["show", "-s", "--format=%T", source.governedHead]) === source.governedTree,
      "RC_SOURCE_ANCESTRY_MISMATCH",
      `${name} governed tree mismatch`,
    );
    refuse(
      git(repo, ["show", "-s", "--format=%P", source.mergeCommit]) ===
        source.mergeParents.join(" "),
      "RC_SOURCE_ANCESTRY_MISMATCH",
      `${name} merge parents mismatch`,
    );
    try {
      execFileSync("git", [
        "-C",
        repo,
        "merge-base",
        "--is-ancestor",
        source.upstreamCommit,
        source.governedHead,
      ]);
    } catch {
      throw new CandidateRefusal(
        "RC_SOURCE_ANCESTRY_MISMATCH",
        `${name} upstream is not an ancestor`,
      );
    }
  }
}

function parseArgs(argv) {
  const result = {};
  for (let index = 0; index < argv.length; index += 2) {
    const key = argv[index];
    const value = argv[index + 1];
    if (!key?.startsWith("--") || value === undefined) throw new Error(`invalid argument: ${key}`);
    result[key.slice(2)] = value;
  }
  return result;
}

export async function run(argv = []) {
  const args = parseArgs(argv);
  const manifestPath = resolve(args.manifest ?? fileURLToPath(manifestUrl));
  const manifest = verifyManifestBytes(await readFile(manifestPath));
  await verifyC1(manifest, resolve(args.c1 ?? fileURLToPath(c1Url)));
  if (args["evidence-root"]) await verifyEvidence(manifest, resolve(args["evidence-root"]));
  if (args["deno-repo"] || args["rusty-v8-repo"] || args["libkrun-repo"]) {
    refuse(
      args["deno-repo"] && args["rusty-v8-repo"] && args["libkrun-repo"],
      "RC_COMPONENT_MISSING",
      "all three governed fork repositories are required",
    );
    verifyGitSources(
      manifest,
      resolve(args["deno-repo"]),
      resolve(args["rusty-v8-repo"]),
      resolve(args["libkrun-repo"]),
    );
  }
  process.stdout.write(
    `${manifest.candidateId} ${manifest.selfDigest.sha256} PASSED offline-verification-only\n`,
  );
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  await run(process.argv.slice(2));
}
