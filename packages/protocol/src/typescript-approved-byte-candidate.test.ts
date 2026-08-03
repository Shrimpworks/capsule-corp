import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  APPROVED_BYTE_CAPS,
  type ApprovedByteFixtureCandidate,
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
  verifyApprovedByteFixtureKnownAnswers(candidate, knownAnswers);
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

test("candidate owns defensive copies and exact ASCII path order", () => {
  const first = Buffer.from("const a: number = 1;\n");
  const second = Buffer.from("const b: number = 2;\n");
  const candidate = buildApprovedByteFixtureCandidate({
    entrypoint: "A.ts",
    sources: [
      { logicalPath: "a.ts", originalBytes: second, emittedBytes: second },
      { logicalPath: "A.ts", originalBytes: first, emittedBytes: first },
    ],
  });
  const retained = candidate.bytes.originalManifest.slice();
  first.fill(0);
  second.fill(0);
  assert.deepEqual(candidate.bytes.originalManifest, retained);
  assert.throws(() => {
    (candidate.digests.originalManifest as unknown as number[])[0] = 0;
  }, TypeError);
});

test("Unicode, newline, BOM, and invalid UTF-8 remain byte distinctions", () => {
  const lf = Buffer.from('const value: string = "é/é";\n');
  const crlf = Buffer.from('const value: string = "é/é";\r\n');
  const left = buildApprovedByteFixtureCandidate({
    entrypoint: "value.ts",
    sources: [{ logicalPath: "value.ts", originalBytes: lf, emittedBytes: lf }],
  });
  const right = buildApprovedByteFixtureCandidate({
    entrypoint: "value.ts",
    sources: [{ logicalPath: "value.ts", originalBytes: crlf, emittedBytes: crlf }],
  });
  assert.notDeepEqual(left.digests.originalManifest, right.digests.originalManifest);
  assert.throws(() => fixture(Buffer.from([0xef, 0xbb, 0xbf, 0x20])), /BOM/u);
  assert.throws(() => fixture(Buffer.from([0xc3, 0x28])), /strict UTF-8/u);
});

test("every retained cap is inclusive and cap-plus-one refuses", () => {
  const exact = Buffer.alloc(APPROVED_BYTE_CAPS.originalFileBytes, 0x20);
  fixture(exact);
  assert.throws(() => fixture(Buffer.alloc(exact.length + 1, 0x20)), /original source/u);
  assert.throws(
    () =>
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
    /emitted JavaScript/u,
  );
  assert.throws(
    () =>
      buildApprovedByteFixtureCandidate({
        entrypoint: "f00.ts",
        sources: Array.from({ length: APPROVED_BYTE_CAPS.sourceFiles + 1 }, (_, index) => ({
          logicalPath: `f${String(index).padStart(2, "0")}.ts`,
          originalBytes: new Uint8Array(),
          emittedBytes: new Uint8Array(),
        })) as never,
      }),
    /file count/u,
  );
  const quarter = Buffer.alloc(APPROVED_BYTE_CAPS.originalAggregateBytes / 4, 0x20);
  buildApprovedByteFixtureCandidate({
    entrypoint: "f0.ts",
    sources: [0, 1, 2, 3].map((index) => ({
      logicalPath: `f${index}.ts`,
      originalBytes: quarter,
      emittedBytes: quarter,
    })) as never,
  });
  assert.throws(
    () =>
      buildApprovedByteFixtureCandidate({
        entrypoint: "f0.ts",
        sources: [0, 1, 2, 3, 4].map((index) => ({
          logicalPath: `f${index}.ts`,
          originalBytes: index === 4 ? Buffer.from(" ") : quarter,
          emittedBytes: index === 4 ? new Uint8Array() : quarter,
        })) as never,
      }),
    /original aggregate/u,
  );
  assert.throws(
    () =>
      buildApprovedByteFixtureCandidate({
        entrypoint: "f0.ts",
        sources: [0, 1, 2, 3, 4].map((index) => ({
          logicalPath: `f${index}.ts`,
          originalBytes: index === 4 ? new Uint8Array() : quarter,
          emittedBytes: index === 4 ? Buffer.from(" ") : quarter,
        })) as never,
      }),
    /emitted aggregate/u,
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
    assert.throws(
      () =>
        verifyApprovedByteFixtureKnownAnswers(
          mutated as unknown as ApprovedByteFixtureCandidate,
          knownAnswers,
        ),
      /mismatch/u,
    );
  }
  const crossDomain = cloneCandidate(ordinary);
  crossDomain.digests.originalManifest = crossDomain.digests.executableManifest;
  assert.throws(
    () =>
      verifyApprovedByteFixtureKnownAnswers(
        crossDomain as unknown as ApprovedByteFixtureCandidate,
        knownAnswers,
      ),
    /nominal digest binding/u,
  );
});

async function ordinaryCandidate(): Promise<ApprovedByteFixtureCandidate> {
  return buildApprovedByteFixtureCandidate({
    entrypoint: "ordinary.ts",
    sources: [
      {
        logicalPath: "ordinary.ts",
        originalBytes: await readFile(new URL("files/ordinary.ts", fixtureRoot)),
        emittedBytes: await readFile(new URL("files/ordinary.js", fixtureRoot)),
      },
    ],
  });
}

function fixture(bytes: Uint8Array): ApprovedByteFixtureCandidate {
  return buildApprovedByteFixtureCandidate({
    entrypoint: "ordinary.ts",
    sources: [{ logicalPath: "ordinary.ts", originalBytes: bytes, emittedBytes: bytes }],
  });
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
      originalManifest: [...value.digests.originalManifest],
      executableManifest: [...value.digests.executableManifest],
      transformationRecordSet: [...value.digests.transformationRecordSet],
    },
  };
}
