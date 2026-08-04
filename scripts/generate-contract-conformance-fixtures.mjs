import { createHash } from "node:crypto";
import { mkdir, readdir, readFile, writeFile } from "node:fs/promises";
import { addApprovalAttemptRulesAndCases } from "./generate-approval-attempt-conformance-fixtures.mjs";
import { addMjsSourceFoundationRulesAndCases } from "./generate-mjs-source-foundation-conformance-fixtures.mjs";
import { addPlanRegistrationRulesAndCases } from "./generate-plan-registration-conformance-fixtures.mjs";

const corpusRoot = new URL("../schemas/conformance/v0/", import.meta.url);
const sharedDirectory = new URL("shared/", corpusRoot);
const jobProposalMediaType = "application/capsule.job-proposal+json;v=0";
const executionPlanMediaType = "application/capsule.execution-plan+cbor;v=0";
const planRegistrationMediaType = "application/capsule.plan-registration+cbor;v=0";
const malformed = "MALFORMED";
const schema = "SCHEMA";
const boundaryRequirements = [
  { decision: "accept", variant: "exact-maximum" },
  { decision: "reject", variant: "cap-plus-one" },
];
const ordinaryRequirements = [
  { decision: "accept", variant: "ordinary" },
  { decision: "reject", variant: "malformed" },
];
const rawImplementations = implementationStatus("pending", "pending", "not-applicable");
const jsonRawImplementations = implementationStatus("pending", "verified", "not-applicable");
const internalMediaTypeImplementations = implementationStatus(
  "verified",
  "pending",
  "not-applicable",
);
const scalarImplementations = implementationStatus("verified", "pending", "pending");
const cborImplementations = implementationStatus("verified", "pending", "pending");
const proposalImplementations = implementationStatus("pending", "verified", "not-applicable");
const proposalSchemaImplementations = implementationStatus("pending", "verified", "not-applicable");
const rules = [];
const cases = [];
const fixtures = new Map();

addMediaTypeRulesAndCases();
addMjsSourceFoundationRulesAndCases({ addCase, addRule, cborEncode, retainFixture });
addJsonRulesAndCases();
addScalarRulesAndCases();
addCborRulesAndCases();
addProposalRulesAndCases();
addPlanRegistrationRulesAndCases({
  addCase,
  addRule,
  cborEncode,
  jsonDocument,
  retainFixture,
  scalarRoleContext,
});
addApprovalAttemptRulesAndCases({
  addCase,
  addRule,
  cborEncode,
  jsonDocument,
  retainFixture,
});

const manifestBytes = utf8(
  `${JSON.stringify({ manifestVersion: "capsule.conformance/v0", rules, cases }, null, 2)}\n`,
);
if (process.argv.includes("--check")) {
  await checkGeneratedCorpus(manifestBytes);
} else {
  await mkdir(sharedDirectory, { recursive: true });
  for (const [path, bytes] of [...fixtures].sort(([left], [right]) => left.localeCompare(right))) {
    const destination = new URL(path, corpusRoot);
    await mkdir(new URL("./", destination), { recursive: true });
    await writeFile(destination, bytes);
  }
  await writeFile(new URL("manifest.json", corpusRoot), manifestBytes);
}
process.stdout.write(
  `${process.argv.includes("--check") ? "verified" : "generated"} conformance corpus: ${rules.length} rules, ${cases.length} cases, ${fixtures.size} fixtures\n`,
);

async function checkGeneratedCorpus(expectedManifest) {
  const expected = new Map(fixtures);
  expected.set("manifest.json", expectedManifest);
  const actualPaths = await listGeneratedCorpusFiles(corpusRoot);
  const expectedPaths = [...expected.keys()].sort();
  if (JSON.stringify(actualPaths) !== JSON.stringify(expectedPaths)) {
    const actual = new Set(actualPaths);
    const wanted = new Set(expectedPaths);
    const extra = actualPaths.filter((path) => !wanted.has(path));
    const missing = expectedPaths.filter((path) => !actual.has(path));
    throw new Error(
      `generated conformance corpus file set is stale; extra=${JSON.stringify(extra)} missing=${JSON.stringify(missing)}`,
    );
  }
  for (const path of expectedPaths) {
    const actual = await readFile(new URL(path, corpusRoot));
    if (!actual.equals(Buffer.from(expected.get(path)))) {
      throw new Error(`generated conformance fixture is stale: ${path}`);
    }
  }
}

async function listGeneratedCorpusFiles(root) {
  const paths = [];
  await visit(root, "");
  return paths.sort();

  async function visit(directory, prefix) {
    for (const entry of await readdir(directory, { withFileTypes: true })) {
      const path = prefix.length === 0 ? entry.name : `${prefix}/${entry.name}`;
      if (path === "manifest.schema.json") continue;
      if (entry.isDirectory()) await visit(new URL(`${entry.name}/`, directory), path);
      else if (entry.isFile()) paths.push(path);
      else throw new Error(`unexpected generated corpus entry: ${path}`);
    }
  }
}

function addMediaTypeRulesAndCases() {
  const profiles = [
    ["job-proposal", "JobProposal", jobProposalMediaType],
    ["execution-plan", "ExecutionPlan", executionPlanMediaType],
    ["plan-registration", "PlanRegistration", planRegistrationMediaType],
  ];
  for (const [slug, object, mediaType] of profiles) {
    const implementations =
      object === "JobProposal" ? rawImplementations : internalMediaTypeImplementations;
    const ruleId = `${slug}.media-type`;
    addRule(
      ruleId,
      "ADR-0023#boundary-ownership-and-media-types",
      `${object} uses only its versioned object-specific media type.`,
      ordinaryRequirements,
    );
    addCase({
      id: `${ruleId}.exact`,
      description: `Accept the exact ${object} media type.`,
      ruleIds: [ruleId],
      object,
      wireFormat: "media-type",
      mediaType,
      variant: "ordinary",
      path: `shared/media-${slug}-exact.txt`,
      bytes: mediaType,
      owner: "media-type-parser",
      implementations,
    });
    const uppercaseType = `${mediaType.slice(0, mediaType.indexOf(";")).toUpperCase()};v=0`;
    addCase({
      id: `${ruleId}.case-insensitive-type`,
      description: `Accept ASCII case differences in the ${object} type and subtype.`,
      ruleIds: [ruleId],
      object,
      wireFormat: "media-type",
      mediaType: uppercaseType,
      variant: "ordinary",
      path: `shared/media-${slug}-uppercase.txt`,
      bytes: uppercaseType,
      owner: "media-type-parser",
      implementations,
    });
    addMediaTypeReject(
      slug,
      object,
      ruleId,
      mediaType.split(";")[0],
      "missing-version",
      malformed,
      implementations,
    );
    addMediaTypeReject(
      slug,
      object,
      ruleId,
      `${mediaType};charset=utf-8`,
      "additional-parameter",
      malformed,
      implementations,
    );
    addMediaTypeReject(
      slug,
      object,
      ruleId,
      mediaType.replace("v=0", "v=1"),
      "unknown-version",
      "UNSUPPORTED",
      implementations,
    );
  }
}

function addMediaTypeReject(
  slug,
  object,
  ruleId,
  mediaType,
  suffix,
  classification,
  implementations,
) {
  addCase({
    id: `${ruleId}.${suffix}`,
    description: `Reject ${suffix.replaceAll("-", " ")} for ${object}.`,
    ruleIds: [ruleId],
    object,
    wireFormat: "media-type",
    mediaType,
    variant: "malformed",
    path: `shared/media-${slug}-${suffix}.txt`,
    bytes: mediaType,
    decision: "reject",
    classification,
    owner: "media-type-parser",
    implementations,
  });
}

function addJsonRulesAndCases() {
  const boundaryRules = [
    ["job-proposal.raw.bytes", "Raw JSON is at most 2,097,152 bytes."],
    ["job-proposal.raw.depth", "JSON value depth is at most 32 including the root value."],
    ["job-proposal.raw.nodes", "JSON contains at most the derived 8,193 value nodes."],
    ["job-proposal.raw.total-members", "JSON contains at most 4,096 object members."],
    ["job-proposal.raw.object-members", "One JSON object contains at most 256 members."],
    ["job-proposal.raw.total-elements", "JSON contains at most 4,096 array elements."],
    ["job-proposal.raw.array-elements", "One JSON array contains at most 256 elements."],
    [
      "job-proposal.raw.decoded-text",
      "Decoded keys and values contain at most 1,572,864 UTF-8 bytes.",
    ],
    ["job-proposal.raw.key-bytes", "One decoded object key contains at most 128 UTF-8 bytes."],
    ["job-proposal.raw.string-bytes", "One non-source string contains at most 65,536 UTF-8 bytes."],
  ];
  for (const [id, description] of boundaryRules) {
    addRule(id, "ADR-0023#strict-jobproposal-raw-profile", description, boundaryRequirements);
  }

  addJsonBoundary(
    "job-proposal.raw.bytes",
    "raw-bytes",
    padJsonToLength("0", 2_097_152),
    padJsonToLength("0", 2_097_153),
  );
  addJsonBoundary("job-proposal.raw.depth", "depth", nestedJsonArray(32), nestedJsonArray(33));
  addJsonBoundary(
    "job-proposal.raw.nodes",
    "nodes",
    jsonBytes(jsonNodeBoundary(false)),
    jsonBytes(jsonNodeBoundary(true)),
  );
  addJsonBoundary(
    "job-proposal.raw.total-members",
    "total-members",
    jsonBytes(jsonMemberBoundary(false)),
    jsonBytes(jsonMemberBoundary(true)),
  );
  addJsonBoundary(
    "job-proposal.raw.object-members",
    "object-members",
    jsonBytes(numberedObject(256)),
    jsonBytes(numberedObject(257)),
  );
  addJsonBoundary(
    "job-proposal.raw.total-elements",
    "total-elements",
    jsonBytes(jsonElementBoundary(false)),
    jsonBytes(jsonElementBoundary(true)),
  );
  addJsonBoundary(
    "job-proposal.raw.array-elements",
    "array-elements",
    jsonBytes(Array.from({ length: 256 }, () => 0)),
    jsonBytes(Array.from({ length: 257 }, () => 0)),
  );
  addJsonBoundary(
    "job-proposal.raw.decoded-text",
    "decoded-text",
    jsonBytes(Array.from({ length: 24 }, () => "a".repeat(65_536))),
    jsonBytes([...Array.from({ length: 23 }, () => "a".repeat(65_536)), "a".repeat(65_535), "aa"]),
  );
  addJsonBoundary(
    "job-proposal.raw.key-bytes",
    "key-bytes",
    jsonBytes({ ["a".repeat(128)]: 0 }),
    jsonBytes({ ["a".repeat(129)]: 0 }),
  );
  addJsonBoundary(
    "job-proposal.raw.string-bytes",
    "string-bytes",
    jsonBytes("a".repeat(65_536)),
    jsonBytes("a".repeat(65_537)),
  );

  const syntaxRules = [
    [
      "job-proposal.raw.single-document",
      "Accept exactly one JSON document with only trailing whitespace.",
    ],
    ["job-proposal.raw.utf8", "Accept only shortest-form valid UTF-8."],
    ["job-proposal.raw.bom", "Reject a UTF-8 BOM."],
    ["job-proposal.raw.surrogate", "Reject unpaired UTF-16 surrogate escapes."],
    ["job-proposal.raw.duplicate-keys", "Reject duplicate decoded object keys."],
    [
      "job-proposal.raw.integer-grammar",
      "Accept only the restricted base-10 integer token grammar.",
    ],
  ];
  for (const [id, description] of syntaxRules) {
    addRule(id, "ADR-0023#strict-jobproposal-raw-profile", description, ordinaryRequirements);
  }
  addRule(
    "job-proposal.raw.integer-range",
    "ADR-0023#strict-jobproposal-raw-profile",
    "Accept integers only inside the inclusive cross-language safe range.",
    boundaryRequirements,
  );

  addJsonCase({
    id: "job-proposal.raw.ordinary.accept",
    description:
      "Accept valid UTF-8 JSON with whitespace, reordered keys, paired scalars, and safe zero.",
    ruleIds: syntaxRules.map(([id]) => id),
    variant: "ordinary",
    path: "shared/json-ordinary-accept.bin",
    bytes: '{\n  "z": 0, "a": "\\ud83d\\ude80"\n}\n',
  });
  addJsonReject("job-proposal.raw.single-document", "empty", new Uint8Array());
  addJsonReject("job-proposal.raw.single-document", "trailing-data", "{}x");
  addJsonReject("job-proposal.raw.single-document", "second-document", "{} {}");
  addJsonReject(
    "job-proposal.raw.utf8",
    "invalid-continuation",
    Uint8Array.of(0x22, 0xc3, 0x28, 0x22),
  );
  addJsonReject("job-proposal.raw.utf8", "overlong", Uint8Array.of(0x22, 0xc0, 0xaf, 0x22));
  addJsonReject(
    "job-proposal.raw.bom",
    "utf8-bom",
    concatenate([Uint8Array.of(0xef, 0xbb, 0xbf), utf8("{}")]),
  );
  addJsonReject("job-proposal.raw.surrogate", "unpaired-surrogate", '"\\ud800"');
  addJsonReject(
    "job-proposal.raw.duplicate-keys",
    "escaped-equivalent",
    '{"kind":0,"k\\u0069nd":1}',
  );
  for (const [suffix, bytes] of [
    ["fraction", "1.0"],
    ["exponent", "1e0"],
    ["leading-plus", "+1"],
    ["leading-zero", "01"],
    ["negative-zero", "-0"],
  ]) {
    addJsonReject("job-proposal.raw.integer-grammar", suffix, bytes);
  }
  addJsonIntegerRange("safe-positive", "9007199254740991", "accept", "exact-maximum");
  addJsonIntegerRange("unsafe-positive", "9007199254740992", "reject", "cap-plus-one");
  addJsonIntegerRange("safe-negative", "-9007199254740991", "accept", "exact-maximum");
  addJsonIntegerRange("unsafe-negative", "-9007199254740992", "reject", "cap-plus-one");
}

function addJsonBoundary(ruleId, slug, exactBytes, overBytes) {
  addJsonCase({
    id: `${ruleId}.exact-maximum`,
    description: `Accept the exact ${slug.replaceAll("-", " ")} maximum at the raw boundary.`,
    ruleIds: [ruleId],
    variant: "exact-maximum",
    path: `shared/json-${slug}-exact.bin`,
    bytes: exactBytes,
  });
  addJsonCase({
    id: `${ruleId}.cap-plus-one`,
    description:
      ruleId === "job-proposal.raw.nodes"
        ? "Reject 8,194 nodes; the derived limit necessarily co-triggers the total-element cap."
        : `Reject ${slug.replaceAll("-", " ")} at cap plus one.`,
    ruleIds: [ruleId],
    variant: "cap-plus-one",
    path: `shared/json-${slug}-over.bin`,
    bytes: overBytes,
    decision: "reject",
    classification: malformed,
  });
}

function addJsonReject(ruleId, suffix, bytes) {
  addJsonCase({
    id: `${ruleId}.${suffix}`,
    description: `Reject ${suffix.replaceAll("-", " ")} JSON bytes.`,
    ruleIds: [ruleId],
    variant: "malformed",
    path: `shared/json-${suffix}.bin`,
    bytes,
    decision: "reject",
    classification: malformed,
  });
}

function addJsonIntegerRange(suffix, bytes, decision, variant) {
  addJsonCase({
    id: `job-proposal.raw.integer-range.${suffix}`,
    description: `${decision === "accept" ? "Accept" : "Reject"} ${suffix.replaceAll("-", " ")} integer bytes.`,
    ruleIds: ["job-proposal.raw.integer-range"],
    variant,
    path: `shared/json-integer-${suffix}.bin`,
    bytes,
    decision,
    classification: decision === "reject" ? malformed : null,
  });
}

function addJsonCase(options) {
  addCase({
    object: "JobProposal",
    wireFormat: "json",
    mediaType: jobProposalMediaType,
    owner: "public-raw-decoder",
    implementations: jsonRawImplementations,
    ...options,
  });
}

function addCborRulesAndCases() {
  const profiles = [
    {
      slug: "execution-plan",
      object: "ExecutionPlan",
      mediaType: executionPlanMediaType,
      owner: "supervisor-plan-predecoder",
      rawBytes: 65_536,
      depth: 8,
      items: 256,
      mapEntries: 64,
      arrayElements: 8,
    },
    {
      slug: "plan-registration",
      object: "PlanRegistration",
      mediaType: planRegistrationMediaType,
      owner: "supervisor-registration-predecoder",
      rawBytes: 4_096,
      depth: 4,
      items: 33,
      mapEntries: 16,
      arrayElements: 0,
    },
  ];
  for (const profile of profiles) {
    for (const [dimension, description] of [
      ["raw-bytes", `Raw deterministic-CBOR payload is at most ${profile.rawBytes} bytes.`],
      ["depth", `CBOR data-item depth is at most ${profile.depth} including the root.`],
      ["items", `CBOR contains at most ${profile.items} total data items.`],
      ["map-entries", `One CBOR map contains at most ${profile.mapEntries} entries.`],
      ["array-elements", `One CBOR array contains at most ${profile.arrayElements} elements.`],
    ]) {
      addRule(
        `${profile.slug}.cbor.${dimension}`,
        "ADR-0023#internal-cbor-predecoder-budgets",
        description,
        boundaryRequirements,
      );
    }
    addRule(
      `${profile.slug}.cbor.deterministic`,
      "ADR-0023#internal-cbor-predecoder-budgets",
      "Accept only the bounded deterministic-CBOR profile before object decoding.",
      ordinaryRequirements,
    );

    addCborBoundary(
      profile,
      "raw-bytes",
      cborByteStringWithTotalLength(profile.rawBytes),
      cborByteStringWithTotalLength(profile.rawBytes + 1),
    );
    addCborBoundary(
      profile,
      "depth",
      cborEncode(
        profile.object === "PlanRegistration"
          ? nestedCborMap(profile.depth)
          : nestedCborArray(profile.depth),
      ),
      cborEncode(
        profile.object === "PlanRegistration"
          ? nestedCborMap(profile.depth + 1)
          : nestedCborArray(profile.depth + 1),
      ),
    );
    addCborBoundary(
      profile,
      "items",
      profile.object === "ExecutionPlan"
        ? cborEncode(planItemBoundary(false))
        : cborEncode(registrationItemBoundary(false)),
      profile.object === "ExecutionPlan"
        ? cborEncode(planItemBoundary(true))
        : cborEncode(registrationItemBoundary(true)),
      profile.object === "PlanRegistration"
        ? "The 34-item rejection also exceeds the zero-array-element cap."
        : undefined,
    );
    addCborBoundary(
      profile,
      "map-entries",
      cborEncode(numberedMap(profile.mapEntries)),
      cborEncode(numberedMap(profile.mapEntries + 1)),
      profile.object === "PlanRegistration"
        ? "The 17-entry rejection also exceeds the 33-item aggregate cap."
        : undefined,
    );
    addCborBoundary(
      profile,
      "array-elements",
      cborEncode(Array.from({ length: profile.arrayElements }, () => 0)),
      cborEncode(Array.from({ length: profile.arrayElements + 1 }, () => 0)),
    );
    addCborDeterministicCases(profile);
  }
}

function addCborBoundary(profile, dimension, exactBytes, overBytes, overlap) {
  const ruleId = `${profile.slug}.cbor.${dimension}`;
  addCborCase(profile, {
    id: `${ruleId}.exact-maximum`,
    description: `Accept the exact ${dimension.replaceAll("-", " ")} maximum at the predecoder.`,
    ruleIds: [ruleId],
    variant: "exact-maximum",
    path: `shared/cbor-${profile.slug}-${dimension}-exact.bin`,
    bytes: exactBytes,
  });
  addCborCase(profile, {
    id: `${ruleId}.cap-plus-one`,
    description: `Reject ${dimension.replaceAll("-", " ")} at cap plus one.${overlap ? ` ${overlap}` : ""}`,
    ruleIds: [ruleId],
    variant: "cap-plus-one",
    path: `shared/cbor-${profile.slug}-${dimension}-over.bin`,
    bytes: overBytes,
    decision: "reject",
    classification: malformed,
  });
}

function addCborDeterministicCases(profile) {
  const ruleId = `${profile.slug}.cbor.deterministic`;
  addCborCase(profile, {
    id: `${ruleId}.canonical-map`,
    description: "Accept a definite, preferred, canonically ordered CBOR map.",
    ruleIds: [ruleId],
    variant: "ordinary",
    path: "shared/cbor-profile-canonical-map.bin",
    bytes: Uint8Array.of(0xa1, 0x01, 0x00),
  });
  for (const [suffix, description, bytes] of [
    ["indefinite-container", "Reject an indefinite-length array.", Uint8Array.of(0x9f, 0x00, 0xff)],
    ["float", "Reject a floating-point value.", Uint8Array.of(0xf9, 0x3c, 0x00)],
    ["bignum", "Reject a bignum tag.", Uint8Array.of(0xc2, 0x41, 0x01)],
    ["tag", "Reject an arbitrary semantic tag.", Uint8Array.of(0xc0, 0x00)],
    [
      "duplicate-key",
      "Reject duplicate decoded map keys.",
      Uint8Array.of(0xa2, 0x01, 0x00, 0x01, 0x01),
    ],
    ["trailing-data", "Reject data after the first CBOR item.", Uint8Array.of(0x00, 0x00)],
    ["nonpreferred-integer", "Reject a nonpreferred integer encoding.", Uint8Array.of(0x18, 0x00)],
    ["invalid-utf8", "Reject invalid UTF-8 in a text string.", Uint8Array.of(0x61, 0xff)],
    [
      "map-order",
      "Reject noncanonical deterministic map order.",
      Uint8Array.of(0xa2, 0x02, 0x00, 0x01, 0x00),
    ],
  ]) {
    addCborCase(profile, {
      id: `${ruleId}.${suffix}`,
      description,
      ruleIds: [ruleId],
      variant: "malformed",
      path: `shared/cbor-profile-${suffix}.bin`,
      bytes,
      decision: "reject",
      classification: malformed,
    });
  }
}

function addCborCase(profile, options) {
  addCase({
    object: profile.object,
    wireFormat: "cbor",
    mediaType: profile.mediaType,
    owner: profile.owner,
    implementations: cborImplementations,
    ...options,
  });
}

function addScalarRulesAndCases() {
  addRule(
    "candidate-id.width",
    "ADR-0023#scalar-zero-rules",
    "Candidate IDs are exactly 16 bytes.",
    boundaryRequirements,
  );
  addRule(
    "candidate-id.nonzero",
    "ADR-0023#scalar-zero-rules",
    "Candidate IDs reject the all-zero value.",
    [
      { decision: "accept", variant: "nonzero" },
      { decision: "reject", variant: "all-zero" },
    ],
  );
  addRule(
    "sha256-digest.width",
    "ADR-0023#scalar-zero-rules",
    "SHA-256 digest fields are exactly 32 bytes.",
    boundaryRequirements,
  );
  addRule(
    "sha256-digest.bit-pattern",
    "ADR-0023#scalar-zero-rules",
    "Every 32-byte digest bit pattern is structurally valid, including all zeroes.",
    [
      { decision: "accept", variant: "all-zero" },
      { decision: "accept", variant: "nonzero" },
    ],
  );
  addRule(
    "sha256-digest.domain",
    "ADR-0023#scalar-zero-rules",
    "Digest width never substitutes for the required nominal role.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "wrong-domain" },
    ],
  );

  const nonzeroId = Uint8Array.from({ length: 16 }, (_, index) => index + 1);
  addScalarCase({
    id: "candidate-id.width.exact",
    description: "Accept exactly 16 candidate-ID bytes.",
    ruleId: "candidate-id.width",
    object: "CandidateId",
    variant: "exact-maximum",
    path: "shared/id-nonzero-16.bin",
    bytes: nonzeroId,
  });
  addScalarReject("candidate-id.width", "short", "CandidateId", new Uint8Array(15), "malformed");
  addScalarReject("candidate-id.width", "long", "CandidateId", new Uint8Array(17), "cap-plus-one");
  addScalarCase({
    id: "candidate-id.nonzero.accept",
    description: "Accept a nonzero 16-byte candidate ID.",
    ruleId: "candidate-id.nonzero",
    object: "CandidateId",
    variant: "nonzero",
    path: "shared/id-nonzero-16.bin",
    bytes: nonzeroId,
  });
  addScalarCase({
    id: "candidate-id.nonzero.reject-zero",
    description: "Reject an all-zero 16-byte candidate ID.",
    ruleId: "candidate-id.nonzero",
    object: "CandidateId",
    variant: "all-zero",
    path: "shared/id-zero-16.bin",
    bytes: new Uint8Array(16),
    decision: "reject",
    classification: schema,
  });

  const nonzeroDigest = Uint8Array.from({ length: 32 }, (_, index) => index + 1);
  addScalarCase({
    id: "sha256-digest.width.exact",
    description: "Accept exactly 32 digest bytes.",
    ruleId: "sha256-digest.width",
    object: "Sha256Digest",
    variant: "exact-maximum",
    path: "shared/digest-nonzero-32.bin",
    bytes: nonzeroDigest,
  });
  addScalarReject("sha256-digest.width", "short", "Sha256Digest", new Uint8Array(31), "malformed");
  addScalarReject(
    "sha256-digest.width",
    "long",
    "Sha256Digest",
    new Uint8Array(33),
    "cap-plus-one",
  );
  addScalarCase({
    id: "sha256-digest.bit-pattern.zero",
    description: "Accept an all-zero digest structurally.",
    ruleId: "sha256-digest.bit-pattern",
    object: "Sha256Digest",
    variant: "all-zero",
    path: "shared/digest-zero-32.bin",
    bytes: new Uint8Array(32),
  });
  addScalarCase({
    id: "sha256-digest.bit-pattern.nonzero",
    description: "Accept a nonzero digest structurally.",
    ruleId: "sha256-digest.bit-pattern",
    object: "Sha256Digest",
    variant: "nonzero",
    path: "shared/digest-nonzero-32.bin",
    bytes: nonzeroDigest,
  });
  addScalarCase({
    id: "sha256-digest.domain.match",
    description: "Accept digest bytes in their expected nominal role.",
    ruleId: "sha256-digest.domain",
    object: "Sha256Digest",
    variant: "ordinary",
    path: "shared/digest-nonzero-32.bin",
    bytes: nonzeroDigest,
    context: scalarRoleContext("source-manifest", "source-manifest"),
    owner: "role-binding-validator",
  });
  addScalarCase({
    id: "sha256-digest.domain.substitution",
    description: "Reject correct-width digest bytes supplied in the wrong nominal role.",
    ruleId: "sha256-digest.domain",
    object: "Sha256Digest",
    variant: "wrong-domain",
    path: "shared/digest-nonzero-32.bin",
    bytes: nonzeroDigest,
    context: scalarRoleContext("source-manifest", "inline-input"),
    decision: "reject",
    classification: "DOMAIN",
    owner: "role-binding-validator",
  });
}

function addScalarReject(ruleId, suffix, object, bytes, variant) {
  addScalarCase({
    id: `${ruleId}.${suffix}`,
    description: `Reject ${suffix} ${object} bytes.`,
    ruleId,
    object,
    variant,
    path: `shared/${object === "CandidateId" ? "id" : "digest"}-${suffix}.bin`,
    bytes,
    decision: "reject",
    classification: schema,
  });
}

function addScalarCase({ ruleId, ...options }) {
  addCase({
    wireFormat: "raw-bytes",
    mediaType: null,
    ruleIds: [ruleId],
    context: { kind: "none" },
    owner: "shared-scalar-validator",
    implementations: scalarImplementations,
    ...options,
  });
}

function addProposalRulesAndCases() {
  const profileRegistry = retainFixture(
    "contexts/profile-registry.json",
    jsonDocument({
      contextType: "capsule.conformance.profile-registry",
      contextVersion: 0,
      profiles: [
        {
          alias: "fixture-active@1",
          status: "active",
          exactWallTimeMs: { minimum: 1, maximum: 10_000 },
        },
        {
          alias: "fixture-inactive@1",
          status: "inactive",
          exactWallTimeMs: { minimum: 1, maximum: 10_000 },
        },
      ],
    }),
  );
  const userPolicy = retainFixture(
    "contexts/user-policy.json",
    jsonDocument({
      contextType: "capsule.conformance.user-policy",
      contextVersion: 0,
      wallTimeMs: { trustedDefault: 5_000, ceiling: 10_000 },
    }),
  );
  const resolverContext = ({
    sourceManifest = null,
    canonicalInlineInput = null,
    wallTime = null,
  } = {}) => ({
    kind: "proposal-resolution",
    profileRegistry,
    userPolicy,
    oracle: {
      sourceManifest,
      sourceManifestDigest: sourceManifest?.sha256 ?? null,
      canonicalInlineInput,
      inlineInputDigest: canonicalInlineInput?.sha256 ?? null,
      wallTime,
    },
  });

  const proposal = ordinaryProposal();
  const proposalBytes = jsonDocument(proposal);
  const proposalPath = "job-proposal/ordinary.json";
  const sourceManifest = retainFixture(
    "source-manifest/ordinary.cbor",
    encodeSourceManifest(proposal.source),
  );
  const canonicalInlineInput = retainFixture(
    "canonical-inline-input/ordinary.json",
    encodeCanonicalInlineJson(proposal.input.value),
  );
  const knownAnswerContext = resolverContext({
    sourceManifest,
    canonicalInlineInput,
    wallTime: { milliseconds: 5_000, origin: "requested" },
  });

  addRule(
    "job-proposal.schema.authority-omission",
    "ADR-0017#decision",
    "The closed proposal omits unsupported authority instead of accepting false placeholders.",
    ordinaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.schema.authority-omission.accept",
    description: "Accept the closed first-slice proposal with unsupported authority omitted.",
    ruleIds: ["job-proposal.schema.authority-omission"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
    owner: "job-proposal-schema",
  });

  addProposalCase({
    id: "job-proposal.schema.authority-omission.reject-placeholder",
    description: "Reject a false network placeholder because unsupported authority is omitted.",
    ruleIds: ["job-proposal.schema.authority-omission"],
    variant: "malformed",
    path: "job-proposal/authority-network-placeholder.json",
    proposal: { ...proposal, network: false },
    decision: "reject",
    classification: "UNSUPPORTED",
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.schema.authority-omission.reject-missing-source",
    description: "Reject a proposal that omits the required source field.",
    ruleIds: ["job-proposal.schema.authority-omission"],
    variant: "malformed",
    path: "job-proposal/missing-source.json",
    proposal: withoutProperty(proposal, "source"),
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });

  addRule(
    "job-proposal.source.path-grammar",
    "ADR-0023#source-path-and-source-manifest-identity",
    "Source paths use the exact relative case-sensitive ASCII grammar.",
    ordinaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.source.path-grammar.accept",
    description: "Accept case-distinct relative ASCII source paths.",
    ruleIds: ["job-proposal.source.path-grammar"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.source.path-grammar.reject-dot-segment",
    description: "Reject a dot-dot source path segment before semantic resolution.",
    ruleIds: ["job-proposal.source.path-grammar"],
    variant: "malformed",
    path: "job-proposal/source-path-dot-segment.json",
    proposal: withSource(proposal, {
      entrypoint: "main.ts",
      files: { "main.ts": "export {};\n", "src/../escape.ts": "export {};\n" },
    }),
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });

  addRule(
    "job-proposal.source.path-bytes",
    "ADR-0023#source-path-and-source-manifest-identity",
    "One source path contains at most 256 ASCII bytes.",
    boundaryRequirements,
  );
  const exactPath = `${"a".repeat(64)}/${"b".repeat(64)}/${"c".repeat(64)}/${"d".repeat(5)}/${"e".repeat(52)}.ts`;
  const overPath = `${"a".repeat(64)}/${"b".repeat(64)}/${"c".repeat(64)}/${"d".repeat(6)}/${"e".repeat(52)}.ts`;
  addProposalCase({
    id: "job-proposal.source.path-bytes.exact-maximum",
    description: "Accept an exact 256-byte source path and matching entrypoint.",
    ruleIds: ["job-proposal.source.path-bytes"],
    variant: "exact-maximum",
    path: "job-proposal/source-path-exact.json",
    proposal: withSource(proposal, {
      entrypoint: exactPath,
      files: { [exactPath]: "export {};\n" },
    }),
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.source.path-bytes.cap-plus-one",
    description: "Reject a 257-byte source path without truncation or normalization.",
    ruleIds: ["job-proposal.source.path-bytes"],
    variant: "cap-plus-one",
    path: "job-proposal/source-path-over.json",
    proposal: withSource(proposal, {
      entrypoint: overPath,
      files: { [overPath]: "export {};\n" },
    }),
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });

  addRule(
    "job-proposal.source.segment-bytes",
    "ADR-0023#source-path-and-source-manifest-identity",
    "One source path segment contains at most 64 ASCII bytes.",
    boundaryRequirements,
  );
  const exactSegmentPath = `${"s".repeat(64)}/helper.ts`;
  const overSegmentPath = `${"s".repeat(65)}/helper.ts`;
  addProposalCase({
    id: "job-proposal.source.segment-bytes.exact-maximum",
    description: "Accept an exact 64-byte non-entrypoint source path segment.",
    ruleIds: ["job-proposal.source.segment-bytes"],
    variant: "exact-maximum",
    path: "job-proposal/source-segment-exact.json",
    proposal: withSource(proposal, {
      entrypoint: "main.ts",
      files: { "main.ts": "export {};\n", [exactSegmentPath]: "export {};\n" },
    }),
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.source.segment-bytes.cap-plus-one",
    description: "Reject a 65-byte source path segment.",
    ruleIds: ["job-proposal.source.segment-bytes"],
    variant: "cap-plus-one",
    path: "job-proposal/source-segment-over.json",
    proposal: withSource(proposal, {
      entrypoint: "main.ts",
      files: { "main.ts": "export {};\n", [overSegmentPath]: "export {};\n" },
    }),
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });

  addRule(
    "job-proposal.source.file-count",
    "ADR-0023#strict-jobproposal-raw-profile",
    "A proposal contains at most 32 source file entries.",
    boundaryRequirements,
  );
  const exactFiles = numberedSourceFiles(32);
  const overFiles = numberedSourceFiles(33);
  addProposalCase({
    id: "job-proposal.source.file-count.exact-maximum",
    description: "Accept exactly 32 source file entries.",
    ruleIds: ["job-proposal.source.file-count"],
    variant: "exact-maximum",
    path: "job-proposal/source-file-count-exact.json",
    proposal: withSource(proposal, { entrypoint: "f00.ts", files: exactFiles }),
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.source.file-count.cap-plus-one",
    description: "Reject 33 source file entries.",
    ruleIds: ["job-proposal.source.file-count"],
    variant: "cap-plus-one",
    path: "job-proposal/source-file-count-over.json",
    proposal: withSource(proposal, { entrypoint: "f00.ts", files: overFiles }),
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });

  addRule(
    "job-proposal.source.file-bytes",
    "ADR-0023#strict-jobproposal-raw-profile",
    "One decoded source file contains at most 262,144 strict UTF-8 bytes.",
    boundaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.source.file-bytes.exact-maximum",
    description: "Accept one source file containing exactly 262,144 UTF-8 bytes.",
    ruleIds: ["job-proposal.source.file-bytes"],
    variant: "exact-maximum",
    path: "job-proposal/source-file-bytes-exact.json",
    proposal: withSource(proposal, {
      entrypoint: "main.ts",
      files: { "main.ts": "🚀".repeat(65_536) },
    }),
  });
  addProposalCase({
    id: "job-proposal.source.file-bytes.cap-plus-one",
    description: "Reject one source file containing 262,145 UTF-8 bytes.",
    ruleIds: ["job-proposal.source.file-bytes"],
    variant: "cap-plus-one",
    path: "job-proposal/source-file-bytes-over.json",
    proposal: withSource(proposal, {
      entrypoint: "main.ts",
      files: { "main.ts": `${"🚀".repeat(65_536)}a` },
    }),
    decision: "reject",
    classification: "SEMANTIC",
  });

  addRule(
    "job-proposal.source.aggregate-bytes",
    "ADR-0023#source-path-and-source-manifest-identity",
    "The sum of strict UTF-8 source content is at most 1,048,576 bytes.",
    boundaryRequirements,
  );
  const aggregateFiles = Object.fromEntries(
    Array.from({ length: 4 }, (_, index) => [`f${index}.ts`, "a".repeat(262_144)]),
  );
  addProposalCase({
    id: "job-proposal.source.aggregate-bytes.exact-maximum",
    description: "Accept exactly 1,048,576 aggregate source bytes without rewriting them.",
    ruleIds: ["job-proposal.source.aggregate-bytes", "job-proposal.source.manifest-identity"],
    variant: "exact-maximum",
    path: "job-proposal/source-aggregate-exact.json",
    proposal: withSource(proposal, { entrypoint: "f0.ts", files: aggregateFiles }),
  });
  addProposalCase({
    id: "job-proposal.source.aggregate-bytes.cap-plus-one",
    description: "Reject 1,048,577 aggregate source bytes without clamping a file.",
    ruleIds: ["job-proposal.source.aggregate-bytes"],
    variant: "cap-plus-one",
    path: "job-proposal/source-aggregate-over.json",
    proposal: withSource(proposal, {
      entrypoint: "f0.ts",
      files: { ...aggregateFiles, "over.ts": "a" },
    }),
    decision: "reject",
    classification: "SEMANTIC",
  });

  addRule(
    "job-proposal.source.manifest-identity",
    "ADR-0023#source-path-and-source-manifest-identity",
    "SourceManifest entries sort by unsigned ASCII path bytes and bind exact content bytes.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "accept", variant: "exact-maximum" },
    ],
  );
  addProposalCase({
    id: "job-proposal.source.manifest-identity.known-answer",
    description: "Derive the retained ordered SourceManifest bytes and SHA-256 digest.",
    ruleIds: ["job-proposal.source.manifest-identity"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
  });

  addRule(
    "job-proposal.source.entrypoint-membership",
    "ADR-0023#source-path-and-source-manifest-identity",
    "The entrypoint exactly equals one source-file key.",
    ordinaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.source.entrypoint-membership.accept",
    description: "Accept an entrypoint that exactly matches a source path.",
    ruleIds: ["job-proposal.source.entrypoint-membership"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
  });
  addProposalCase({
    id: "job-proposal.source.entrypoint-membership.reject-missing",
    description: "Reject a valid entrypoint path that is absent from the source map.",
    ruleIds: ["job-proposal.source.entrypoint-membership"],
    variant: "malformed",
    path: "job-proposal/entrypoint-not-member.json",
    proposal: withSource(proposal, { ...proposal.source, entrypoint: "missing.ts" }),
    decision: "reject",
    classification: "SEMANTIC",
  });

  addRule(
    "job-proposal.input.canonical-identity",
    "ADR-0023#canonical-inline-json-identity",
    "Canonical inline JSON has one exact sorted, escaped, non-normalizing byte identity.",
    ordinaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.input.canonical-identity.known-answer",
    description: "Derive the retained canonical inline JSON bytes and SHA-256 digest.",
    ruleIds: ["job-proposal.input.canonical-identity"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
  });
  const equivalentBytes = utf8(`{
  "labels": {"example":"known-answer"},
  "outputs": [{"maxBytes":65536,"kind":"inline-json","slot":"transformed-json"}],
  "requestedLimits": {"wallTimeMs":5000},
  "input": {"value":{"z":[3,0,-2,true,null],"a":"quote:\\" slash:/ backslash:\\\\ control:\\u0000\\u001f rocket:\\ud83d\\ude80","A":{"b":false,"a":"\\u00e9"}},"kind":"inline-json","slot":"primary-data"},
  "runtimeProfile": "fixture-active@1",
  "source": {"files":{"A.ts":"export const a = 1;\\r\\n","src/main.ts":"import { a } from \\"../A.ts\\";\\nconsole.log(a);\\n","z.ts":"export const z = \\"\\ud83d\\ude80\\";\\n"},"entrypoint":"src/main.ts"},
  "kind": "JobProposal",
  "apiVersion": "capsule.dev/v0"
}
`);
  addProposalCase({
    id: "job-proposal.input.canonical-identity.equivalent-public-json",
    description: "Different public key order and escape spelling derive the same canonical bytes.",
    ruleIds: ["job-proposal.input.canonical-identity"],
    variant: "ordinary",
    path: "job-proposal/equivalent-public-json.json",
    bytes: equivalentBytes,
    context: knownAnswerContext,
  });
  addProposalCase({
    id: "job-proposal.input.canonical-identity.reject-unsupported-key",
    description: "Reject a canonical inline-input object key outside the candidate ASCII grammar.",
    ruleIds: ["job-proposal.input.canonical-identity"],
    variant: "malformed",
    path: "job-proposal/inline-input-unicode-key.json",
    proposal: withInputValue(proposal, { é: 1 }),
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });

  addRule(
    "job-proposal.input.canonical-bytes",
    "ADR-0023#canonical-inline-json-identity",
    "Canonical inline-input JSON contains at most 262,144 bytes.",
    boundaryRequirements,
  );
  const exactCanonicalInput = [
    "a".repeat(65_536),
    "b".repeat(65_536),
    "c".repeat(65_536),
    "d".repeat(65_523),
  ];
  const overCanonicalInput = [...exactCanonicalInput.slice(0, 3), "d".repeat(65_524)];
  addProposalCase({
    id: "job-proposal.input.canonical-bytes.exact-maximum",
    description: "Accept exactly 262,144 canonical inline-input bytes.",
    ruleIds: ["job-proposal.input.canonical-bytes"],
    variant: "exact-maximum",
    path: "job-proposal/inline-input-bytes-exact.json",
    proposal: withInputValue(proposal, exactCanonicalInput),
  });
  addProposalCase({
    id: "job-proposal.input.canonical-bytes.cap-plus-one",
    description: "Reject 262,145 canonical inline-input bytes without resizing the cap.",
    ruleIds: ["job-proposal.input.canonical-bytes"],
    variant: "cap-plus-one",
    path: "job-proposal/inline-input-bytes-over.json",
    proposal: withInputValue(proposal, overCanonicalInput),
    decision: "reject",
    classification: "SEMANTIC",
  });

  addRule(
    "job-proposal.slots.fixed-roles",
    "ADR-0017#decision",
    "The first slice uses only primary-data input and transformed-json output roles.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "wrong-domain" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addProposalCase({
    id: "job-proposal.slots.fixed-roles.accept",
    description: "Accept the fixed input and output slot roles.",
    ruleIds: ["job-proposal.slots.fixed-roles"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.slots.fixed-roles.reject-input-substitution",
    description: "Reject the output slot substituted into the input role.",
    ruleIds: ["job-proposal.slots.fixed-roles"],
    variant: "wrong-domain",
    path: "job-proposal/input-slot-wrong-domain.json",
    proposal: { ...proposal, input: { ...proposal.input, slot: "transformed-json" } },
    decision: "reject",
    classification: "DOMAIN",
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.slots.fixed-roles.reject-unknown",
    description: "Reject an unknown logical input slot.",
    ruleIds: ["job-proposal.slots.fixed-roles"],
    variant: "malformed",
    path: "job-proposal/input-slot-unknown.json",
    proposal: { ...proposal, input: { ...proposal.input, slot: "other-data" } },
    decision: "reject",
    classification: schema,
    owner: "job-proposal-schema",
  });
  addProposalCase({
    id: "job-proposal.slots.fixed-roles.reject-output-substitution",
    description: "Reject the input slot substituted into the output role.",
    ruleIds: ["job-proposal.slots.fixed-roles"],
    variant: "wrong-domain",
    path: "job-proposal/output-slot-wrong-domain.json",
    proposal: {
      ...proposal,
      outputs: [{ ...proposal.outputs[0], slot: "primary-data" }],
    },
    decision: "reject",
    classification: "DOMAIN",
    owner: "job-proposal-schema",
  });

  addLabelCases(proposal);

  addRule(
    "job-proposal.profile.resolution",
    "PHASE-2B-BOUNDARY-DECISIONS#task-23-add-proposalsourceinput-fixtures",
    "The requested runtime profile alias resolves to one active fixed-context entry.",
    ordinaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.profile.resolution.active",
    description: "Resolve the active fixture-only profile from the fixed registry context.",
    ruleIds: ["job-proposal.profile.resolution"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
  });
  addProposalCase({
    id: "job-proposal.profile.resolution.reject-unknown",
    description: "Reject a structurally valid alias absent from the fixed profile context.",
    ruleIds: ["job-proposal.profile.resolution"],
    variant: "malformed",
    path: "job-proposal/profile-unknown.json",
    proposal: { ...proposal, runtimeProfile: "fixture-unknown@1" },
    decision: "reject",
    classification: "BINDING",
  });
  addProposalCase({
    id: "job-proposal.profile.resolution.reject-inactive",
    description: "Reject a known but inactive profile under the fixed resolver context.",
    ruleIds: ["job-proposal.profile.resolution"],
    variant: "malformed",
    path: "job-proposal/profile-inactive.json",
    proposal: { ...proposal, runtimeProfile: "fixture-inactive@1" },
    decision: "reject",
    classification: "POLICY",
  });

  addRule(
    "job-proposal.policy.wall-time",
    "ADR-0009#decision",
    "Wall time is requested or defaulted exactly, and ceiling plus one rejects without clamping.",
    boundaryRequirements,
  );
  addProposalCase({
    id: "job-proposal.policy.wall-time.requested",
    description: "Preserve a requested wall time below the ceiling exactly.",
    ruleIds: ["job-proposal.policy.wall-time"],
    variant: "ordinary",
    path: proposalPath,
    bytes: proposalBytes,
    context: knownAnswerContext,
  });
  const defaultedProposal = { ...proposal, requestedLimits: {} };
  addProposalCase({
    id: "job-proposal.policy.wall-time.trusted-default",
    description: "Resolve omitted wall time to the unchanged trusted default.",
    ruleIds: ["job-proposal.policy.wall-time"],
    variant: "ordinary",
    path: "job-proposal/wall-time-defaulted.json",
    proposal: defaultedProposal,
    context: resolverContext({
      wallTime: { milliseconds: 5_000, origin: "trusted-default" },
    }),
  });
  const ceilingProposal = {
    ...proposal,
    requestedLimits: { wallTimeMs: 10_000 },
  };
  addProposalCase({
    id: "job-proposal.policy.wall-time.exact-ceiling",
    description: "Accept the exact 10,000 ms trusted ceiling unchanged.",
    ruleIds: ["job-proposal.policy.wall-time"],
    variant: "exact-maximum",
    path: "job-proposal/wall-time-exact-ceiling.json",
    proposal: ceilingProposal,
    context: resolverContext({ wallTime: { milliseconds: 10_000, origin: "requested" } }),
  });
  addProposalCase({
    id: "job-proposal.policy.wall-time.ceiling-plus-one",
    description: "Reject 10,001 ms rather than clamping it to the trusted ceiling.",
    ruleIds: ["job-proposal.policy.wall-time"],
    variant: "cap-plus-one",
    path: "job-proposal/wall-time-over-ceiling.json",
    proposal: { ...proposal, requestedLimits: { wallTimeMs: 10_001 } },
    decision: "reject",
    classification: "POLICY",
  });

  function addProposalCase({ proposal: caseProposal, ...options }) {
    const owner = options.owner ?? "job-proposal-semantic-validator";
    addCase({
      object: "JobProposal",
      wireFormat: "json",
      mediaType: jobProposalMediaType,
      context: resolverContext(),
      owner,
      implementations:
        owner === "job-proposal-schema" ? proposalSchemaImplementations : proposalImplementations,
      bytes: caseProposal === undefined ? options.bytes : jsonDocument(caseProposal),
      ...options,
    });
  }
}

function addLabelCases(proposal) {
  for (const [dimension, description] of [
    ["count", "Labels contain at most 8 entries."],
    ["key-bytes", "One label key contains at most 32 ASCII bytes."],
    ["value-bytes", "One label value contains at most 128 printable-ASCII bytes."],
  ]) {
    addRule(
      `job-proposal.labels.${dimension}`,
      "ADR-0023#strict-jobproposal-raw-profile",
      description,
      boundaryRequirements,
    );
  }
  addRule(
    "job-proposal.labels.printable-ascii",
    "ADR-0023#strict-jobproposal-raw-profile",
    "Label values contain only printable ASCII bytes.",
    ordinaryRequirements,
  );

  const casesForLabels = [
    {
      id: "job-proposal.labels.count.exact-maximum",
      ruleIds: ["job-proposal.labels.count"],
      variant: "exact-maximum",
      path: "job-proposal/labels-count-exact.json",
      labels: Object.fromEntries(Array.from({ length: 8 }, (_, index) => [`label-${index}`, "v"])),
      description: "Accept exactly 8 labels.",
    },
    {
      id: "job-proposal.labels.count.cap-plus-one",
      ruleIds: ["job-proposal.labels.count"],
      variant: "cap-plus-one",
      path: "job-proposal/labels-count-over.json",
      labels: Object.fromEntries(Array.from({ length: 9 }, (_, index) => [`label-${index}`, "v"])),
      description: "Reject 9 labels.",
      decision: "reject",
      classification: schema,
    },
    {
      id: "job-proposal.labels.key-bytes.exact-maximum",
      ruleIds: ["job-proposal.labels.key-bytes"],
      variant: "exact-maximum",
      path: "job-proposal/label-key-exact.json",
      labels: { [`a${"b".repeat(31)}`]: "v" },
      description: "Accept a 32-byte label key.",
    },
    {
      id: "job-proposal.labels.key-bytes.cap-plus-one",
      ruleIds: ["job-proposal.labels.key-bytes"],
      variant: "cap-plus-one",
      path: "job-proposal/label-key-over.json",
      labels: { [`a${"b".repeat(32)}`]: "v" },
      description: "Reject a 33-byte label key.",
      decision: "reject",
      classification: schema,
    },
    {
      id: "job-proposal.labels.value-bytes.exact-maximum",
      ruleIds: ["job-proposal.labels.value-bytes"],
      variant: "exact-maximum",
      path: "job-proposal/label-value-exact.json",
      labels: { example: "x".repeat(128) },
      description: "Accept a 128-byte printable-ASCII label value.",
    },
    {
      id: "job-proposal.labels.value-bytes.cap-plus-one",
      ruleIds: ["job-proposal.labels.value-bytes"],
      variant: "cap-plus-one",
      path: "job-proposal/label-value-over.json",
      labels: { example: "x".repeat(129) },
      description: "Reject a 129-byte label value.",
      decision: "reject",
      classification: schema,
    },
    {
      id: "job-proposal.labels.printable-ascii.accept",
      ruleIds: ["job-proposal.labels.printable-ascii"],
      variant: "ordinary",
      path: "job-proposal/label-printable.json",
      labels: { example: "AZaz09 !~" },
      description: "Accept printable ASCII label bytes.",
    },
    {
      id: "job-proposal.labels.printable-ascii.reject-control",
      ruleIds: ["job-proposal.labels.printable-ascii"],
      variant: "malformed",
      path: "job-proposal/label-control.json",
      labels: { example: "line\nbreak" },
      description: "Reject a control byte in a label value.",
      decision: "reject",
      classification: schema,
    },
  ];
  for (const { labels, ...entry } of casesForLabels) {
    addCase({
      object: "JobProposal",
      wireFormat: "json",
      mediaType: jobProposalMediaType,
      context: { kind: "none" },
      owner: "job-proposal-schema",
      implementations: proposalSchemaImplementations,
      bytes: jsonDocument({ ...proposal, labels }),
      ...entry,
    });
  }
}

function ordinaryProposal() {
  return {
    apiVersion: "capsule.dev/v0",
    kind: "JobProposal",
    source: {
      entrypoint: "src/main.ts",
      files: {
        "z.ts": 'export const z = "🚀";\n',
        "src/main.ts": 'import { a } from "../A.ts";\nconsole.log(a);\n',
        "A.ts": "export const a = 1;\r\n",
      },
    },
    runtimeProfile: "fixture-active@1",
    input: {
      slot: "primary-data",
      kind: "inline-json",
      value: {
        z: [3, 0, -2, true, null],
        a: 'quote:" slash:/ backslash:\\ control:\u0000\u001f rocket:🚀',
        A: { b: false, a: "é" },
      },
    },
    requestedLimits: { wallTimeMs: 5_000 },
    outputs: [{ slot: "transformed-json", kind: "inline-json", maxBytes: 65_536 }],
    labels: { example: "known-answer" },
  };
}

function withSource(proposal, source) {
  return { ...proposal, source };
}

function withInputValue(proposal, value) {
  return { ...proposal, input: { ...proposal.input, value } };
}

function withoutProperty(value, property) {
  const copy = { ...value };
  delete copy[property];
  return copy;
}

function numberedSourceFiles(count) {
  return Object.fromEntries(
    Array.from({ length: count }, (_, index) => [`f${String(index).padStart(2, "0")}.ts`, ""]),
  );
}

function encodeSourceManifest(source) {
  const entries = Object.entries(source.files)
    .sort(([left], [right]) => Buffer.compare(utf8(left), utf8(right)))
    .map(([path, content]) => {
      const contentBytes = utf8(content);
      return [path, sha256Bytes(contentBytes), contentBytes.length];
    });
  const aggregateBytes = entries.reduce((total, entry) => total + entry[2], 0);
  return cborEncode(
    new Map([
      [1, "capsule.source-manifest"],
      [2, 0],
      [3, source.entrypoint],
      [4, entries],
      [5, aggregateBytes],
    ]),
  );
}

function encodeCanonicalInlineJson(value) {
  if (value === null) {
    return utf8("null");
  }
  if (typeof value === "boolean") {
    return utf8(value ? "true" : "false");
  }
  if (typeof value === "number") {
    if (!Number.isSafeInteger(value) || Object.is(value, -0)) {
      throw new Error(`unsupported canonical inline JSON number: ${value}`);
    }
    return utf8(String(value));
  }
  if (typeof value === "string") {
    return utf8(escapeCanonicalJsonString(value));
  }
  if (Array.isArray(value)) {
    return concatenate([
      utf8("["),
      ...value.flatMap((child, index) => [
        index === 0 ? new Uint8Array() : utf8(","),
        encodeCanonicalInlineJson(child),
      ]),
      utf8("]"),
    ]);
  }
  if (typeof value === "object") {
    const entries = Object.entries(value).sort(([left], [right]) =>
      Buffer.compare(utf8(left), utf8(right)),
    );
    return concatenate([
      utf8("{"),
      ...entries.flatMap(([key, child], index) => [
        index === 0 ? new Uint8Array() : utf8(","),
        utf8(escapeCanonicalJsonString(key)),
        utf8(":"),
        encodeCanonicalInlineJson(child),
      ]),
      utf8("}"),
    ]);
  }
  throw new Error(`unsupported canonical inline JSON value: ${String(value)}`);
}

function escapeCanonicalJsonString(value) {
  let result = '"';
  for (const scalar of value) {
    const codePoint = scalar.codePointAt(0);
    if (scalar === '"') {
      result += '\\"';
    } else if (scalar === "\\") {
      result += "\\\\";
    } else if (codePoint <= 0x1f) {
      result += `\\u00${codePoint.toString(16).padStart(2, "0")}`;
    } else {
      result += scalar;
    }
  }
  return `${result}"`;
}

function sha256Bytes(bytes) {
  return createHash("sha256").update(bytes).digest();
}

function jsonDocument(value) {
  return utf8(`${JSON.stringify(value, null, 2)}\n`);
}

function addRule(id, source, description, requiredCases) {
  rules.push({ id, source, description, requiredCases });
}

function addCase({
  id,
  description,
  ruleIds,
  object,
  wireFormat,
  mediaType,
  variant,
  path,
  bytes,
  context = { kind: "none" },
  decision = "accept",
  classification = null,
  owner,
  implementations,
  authorityStateChanged = false,
  timeHighWaterChanged,
  trustStateTightened,
  fakeBackendEffectPermitted,
  stateDelta,
}) {
  const fixture = retainFixture(path, bytes);
  const expected = { decision, classification, owner, authorityStateChanged };
  if (timeHighWaterChanged !== undefined) {
    expected.timeHighWaterChanged = timeHighWaterChanged;
  }
  if (trustStateTightened !== undefined) {
    expected.trustStateTightened = trustStateTightened;
  }
  if (fakeBackendEffectPermitted !== undefined) {
    expected.fakeBackendEffectPermitted = fakeBackendEffectPermitted;
  }
  if (stateDelta !== undefined) {
    expected.stateDelta = stateDelta;
  }
  cases.push({
    id,
    description,
    ruleIds,
    object,
    wireFormat,
    mediaType,
    variant,
    fixture,
    context,
    expected,
    implementations,
  });
}

function retainFixture(path, bytes) {
  const retainedBytes = typeof bytes === "string" ? utf8(bytes) : Buffer.from(bytes);
  const existing = fixtures.get(path);
  if (existing && !existing.equals(retainedBytes)) {
    throw new Error(`fixture path ${path} was assigned different bytes`);
  }
  fixtures.set(path, retainedBytes);
  return {
    path,
    sha256: createHash("sha256").update(retainedBytes).digest("hex"),
    byteLength: retainedBytes.length,
  };
}

function implementationStatus(go, typescript, swift) {
  return { "fixture-integrity": "verified", go, typescript, swift };
}

function scalarRoleContext(expectedRole, providedRole, scalarKind = "digest") {
  return { kind: "scalar-role", scalarKind, expectedRole, providedRole };
}

function padJsonToLength(value, length) {
  return utf8(`${value}${" ".repeat(length - Buffer.byteLength(value))}`);
}

function nestedJsonArray(depth) {
  let value = "0";
  for (let current = 1; current < depth; current += 1) {
    value = `[${value}]`;
  }
  return utf8(value);
}

function jsonNodeBoundary(over) {
  const root = {};
  let nestedObjectIndex = 0;
  for (let rootIndex = 0; rootIndex < 256; rootIndex += 1) {
    root[`r${rootIndex}`] = Array.from({ length: 15 }, () => {
      const wrapScalar = nestedObjectIndex < 256 || (over && nestedObjectIndex === 256);
      const child = { value: wrapScalar ? [0] : 0 };
      nestedObjectIndex += 1;
      return child;
    });
  }
  return root;
}

function jsonMemberBoundary(over) {
  return Array.from({ length: 17 }, (_, index) =>
    numberedObject(index === 16 ? 256 : 240 + (over && index === 0 ? 1 : 0)),
  );
}

function jsonElementBoundary(over) {
  return Array.from({ length: 16 }, (_, index) =>
    Array.from({ length: 255 + (over && index === 0 ? 1 : 0) }, () => 0),
  );
}

function numberedObject(count) {
  return Object.fromEntries(Array.from({ length: count }, (_, index) => [`k${index}`, 0]));
}

function jsonBytes(value) {
  return utf8(JSON.stringify(value));
}

function cborByteStringWithTotalLength(totalLength) {
  for (const headerLength of [1, 2, 3, 5]) {
    const payloadLength = totalLength - headerLength;
    if (payloadLength >= 0 && encodeCborArgument(2, payloadLength).length === headerLength) {
      return concatenate([encodeCborArgument(2, payloadLength), new Uint8Array(payloadLength)]);
    }
  }
  throw new Error(`cannot construct CBOR byte string of ${totalLength} total bytes`);
}

function nestedCborArray(depth) {
  let value = 0;
  for (let current = 1; current < depth; current += 1) {
    value = [value];
  }
  return value;
}

function nestedCborMap(depth) {
  let value = 0;
  for (let current = 1; current < depth; current += 1) {
    value = new Map([[0, value]]);
  }
  return value;
}

function planItemBoundary(over) {
  const value = new Map();
  for (let index = 0; index < 64; index += 1) {
    const elementCount = over || index < 63 ? 2 : 1;
    value.set(
      index,
      Array.from({ length: elementCount }, () => 0),
    );
  }
  return value;
}

function registrationItemBoundary(over) {
  const value = numberedMap(16);
  if (over) {
    value.set(0, [0]);
  }
  return value;
}

function numberedMap(count) {
  return new Map(Array.from({ length: count }, (_, index) => [index, 0]));
}

function cborEncode(value) {
  if (value instanceof Uint8Array) {
    return concatenate([encodeCborArgument(2, value.length), value]);
  }
  if (typeof value === "number") {
    if (!Number.isSafeInteger(value)) {
      throw new Error(`CBOR integer is outside the safe range: ${value}`);
    }
    return value >= 0 ? encodeCborArgument(0, value) : encodeCborArgument(1, -1 - value);
  }
  if (typeof value === "string") {
    const bytes = utf8(value);
    return concatenate([encodeCborArgument(3, bytes.length), bytes]);
  }
  if (Array.isArray(value)) {
    return concatenate([
      encodeCborArgument(4, value.length),
      ...value.map((child) => cborEncode(child)),
    ]);
  }
  if (value instanceof Map) {
    const entries = [...value].map(([key, child]) => [cborEncode(key), cborEncode(child)]);
    entries.sort(([left], [right]) =>
      left.length === right.length ? Buffer.compare(left, right) : left.length - right.length,
    );
    return concatenate([encodeCborArgument(5, entries.length), ...entries.flat()]);
  }
  throw new Error(`unsupported CBOR fixture value: ${String(value)}`);
}

function encodeCborArgument(majorType, argument) {
  if (!Number.isSafeInteger(argument) || argument < 0) {
    throw new Error(`CBOR argument is outside the safe unsigned range: ${argument}`);
  }
  if (argument < 24) {
    return Uint8Array.of((majorType << 5) | argument);
  }
  if (argument <= 0xff) {
    return Uint8Array.of((majorType << 5) | 24, argument);
  }
  if (argument <= 0xffff) {
    return Uint8Array.of((majorType << 5) | 25, argument >> 8, argument & 0xff);
  }
  if (argument <= 0xffffffff) {
    return Uint8Array.of(
      (majorType << 5) | 26,
      Math.floor(argument / 0x1000000) & 0xff,
      Math.floor(argument / 0x10000) & 0xff,
      Math.floor(argument / 0x100) & 0xff,
      argument & 0xff,
    );
  }
  if (argument <= Number.MAX_SAFE_INTEGER) {
    const bytes = new Uint8Array(9);
    bytes[0] = (majorType << 5) | 27;
    let remaining = argument;
    for (let index = 8; index >= 1; index -= 1) {
      bytes[index] = remaining % 256;
      remaining = Math.floor(remaining / 256);
    }
    return bytes;
  }
  throw new Error(`CBOR fixture argument is too large: ${argument}`);
}

function utf8(value) {
  return Buffer.from(value, "utf8");
}

function concatenate(parts) {
  return Buffer.concat(parts.map((part) => Buffer.from(part)));
}
