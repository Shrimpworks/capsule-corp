import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  createTrustedProfileRegistryContext,
  createTrustedUserPolicyContext,
  type DecodedJobProposalCandidate,
  decodeJobProposal,
  resolveJobProposal,
  type TrustedProfileRegistryContext,
  type TrustedUserPolicyContext,
} from "./index.js";

interface ConformanceFixture {
  readonly path: string;
}

interface ProposalResolutionOracle {
  readonly sourceManifest: ConformanceFixture | null;
  readonly sourceManifestDigest: string | null;
  readonly canonicalInlineInput: ConformanceFixture | null;
  readonly inlineInputDigest: string | null;
  readonly wallTime: {
    readonly milliseconds: number;
    readonly origin: "requested" | "trusted-default";
  } | null;
}

interface ConformanceCase {
  readonly id: string;
  readonly object: string;
  readonly fixture: ConformanceFixture;
  readonly context: {
    readonly kind: string;
    readonly oracle?: ProposalResolutionOracle;
  };
  readonly expected: {
    readonly decision: "accept" | "reject";
    readonly classification: string | null;
    readonly owner: string;
  };
  readonly implementations: { readonly typescript: string };
}

interface ConformanceManifest {
  readonly cases: readonly ConformanceCase[];
}

interface ProfileRegistryFixture {
  readonly profiles: readonly {
    readonly alias: string;
    readonly status: "active" | "inactive";
    readonly exactWallTimeMs: { readonly minimum: number; readonly maximum: number };
  }[];
}

interface UserPolicyFixture {
  readonly wallTimeMs: { readonly trustedDefault: number; readonly ceiling: number };
}

const corpusRoot = new URL("../../../schemas/conformance/v0/", import.meta.url);
const manifest = JSON.parse(
  await readFile(new URL("manifest.json", corpusRoot), "utf8"),
) as ConformanceManifest;
const semanticCases = manifest.cases.filter(
  (entry) =>
    entry.object === "JobProposal" && entry.expected.owner === "job-proposal-semantic-validator",
);
const registryFixture = JSON.parse(
  await readFile(new URL("contexts/profile-registry.json", corpusRoot), "utf8"),
) as ProfileRegistryFixture;
const policyFixture = JSON.parse(
  await readFile(new URL("contexts/user-policy.json", corpusRoot), "utf8"),
) as UserPolicyFixture;

const profileRegistry = createTrustedProfileRegistryContext(registryFixture.profiles);
const userPolicy = createTrustedUserPolicyContext(policyFixture);

test("Task 3B marks exactly its semantic-resolution manifest slice verified", () => {
  assert.equal(semanticCases.length, 15);
  assert.equal(
    semanticCases.every((entry) => entry.implementations.typescript === "verified"),
    true,
  );
  assert.equal(new Set(semanticCases.map((entry) => entry.fixture.path)).size, 12);
});

for (const entry of semanticCases) {
  test(`Task 3B conformance: ${entry.id}`, async () => {
    const bytes = new Uint8Array(await readFile(new URL(entry.fixture.path, corpusRoot)));
    const decoded = decodeJobProposal(bytes);
    assert.equal(decoded.ok, true, `${entry.id} must cross only the Task 3A boundary`);
    if (!decoded.ok) {
      return;
    }

    const result = resolveJobProposal(decoded.proposal, profileRegistry, userPolicy);
    assert.equal(result.ok, entry.expected.decision === "accept", entry.id);
    if (!result.ok) {
      assert.deepEqual(
        [result.refusal.owner, result.refusal.classification],
        [entry.expected.owner, entry.expected.classification],
      );
      return;
    }

    assert.equal(entry.context.kind, "proposal-resolution");
    if (entry.context.kind !== "proposal-resolution" || !entry.context.oracle) {
      return;
    }
    await assertOracle(result.resolved, entry.context.oracle);
  });
}

test("resolver retains sorted exact source and canonical input known answers", async () => {
  const resolved = await resolveFixture("job-proposal/ordinary.json");
  assert.deepEqual(
    resolved.source.files.map((file) => file.path),
    ["main.mjs"],
  );
  assert.equal(resolved.source.aggregateByteLength, 50);
  assert.equal(resolved.source.manifest.exactBytes.byteLength, 89);
  assert.equal(
    resolved.source.manifest.digest,
    "c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14",
  );
  assert.equal(resolved.inlineInput.exactBytes.byteLength, 118);
  assert.equal(
    resolved.inlineInput.digest,
    "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72",
  );
  assert.deepEqual(resolved.wallTime, { milliseconds: 5_000, origin: "requested" });
  assert.deepEqual(resolved.planInputs, {
    sourceManifestDigest: resolved.source.manifest.digest,
    sourceEntrypoint: "main.mjs",
    sourceByteLength: 50,
    inputSlot: "primary-data",
    inlineInputDigest: resolved.inlineInput.digest,
    inlineInputByteLength: 118,
    runtimeProfileAlias: "fixture-active@1",
    wallTimeMs: 5_000,
    wallTimeOrigin: "requested",
    outputSlot: "transformed-json",
    outputMaxJsonBytes: 65_536,
  });
});

test("resolver and trusted context constructors return deeply immutable passive values", async () => {
  const resolved = await resolveFixture("job-proposal/ordinary.json");
  assert.equal(Object.isFrozen(profileRegistry), true);
  assert.equal(Object.isFrozen(profileRegistry.profiles), true);
  assert.equal(Object.isFrozen(profileRegistry.profiles[0] as object), true);
  assert.equal(Object.isFrozen(userPolicy.wallTimeMs), true);
  assert.equal(Object.isFrozen(resolved), true);
  assert.equal(Object.isFrozen(resolved.source), true);
  assert.equal(Object.isFrozen(resolved.source.files), true);
  const firstSourceFile = resolved.source.files[0];
  assert.ok(firstSourceFile);
  assert.equal(Object.isFrozen(firstSourceFile.utf8Bytes.bytes), true);
  assert.equal(Object.isFrozen(resolved.source.manifest.exactBytes.bytes), true);
  assert.equal(Object.isFrozen(resolved.inlineInput.exactBytes.bytes), true);
  assert.equal(Object.isFrozen(resolved.runtimeProfile), true);
  assert.equal(Object.isFrozen(resolved.wallTime), true);
  assert.equal(Object.isFrozen(resolved.planInputs), true);

  const firstByte = resolved.inlineInput.exactBytes.bytes[0];
  assert.equal(Reflect.set(resolved.inlineInput.exactBytes.bytes, "0", 0), false);
  assert.equal(resolved.inlineInput.exactBytes.bytes[0], firstByte);
});

test("resolver rejects values not issued by its decoder or trusted context constructors", async () => {
  const decoded = decodeJobProposal(
    new Uint8Array(await readFile(new URL("job-proposal/ordinary.json", corpusRoot))),
  );
  assert.equal(decoded.ok, true);
  if (!decoded.ok) {
    return;
  }

  const copiedProposal = { ...decoded.proposal } as DecodedJobProposalCandidate;
  assert.deepEqual(resolveJobProposal(copiedProposal, profileRegistry, userPolicy), {
    ok: false,
    refusal: {
      owner: "job-proposal-semantic-validator",
      classification: "BINDING",
      code: "CANDIDATE_PROVENANCE",
    },
  });

  const copiedRegistry = {
    profiles: profileRegistry.profiles,
  } as TrustedProfileRegistryContext;
  assert.deepEqual(resolveJobProposal(decoded.proposal, copiedRegistry, userPolicy), {
    ok: false,
    refusal: {
      owner: "job-proposal-semantic-validator",
      classification: "BINDING",
      code: "TRUSTED_CONTEXT_PROVENANCE",
    },
  });

  const copiedPolicy = { wallTimeMs: userPolicy.wallTimeMs } as TrustedUserPolicyContext;
  assert.deepEqual(resolveJobProposal(decoded.proposal, profileRegistry, copiedPolicy), {
    ok: false,
    refusal: {
      owner: "job-proposal-semantic-validator",
      classification: "BINDING",
      code: "TRUSTED_CONTEXT_PROVENANCE",
    },
  });
});

async function resolveFixture(path: string) {
  const decoded = decodeJobProposal(new Uint8Array(await readFile(new URL(path, corpusRoot))));
  assert.equal(decoded.ok, true);
  if (!decoded.ok) {
    throw new Error("fixture did not cross the raw/schema boundary");
  }
  const result = resolveJobProposal(decoded.proposal, profileRegistry, userPolicy);
  assert.equal(result.ok, true);
  if (!result.ok) {
    throw new Error("fixture did not cross the semantic boundary");
  }
  return result.resolved;
}

async function assertOracle(
  resolved: Awaited<ReturnType<typeof resolveFixture>>,
  oracle: ProposalResolutionOracle,
): Promise<void> {
  if (oracle.sourceManifest) {
    assert.deepEqual(
      Uint8Array.from(resolved.source.manifest.exactBytes.bytes),
      new Uint8Array(await readFile(new URL(oracle.sourceManifest.path, corpusRoot))),
    );
    assert.equal(resolved.source.manifest.digest, oracle.sourceManifestDigest);
  }
  if (oracle.canonicalInlineInput) {
    assert.deepEqual(
      Uint8Array.from(resolved.inlineInput.exactBytes.bytes),
      new Uint8Array(await readFile(new URL(oracle.canonicalInlineInput.path, corpusRoot))),
    );
    assert.equal(resolved.inlineInput.digest, oracle.inlineInputDigest);
  }
  if (oracle.wallTime) {
    assert.deepEqual(resolved.wallTime, oracle.wallTime);
  }
}
