import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  APPROVED_BYTE_CAPS,
  type ApprovedByteFixtureBuildRefusalCode,
  type ApprovedByteFixtureCandidate,
  asApprovedByteDigest,
  buildApprovedByteFixtureCandidate,
  verifyApprovedByteFixtureKnownAnswers,
} from "./typescript-approved-byte-candidate.js";

const fixtureRoot = new URL(
  "../../../schemas/conformance/typescript-approved-byte-v0/",
  import.meta.url,
);
const knownAnswers = {
  transformerProfile: "ab82cd553d52490fb0e2cc2e6cfc8cc106440a601c99c70cebad43537efc7f97",
  normalizedOptions: "8a43f134b568e983a7e4e24a763209a61405dc535a710e4799611033e4983b2e",
  originalManifest: "1010ae00c786a6266348173c7760e0190be4cc280be1f71c8549f09727e4b183",
  executableManifest: "295138062d0785785373b8c468fee75f77a28131d0974f30f69c4050425e9814",
  transformationRecords: ["3ddf8b242afbddaa1856f79a24a46f7e8f9c674699ef6b1bd63215d49e512c39"],
  transformationRecordSet: "5738283a5accdbd8b736af81982dc46068172ec502f5c43e4113fe7de10c76eb",
  planSourceBindings: "deb9342fc8c2c0ea18ff280aa3c409246d7c96fab9dc46dddb54862640e4cc28",
} as const;

test("passive approved-byte builder agrees with all retained known answers", async () => {
  const candidate = await ordinaryCandidate();
  const verified = verifyApprovedByteFixtureKnownAnswers(candidate, knownAnswers);
  assert.equal(verified.ok, true);
  for (const [role, file] of [
    ["transformerProfile", "objects/transformer-profile.cbor"],
    ["normalizedOptions", "objects/normalized-options.cbor"],
    ["originalManifest", "objects/original-manifest.cbor"],
    ["executableManifest", "objects/executable-manifest.cbor"],
    ["transformationRecordSet", "objects/transformation-record-set.cbor"],
    ["planSourceBindings", "objects/plan-source-bindings.cbor"],
  ] as const) {
    assert.deepEqual(
      Buffer.from(candidate.bytes[role]),
      await readFile(new URL(file, fixtureRoot)),
    );
  }
  assert.deepEqual(
    Buffer.from(candidate.bytes.transformationRecords[0] ?? []),
    await readFile(new URL("objects/transformation-record.cbor", fixtureRoot)),
  );
});

test("asApprovedByteDigest accepts only exactly 32 bytes", () => {
  assert.throws(
    () => asApprovedByteDigest(new Uint8Array(31), "original-source"),
    /must contain exactly 32 bytes/u,
  );
  assert.throws(
    () => asApprovedByteDigest(new Uint8Array(33), "original-source"),
    /must contain exactly 32 bytes/u,
  );
  const digest = asApprovedByteDigest(new Uint8Array(32), "original-source");
  assert.equal(digest.length, 32);
});

test("candidate owns defensive copies and exact ASCII path order", () => {
  const first = Buffer.from("const a: number = 1;\n");
  const second = Buffer.from("const b: number = 2;\n");
  const third = Buffer.from("const z: number = 3;\n");
  const result = buildApprovedByteFixtureCandidate({
    // The entrypoint is a third, distinct path so its own manifest field
    // (encoded separately from the entries list) cannot mask an entry-order
    // bug in "a.ts" vs "A.ts" below.
    entrypoint: "z.ts",
    sources: [
      { logicalPath: "a.ts", originalBytes: second, emittedBytes: second },
      { logicalPath: "A.ts", originalBytes: first, emittedBytes: first },
      { logicalPath: "z.ts", originalBytes: third, emittedBytes: third },
    ],
  });
  assert.equal(result.ok, true);
  if (!result.ok) {
    return;
  }
  const candidate = result.candidate;
  const retained = candidate.bytes.originalManifest.slice();
  first.fill(0);
  second.fill(0);
  third.fill(0);
  assert.deepEqual(candidate.bytes.originalManifest, retained);
  assert.throws(() => {
    (candidate.digests.originalManifest as unknown as number[])[0] = 0;
  }, TypeError);

  // Sources were given out of ASCII order ("a.ts" before "A.ts"); the
  // manifest's entries must be in exact ASCII order ("A.ts" = 0x41 before
  // "a.ts" = 0x61), not input order.
  const manifestBytes = Buffer.from(candidate.bytes.originalManifest);
  const upperIndex = manifestBytes.indexOf("A.ts");
  const lowerIndex = manifestBytes.indexOf("a.ts");
  assert.ok(upperIndex >= 0, '"A.ts" must appear in the encoded manifest');
  assert.ok(lowerIndex >= 0, '"a.ts" must appear in the encoded manifest');
  assert.ok(upperIndex < lowerIndex, "manifest entries must be in exact ASCII path order");
});

test("Unicode, newline, BOM, and invalid UTF-8 remain byte distinctions", () => {
  const lf = Buffer.from('const value: string = "é/é";\n');
  const crlf = Buffer.from('const value: string = "é/é";\r\n');
  const left = buildApprovedByteFixtureCandidate({
    entrypoint: "value.ts",
    sources: [{ logicalPath: "value.ts", originalBytes: lf, emittedBytes: lf }],
  });
  const right = buildApprovedByteFixtureCandidate({
    entrypoint: "value.ts",
    sources: [{ logicalPath: "value.ts", originalBytes: crlf, emittedBytes: crlf }],
  });
  assert.equal(left.ok, true);
  assert.equal(right.ok, true);
  if (!left.ok || !right.ok) {
    return;
  }
  assert.notDeepEqual(
    left.candidate.digests.originalManifest,
    right.candidate.digests.originalManifest,
  );
  assertRefusal(fixture(Buffer.from([0xef, 0xbb, 0xbf, 0x20])), "ORIGINAL_SOURCE_BOM");
  assertRefusal(fixture(Buffer.from([0xc3, 0x28])), "ORIGINAL_SOURCE_UTF8");
});

test("every retained cap is inclusive and cap-plus-one refuses", () => {
  const exact = Buffer.alloc(APPROVED_BYTE_CAPS.originalFileBytes, 0x20);
  assert.equal(fixture(exact).ok, true);
  assertRefusal(fixture(Buffer.alloc(exact.length + 1, 0x20)), "ORIGINAL_SOURCE_BYTES");
  assertRefusal(
    buildApprovedByteFixtureCandidate({
      entrypoint: "ordinary.ts",
      sources: [
        {
          logicalPath: "ordinary.ts",
          originalBytes: new Uint8Array(),
          emittedBytes: Buffer.alloc(APPROVED_BYTE_CAPS.emittedFileBytes + 1, 0x20),
        },
      ],
    }),
    "EMITTED_SOURCE_BYTES",
  );
  assertRefusal(
    buildApprovedByteFixtureCandidate({
      entrypoint: "f00.ts",
      sources: Array.from({ length: APPROVED_BYTE_CAPS.sourceFiles + 1 }, (_, index) => ({
        logicalPath: `f${String(index).padStart(2, "0")}.ts`,
        originalBytes: new Uint8Array(),
        emittedBytes: new Uint8Array(),
      })) as never,
    }),
    "SOURCE_FILE_COUNT",
  );
  const quarter = Buffer.alloc(APPROVED_BYTE_CAPS.originalAggregateBytes / 4, 0x20);
  const withinAggregateCap = buildApprovedByteFixtureCandidate({
    entrypoint: "f0.ts",
    sources: [0, 1, 2, 3].map((index) => ({
      logicalPath: `f${index}.ts`,
      originalBytes: quarter,
      emittedBytes: quarter,
    })) as never,
  });
  assert.equal(withinAggregateCap.ok, true);
  assertRefusal(
    buildApprovedByteFixtureCandidate({
      entrypoint: "f0.ts",
      sources: [0, 1, 2, 3, 4].map((index) => ({
        logicalPath: `f${index}.ts`,
        originalBytes: index === 4 ? Buffer.from(" ") : quarter,
        emittedBytes: index === 4 ? new Uint8Array() : quarter,
      })) as never,
    }),
    "ORIGINAL_AGGREGATE_BYTES",
  );
  assertRefusal(
    buildApprovedByteFixtureCandidate({
      entrypoint: "f0.ts",
      sources: [0, 1, 2, 3, 4].map((index) => ({
        logicalPath: `f${index}.ts`,
        originalBytes: index === 4 ? new Uint8Array() : quarter,
        emittedBytes: index === 4 ? Buffer.from(" ") : quarter,
      })) as never,
    }),
    "EMITTED_AGGREGATE_BYTES",
  );
});

test("known-answer verification refuses object, disposition, and digest-domain mutations", async () => {
  const ordinary = await ordinaryCandidate();
  for (const role of [
    "transformerProfile",
    "normalizedOptions",
    "originalManifest",
    "executableManifest",
    "transformationRecordSet",
    "planSourceBindings",
  ] as const) {
    const mutated = cloneCandidate(ordinary);
    mutated.bytes[role][0] = (mutated.bytes[role][0] ?? 0) ^ 1;
    const verified = verifyApprovedByteFixtureKnownAnswers(
      mutated as unknown as ApprovedByteFixtureCandidate,
      knownAnswers,
    );
    assert.equal(verified.ok, false);
    if (verified.ok) {
      continue;
    }
    assert.deepEqual(
      [verified.refusal.owner, verified.refusal.classification, verified.refusal.code],
      ["approved-byte-fixture-verifier", "DOMAIN", "KNOWN_ANSWER_MISMATCH"],
    );
  }
  const crossDomain = cloneCandidate(ordinary);
  crossDomain.digests.originalManifest = crossDomain.digests.executableManifest;
  const verified = verifyApprovedByteFixtureKnownAnswers(
    crossDomain as unknown as ApprovedByteFixtureCandidate,
    knownAnswers,
  );
  assert.equal(verified.ok, false);
  if (verified.ok) {
    return;
  }
  assert.deepEqual(
    [verified.refusal.owner, verified.refusal.classification, verified.refusal.code],
    ["approved-byte-fixture-verifier", "DOMAIN", "DIGEST_BINDING_MISMATCH"],
  );
});

async function ordinaryCandidate(): Promise<ApprovedByteFixtureCandidate> {
  const result = buildApprovedByteFixtureCandidate({
    entrypoint: "ordinary.ts",
    sources: [
      {
        logicalPath: "ordinary.ts",
        originalBytes: await readFile(new URL("files/ordinary.ts", fixtureRoot)),
        emittedBytes: await readFile(new URL("files/ordinary.js", fixtureRoot)),
      },
    ],
  });
  assert.equal(result.ok, true);
  if (!result.ok) {
    throw new Error("expected ordinary fixture to build");
  }
  return result.candidate;
}

function fixture(bytes: Uint8Array) {
  return buildApprovedByteFixtureCandidate({
    entrypoint: "ordinary.ts",
    sources: [{ logicalPath: "ordinary.ts", originalBytes: bytes, emittedBytes: bytes }],
  });
}

function assertRefusal(
  result: ReturnType<typeof buildApprovedByteFixtureCandidate>,
  code: ApprovedByteFixtureBuildRefusalCode,
): void {
  assert.equal(result.ok, false);
  if (result.ok) {
    return;
  }
  assert.equal(result.refusal.owner, "approved-byte-fixture-builder");
  assert.equal(result.refusal.code, code);
}

type MutableCandidate = {
  bytes: {
    transformerProfile: number[];
    normalizedOptions: number[];
    originalManifest: number[];
    executableManifest: number[];
    transformationRecords: number[][];
    transformationRecordSet: number[];
    planSourceBindings: number[];
  };
  digests: {
    transformerProfile: number[];
    normalizedOptions: number[];
    originalManifest: number[];
    executableManifest: number[];
    transformationRecordSet: number[];
  };
};

function cloneCandidate(value: ApprovedByteFixtureCandidate): MutableCandidate {
  return {
    bytes: {
      transformerProfile: [...value.bytes.transformerProfile],
      normalizedOptions: [...value.bytes.normalizedOptions],
      originalManifest: [...value.bytes.originalManifest],
      executableManifest: [...value.bytes.executableManifest],
      transformationRecords: value.bytes.transformationRecords.map((item) => [...item]),
      transformationRecordSet: [...value.bytes.transformationRecordSet],
      planSourceBindings: [...value.bytes.planSourceBindings],
    },
    digests: {
      transformerProfile: [...value.digests.transformerProfile],
      normalizedOptions: [...value.digests.normalizedOptions],
      executableManifest: [...value.digests.executableManifest],
      originalManifest: [...value.digests.originalManifest],
      transformationRecordSet: [...value.digests.transformationRecordSet],
    },
  };
}
