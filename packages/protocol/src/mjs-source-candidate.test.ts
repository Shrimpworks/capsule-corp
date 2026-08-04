import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  decodeSourceManifestV0,
  encodeSourceManifestV0,
  JOB_PROPOSAL_MEDIA_TYPE,
  MJS_SOURCE_MEDIA_TYPE,
  MJS_SOURCE_PROFILE,
  SOURCE_MANIFEST_MEDIA_TYPE,
  validateMjsSourceBytes,
  validateMjsSourceMediaType,
  validateMjsSourceProfile,
  validateSourceManifestMediaType,
} from "./index.js";

interface FixtureReference {
  readonly path: string;
  readonly sha256: string;
  readonly byteLength: number;
}

interface M1Case {
  readonly id: string;
  readonly object: "MjsSource" | "SourceManifest";
  readonly wireFormat: "raw-bytes" | "cbor" | "media-type";
  readonly mediaType: string | null;
  readonly fixture: FixtureReference;
  readonly context:
    | { readonly kind: "none" }
    | { readonly kind: "source-manifest"; readonly source: FixtureReference };
  readonly expected: {
    readonly decision: "accept" | "reject";
    readonly classification: string | null;
  };
  readonly implementations: { readonly typescript: string };
}

const corpusRoot = new URL("../../../schemas/conformance/v0/", import.meta.url);
const manifest = JSON.parse(await readFile(new URL("manifest.json", corpusRoot), "utf8")) as {
  readonly cases: readonly M1Case[];
};
const m1Cases = manifest.cases.filter(
  (entry) =>
    ["MjsSource", "SourceManifest"].includes(entry.object) &&
    entry.implementations.typescript === "verified",
);

test("M1 TypeScript verifier owns every marked source and manifest case", () => {
  assert.equal(m1Cases.length, 40);
  assert.equal(m1Cases.filter((entry) => entry.object === "MjsSource").length, 21);
  assert.equal(m1Cases.filter((entry) => entry.object === "SourceManifest").length, 19);
});

for (const entry of m1Cases) {
  test(`M1 TypeScript conformance: ${entry.id}`, async () => {
    const fixtureBytes = new Uint8Array(await readFile(new URL(entry.fixture.path, corpusRoot)));
    let accepted: boolean;
    if (entry.wireFormat === "media-type") {
      const value = new TextDecoder().decode(fixtureBytes);
      accepted = entry.id.includes(".profile")
        ? validateMjsSourceProfile(value)
        : entry.id.includes(".member-media")
          ? validateMjsSourceMediaType(value)
          : validateSourceManifestMediaType(value);
    } else if (entry.object === "MjsSource") {
      accepted = validateMjsSourceBytes(fixtureBytes).ok;
    } else {
      assert.equal(entry.context.kind, "source-manifest");
      if (entry.context.kind !== "source-manifest") return;
      const sourceBytes = new Uint8Array(
        await readFile(new URL(entry.context.source.path, corpusRoot)),
      );
      accepted = decodeSourceManifestV0(fixtureBytes, sourceBytes, entry.mediaType ?? "").ok;
    }
    assert.equal(accepted, entry.expected.decision === "accept", entry.id);
  });
}

test("M1 known answers agree on ordinary, minimum, and maximum exact identities", async () => {
  const answers = [
    [
      "mjs-source/ordinary.mjs",
      57,
      "2412c601eb0977eac9a5d4d30f459f48ab787c3f80551e080c38103543d89d75",
    ],
    [
      "mjs-source-manifest/ordinary.cbor",
      89,
      "052dce0c353e1efeb70f93405a8757ef6fa4d29f91a4d6bcaa67d00c45abc0d6",
    ],
    [
      "mjs-source-manifest/minimum.cbor",
      87,
      "b29a880488898dbac54c5890ec56b97a8aca461cbf6585af60fa3a56bc9e9044",
    ],
    [
      "mjs-source-manifest/exact-maximum.cbor",
      95,
      "21b9b521ed0ab6973be563b14f1eae4591097423e17d4d00ba0025c3fe4298d8",
    ],
  ] as const;
  for (const [path, byteLength, digest] of answers) {
    const bytes = await readFile(new URL(path, corpusRoot));
    assert.equal(bytes.byteLength, byteLength, path);
    assert.equal(createHash("sha256").update(bytes).digest("hex"), digest, path);
  }
  for (const name of ["minimum", "ordinary", "exact-maximum"]) {
    const sourceBytes = new Uint8Array(
      await readFile(new URL(`mjs-source/${name}.mjs`, corpusRoot)),
    );
    const source = validateMjsSourceBytes(sourceBytes);
    assert.equal(source.ok, true, name);
    if (!source.ok) continue;
    assert.deepEqual(
      encodeSourceManifestV0(source.source),
      new Uint8Array(await readFile(new URL(`mjs-source-manifest/${name}.cbor`, corpusRoot))),
      name,
    );
  }
});

test("M1 encoder and decoder retain defensive, role-separated exact-byte copies", async () => {
  const callerSource = new Uint8Array(
    await readFile(new URL("mjs-source/ordinary.mjs", corpusRoot)),
  );
  const validated = validateMjsSourceBytes(callerSource);
  assert.equal(validated.ok, true);
  if (!validated.ok) return;
  const expectedManifest = new Uint8Array(
    await readFile(new URL("mjs-source-manifest/ordinary.cbor", corpusRoot)),
  );
  assert.deepEqual(encodeSourceManifestV0(validated.source), expectedManifest);

  const callerManifest = Uint8Array.from(expectedManifest);
  const decoded = decodeSourceManifestV0(callerManifest, callerSource, SOURCE_MANIFEST_MEDIA_TYPE);
  assert.equal(decoded.ok, true);
  if (!decoded.ok) return;

  callerSource[0] = (callerSource[0] ?? 0) ^ 0xff;
  callerManifest[0] = (callerManifest[0] ?? 0) ^ 0xff;
  const manifestCopy = decoded.manifest.copyAuthoritativeBytes();
  const digestCopy = decoded.manifest.copyDigestBytes();
  manifestCopy[0] = (manifestCopy[0] ?? 0) ^ 0xff;
  digestCopy[0] = (digestCopy[0] ?? 0) ^ 0xff;
  assert.deepEqual(decoded.manifest.copyAuthoritativeBytes(), expectedManifest);
  assert.equal(
    Buffer.from(decoded.manifest.copyDigestBytes()).toString("hex"),
    "052dce0c353e1efeb70f93405a8757ef6fa4d29f91a4d6bcaa67d00c45abc0d6",
  );
  assert.equal(decoded.manifest.view.members[0].byteLength, 57);
  assert.equal(
    decoded.manifest.view.members[0].contentDigest,
    "2412c601eb0977eac9a5d4d30f459f48ab787c3f80551e080c38103543d89d75",
  );
  assert.equal(Object.isFrozen(decoded.manifest.view), true);
  assert.equal(Object.isFrozen(decoded.manifest.view.members), true);
  assert.equal(MJS_SOURCE_PROFILE, "capsule.mjs-source/v0");
  assert.equal(JOB_PROPOSAL_MEDIA_TYPE, "application/capsule.job-proposal+json;v=0");
  assert.equal(MJS_SOURCE_MEDIA_TYPE, "application/capsule.javascript-source;v=0;module=esm");
});

test("M1 byte-identity variants remain digest-distinct without normalization", async () => {
  const paths = [
    "mjs-source/identity-lf.mjs",
    "mjs-source/identity-crlf.mjs",
    "mjs-source/identity-lone-cr.mjs",
    "mjs-source/identity-composed.mjs",
    "mjs-source/identity-decomposed.mjs",
    "mjs-source/identity-no-trailing-newline.mjs",
  ];
  const digests = new Set<string>();
  for (const path of paths) {
    const bytes = await readFile(new URL(path, corpusRoot));
    digests.add(createHash("sha256").update(bytes).digest("hex"));
  }
  assert.equal(digests.size, paths.length);
});

test("M1 language validation remains an explicit unimplemented hold", async () => {
  const holdCases = manifest.cases.filter((entry) =>
    entry.id.startsWith("mjs-source.language-validator-hold."),
  );
  assert.equal(holdCases.length, 28);
  assert.equal(
    holdCases.every(
      (entry) =>
        entry.implementations.typescript === "pending" &&
        (entry.implementations as { readonly go?: string }).go === "pending",
    ),
    true,
  );
  const decisions = Object.fromEntries(
    holdCases.map((entry) => [entry.id.split(".").at(-1), entry.expected.decision]),
  );
  assert.deepEqual(
    Object.fromEntries(
      [
        "property-import-meta",
        "method-import",
        "template-interpolation-import",
        "eval-string-data",
        "division-regexp-counterexample",
      ].map((id) => [id, decisions[id]]),
    ),
    {
      "property-import-meta": "accept",
      "method-import": "accept",
      "template-interpolation-import": "reject",
      "eval-string-data": "accept",
      "division-regexp-counterexample": "reject",
    },
  );

  const counterexample = await readFile(
    new URL("mjs-source/language-hold-division-regexp-counterexample.mjs", corpusRoot),
  );
  assert.deepEqual(counterexample, Buffer.from('const of = 9; of / import("evil") / divisor;\n'));
  assert.equal(
    createHash("sha256").update(counterexample).digest("hex"),
    "f5d7998154870fabd6f99ccc53c88e4f9edbe82bf50d63abd8e5db6138b663a1",
  );
});
