import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  CandidateRefusal,
  verifyC1,
  verifyManifestBytes,
  verifyParsedManifest,
} from "./verify-governed-runtime-release-candidate.mjs";

const manifestUrl = new URL(
  "../schemas/conformance/governed-deno-core-release-candidate/candidate-manifest.json",
  import.meta.url,
);
const corpusUrl = new URL(
  "../schemas/conformance/governed-deno-core-release-candidate/mutation-corpus.json",
  import.meta.url,
);
const exact = await readFile(manifestUrl);
const manifest = verifyManifestBytes(exact);
const corpus = JSON.parse(await readFile(corpusUrl, "utf8"));

function tokens(pointer) {
  return pointer
    .slice(1)
    .split("/")
    .map((token) => token.replaceAll("~1", "/").replaceAll("~0", "~"));
}

function location(root, pointer) {
  const parts = tokens(pointer);
  const key = parts.pop();
  let parent = root;
  for (const part of parts) parent = parent[part];
  return { parent, key };
}

function valueAt(root, pointer) {
  return tokens(pointer).reduce((value, part) => value[part], root);
}

function applyMutation(root, mutation) {
  const copy = structuredClone(root);
  if (mutation.operation === "swap") {
    const first = location(copy, mutation.pointer);
    const second = location(copy, mutation.otherPointer);
    [first.parent[first.key], second.parent[second.key]] = [
      second.parent[second.key],
      first.parent[first.key],
    ];
    return copy;
  }
  const { parent, key } = location(copy, mutation.pointer);
  if (mutation.operation === "remove") {
    if (Array.isArray(parent)) parent.splice(Number(key), 1);
    else delete parent[key];
  } else if (mutation.operation === "add" && key === "-") {
    parent.push(
      mutation.copyPointer ? structuredClone(valueAt(copy, mutation.copyPointer)) : mutation.value,
    );
  } else {
    parent[key] = mutation.value;
  }
  return copy;
}

test("verifies the exact closed candidate and local passive C1 fixture", async () => {
  assert.equal(
    manifest.selfDigest.sha256,
    "78cf2e99e58a4e79413f22889dd19f794ac7cdce3e4ec5c167d6c2051d19afaa",
  );
  assert.equal(manifest.runtimeSubjects.length, 8);
  assert.equal(manifest.closurePolicy.requiredCandidateEnvelopeFileCount, 32);
  await verifyC1(manifest);
});

for (const mutation of corpus) {
  test(`refuses mutation: ${mutation.name}`, () => {
    let error;
    try {
      if (mutation.operation === "append-byte") {
        verifyManifestBytes(Buffer.concat([exact, Buffer.from(mutation.value)]));
      } else {
        verifyParsedManifest(applyMutation(manifest, mutation));
      }
    } catch (caught) {
      error = caught;
    }
    assert.ok(error instanceof CandidateRefusal, `expected CandidateRefusal, got ${error}`);
    assert.equal(error.code, mutation.expectedRefusal);
  });
}
