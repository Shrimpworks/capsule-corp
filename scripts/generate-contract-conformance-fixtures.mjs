import { createHash } from "node:crypto";
import { mkdir, writeFile } from "node:fs/promises";

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
const scalarImplementations = implementationStatus("pending", "pending", "pending");
const rules = [];
const cases = [];
const fixtures = new Map();

addMediaTypeRulesAndCases();
addJsonRulesAndCases();
addScalarRulesAndCases();

await mkdir(sharedDirectory, { recursive: true });
for (const [path, bytes] of [...fixtures].sort(([left], [right]) => left.localeCompare(right))) {
  const destination = new URL(path, corpusRoot);
  await mkdir(new URL("./", destination), { recursive: true });
  await writeFile(destination, bytes);
}
await writeFile(
  new URL("manifest.json", corpusRoot),
  `${JSON.stringify({ manifestVersion: "capsule.conformance/v0", rules, cases }, null, 2)}\n`,
);
process.stdout.write(
  `generated conformance corpus: ${rules.length} rules, ${cases.length} cases, ${fixtures.size} fixtures\n`,
);

function addMediaTypeRulesAndCases() {
  const profiles = [
    ["job-proposal", "JobProposal", jobProposalMediaType],
    ["execution-plan", "ExecutionPlan", executionPlanMediaType],
    ["plan-registration", "PlanRegistration", planRegistrationMediaType],
  ];
  for (const [slug, object, mediaType] of profiles) {
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
      implementations: rawImplementations,
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
      implementations: rawImplementations,
    });
    addMediaTypeReject(slug, object, ruleId, mediaType.split(";")[0], "missing-version", malformed);
    addMediaTypeReject(
      slug,
      object,
      ruleId,
      `${mediaType};charset=utf-8`,
      "additional-parameter",
      malformed,
    );
    addMediaTypeReject(
      slug,
      object,
      ruleId,
      mediaType.replace("v=0", "v=1"),
      "unknown-version",
      "UNSUPPORTED",
    );
  }
}

function addMediaTypeReject(slug, object, ruleId, mediaType, suffix, classification) {
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
    implementations: rawImplementations,
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
    implementations: rawImplementations,
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
}) {
  const retainedBytes = typeof bytes === "string" ? utf8(bytes) : Buffer.from(bytes);
  const existing = fixtures.get(path);
  if (existing && !existing.equals(retainedBytes)) {
    throw new Error(`fixture path ${path} was assigned different bytes`);
  }
  fixtures.set(path, retainedBytes);
  cases.push({
    id,
    description,
    ruleIds,
    object,
    wireFormat,
    mediaType,
    variant,
    fixture: {
      path,
      sha256: createHash("sha256").update(retainedBytes).digest("hex"),
      byteLength: retainedBytes.length,
    },
    context,
    expected: { decision, classification, owner, authorityStateChanged: false },
    implementations,
  });
}

function implementationStatus(go, typescript, swift) {
  return { "fixture-integrity": "verified", go, typescript, swift };
}

function scalarRoleContext(expectedRole, providedRole) {
  return { kind: "scalar-role", scalarKind: "digest", expectedRole, providedRole };
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

function utf8(value) {
  return Buffer.from(value, "utf8");
}

function concatenate(parts) {
  return Buffer.concat(parts.map((part) => Buffer.from(part)));
}
