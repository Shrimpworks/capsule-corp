import { sha256 } from "./lib/fixture-bytes.mjs";

const profile = "capsule.mjs-source/v0";
const memberMediaType = "application/capsule.javascript-source;v=0;module=esm";
const manifestMediaType = "application/capsule.source-manifest+cbor;v=0";
const mainPath = "main.mjs";
const verified = {
  "fixture-integrity": "verified",
  go: "verified",
  typescript: "verified",
  swift: "pending",
};
const languageHold = {
  "fixture-integrity": "verified",
  go: "pending",
  typescript: "pending",
  swift: "pending",
};
const ordinaryRequirements = [
  { decision: "accept", variant: "ordinary" },
  { decision: "reject", variant: "malformed" },
];
const boundaryRequirements = [
  { decision: "accept", variant: "exact-maximum" },
  { decision: "reject", variant: "cap-plus-one" },
];

export function addMjsSourceFoundationRulesAndCases({
  addCase,
  addRule,
  cborEncode,
  retainFixture,
}) {
  const ordinarySource = bytes('const message = "Capsule 🚀";\r\nexport default message;\n');
  const minimumSource = Buffer.alloc(0);
  const maximumSource = bytes("a".repeat(262_144));
  const ordinary = retainFixture("mjs-source/ordinary.mjs", ordinarySource);
  const minimum = retainFixture("mjs-source/minimum.mjs", minimumSource);
  const maximum = retainFixture("mjs-source/exact-maximum.mjs", maximumSource);

  addRule(
    "mjs-source.closed-identifiers",
    "ADR-0034#exact-source-profile",
    "The passive source foundation closes profile, member media, manifest media, path, and entrypoint values.",
    ordinaryRequirements,
  );
  for (const [id, value, wrong, path] of [
    ["profile", profile, "capsule.mjs-source/v1", "mjs-source/profile"],
    ["member-media", memberMediaType, "application/typescript", "mjs-source/member-media"],
    [
      "manifest-media",
      manifestMediaType,
      "application/capsule.source-manifest+cbor;v=1",
      "mjs-source-manifest/media",
    ],
  ]) {
    addCase({
      id: `mjs-source.closed-identifiers.${id}`,
      description: `Accept the exact ${id} value.`,
      ruleIds: ["mjs-source.closed-identifiers"],
      object: id === "manifest-media" ? "SourceManifest" : "MjsSource",
      wireFormat: "media-type",
      mediaType: value,
      variant: "ordinary",
      path: `${path}-exact.txt`,
      bytes: value,
      owner: id === "manifest-media" ? "source-manifest-validator" : "mjs-source-byte-validator",
      implementations: verified,
    });
    addCase({
      id: `mjs-source.closed-identifiers.${id}-reject`,
      description: `Reject the closed ${id} value.`,
      ruleIds: ["mjs-source.closed-identifiers"],
      object: id === "manifest-media" ? "SourceManifest" : "MjsSource",
      wireFormat: "media-type",
      mediaType: wrong,
      variant: "malformed",
      path: `${path}-reject.txt`,
      bytes: wrong,
      decision: "reject",
      classification: "UNSUPPORTED",
      owner: id === "manifest-media" ? "source-manifest-validator" : "mjs-source-byte-validator",
      implementations: verified,
    });
  }

  addRule(
    "mjs-source.byte-length",
    "ADR-0034#exact-source-profile",
    "Exact main.mjs bytes are 0 through 262144 bytes inclusive.",
    boundaryRequirements,
  );
  for (const [id, variant, sourceBytes, decision, classification] of [
    ["minimum", "minimum", minimumSource, "accept", null],
    ["ordinary", "ordinary", ordinarySource, "accept", null],
    ["exact-maximum", "exact-maximum", maximumSource, "accept", null],
    ["cap-plus-one", "cap-plus-one", bytes("a".repeat(262_145)), "reject", "SEMANTIC"],
  ]) {
    addCase({
      id: `mjs-source.byte-length.${id}`,
      description: `${decision === "accept" ? "Accept" : "Reject"} the ${id} source-byte case.`,
      ruleIds: ["mjs-source.byte-length"],
      object: "MjsSource",
      wireFormat: "raw-bytes",
      mediaType: memberMediaType,
      variant,
      path: `mjs-source/${id === "exact-maximum" ? "exact-maximum" : id}.mjs`,
      bytes: sourceBytes,
      decision,
      classification,
      owner: "mjs-source-byte-validator",
      implementations: verified,
    });
  }

  addRule(
    "mjs-source.utf8-bom",
    "ADR-0034#exact-source-profile",
    "main.mjs is strict UTF-8 and does not begin with a UTF-8 BOM.",
    ordinaryRequirements,
  );
  for (const [id, value, classification] of [
    ["invalid-utf8", Buffer.from([0x61, 0xff]), "MALFORMED"],
    ["unpaired-surrogate-utf8", Buffer.from([0xed, 0xa0, 0x80]), "MALFORMED"],
    ["leading-bom", Buffer.from([0xef, 0xbb, 0xbf, 0x61]), "MALFORMED"],
  ]) {
    addCase({
      id: `mjs-source.utf8-bom.${id}`,
      description: `Reject ${id.replaceAll("-", " ")} source bytes.`,
      ruleIds: ["mjs-source.utf8-bom"],
      object: "MjsSource",
      wireFormat: "raw-bytes",
      mediaType: memberMediaType,
      variant: "malformed",
      path: `mjs-source/${id}.mjs`,
      bytes: value,
      decision: "reject",
      classification,
      owner: "mjs-source-byte-validator",
      implementations: verified,
    });
  }
  addCase({
    id: "mjs-source.utf8-bom.ordinary",
    description: "Accept ordinary strict UTF-8 without a leading BOM.",
    ruleIds: ["mjs-source.utf8-bom"],
    object: "MjsSource",
    wireFormat: "raw-bytes",
    mediaType: memberMediaType,
    variant: "ordinary",
    path: ordinary.path,
    bytes: ordinarySource,
    owner: "mjs-source-byte-validator",
    implementations: verified,
  });

  addRule(
    "mjs-source.byte-identity",
    "ADR-0034#canonical-source-identity",
    "Newline, Unicode, embedded-BOM, and trailing-newline bytes are retained without rewriting.",
    ordinaryRequirements,
  );
  const identityValues = [
    ["lf", "const value = 'line';\n"],
    ["crlf", "const value = 'line';\r\n"],
    ["lone-cr", "const value = 'line';\r"],
    ["line-separator", "const value = 'a\u2028b';\n"],
    ["paragraph-separator", "const value = 'a\u2029b';\n"],
    ["composed", "const value = 'é';\n"],
    ["decomposed", "const value = 'e\u0301';\n"],
    ["embedded-bom", "const value = 'a\ufeffb';\n"],
    ["no-trailing-newline", "const value = 'é';"],
  ];
  for (const [id, value] of identityValues) {
    addCase({
      id: `mjs-source.byte-identity.${id}`,
      description: `Retain the ${id.replaceAll("-", " ")} bytes exactly.`,
      ruleIds: ["mjs-source.byte-identity"],
      object: "MjsSource",
      wireFormat: "raw-bytes",
      mediaType: memberMediaType,
      variant: "ordinary",
      path: `mjs-source/identity-${id}.mjs`,
      bytes: value,
      owner: "mjs-source-byte-validator",
      implementations: verified,
    });
  }
  const identityMutation = mutateDigest(sourceManifest(cborEncode, bytes(identityValues[0][1])));
  addCase({
    id: "mjs-source.byte-identity.digest-mutation",
    description: "Reject a manifest digest after a one-byte source identity change.",
    ruleIds: ["mjs-source.byte-identity"],
    object: "SourceManifest",
    wireFormat: "cbor",
    mediaType: manifestMediaType,
    variant: "malformed",
    path: "mjs-source-manifest/reject-identity-digest.cbor",
    bytes: identityMutation,
    context: {
      kind: "source-manifest",
      source: retainFixture("mjs-source/identity-lf.mjs", bytes(identityValues[0][1])),
    },
    decision: "reject",
    classification: "DOMAIN",
    owner: "source-manifest-validator",
    implementations: verified,
  });

  addSourceManifestCases({
    addCase,
    addRule,
    cborEncode,
    retainFixture,
    ordinary,
    ordinarySource,
    minimum,
    minimumSource,
    maximum,
    maximumSource,
  });
  addLanguageHoldCases({ addCase, addRule });
}

function addSourceManifestCases({
  addCase,
  addRule,
  cborEncode,
  retainFixture,
  ordinary,
  ordinarySource,
  minimum,
  minimumSource,
  maximum,
  maximumSource,
}) {
  addRule(
    "source-manifest.exact-cbor",
    "ADR-0034#canonical-source-identity",
    "SourceManifest v0 is the exact deterministic-CBOR single-member source identity.",
    boundaryRequirements,
  );
  const manifests = [
    ["minimum", "minimum", minimumSource, minimum],
    ["ordinary", "ordinary", ordinarySource, ordinary],
    ["exact-maximum", "exact-maximum", maximumSource, maximum],
  ];
  for (const [id, variant, sourceBytes, source] of manifests) {
    const manifest = sourceManifest(cborEncode, sourceBytes);
    addCase({
      id: `source-manifest.exact-cbor.${id}`,
      description: `Accept the exact ${id} SourceManifest.`,
      ruleIds: ["source-manifest.exact-cbor"],
      object: "SourceManifest",
      wireFormat: "cbor",
      mediaType: manifestMediaType,
      variant,
      path: `mjs-source-manifest/${id}.cbor`,
      bytes: manifest,
      context: { kind: "source-manifest", source },
      owner: "source-manifest-validator",
      implementations: verified,
    });
  }
  addCase({
    id: "source-manifest.exact-cbor.cap-plus-one-raw",
    description: "Reject a 96-byte SourceManifest raw encoding.",
    ruleIds: ["source-manifest.exact-cbor"],
    object: "SourceManifest",
    wireFormat: "cbor",
    mediaType: manifestMediaType,
    variant: "cap-plus-one",
    path: "mjs-source-manifest/cap-plus-one.cbor",
    bytes: Buffer.concat([sourceManifest(cborEncode, maximumSource), Buffer.from([0])]),
    context: { kind: "source-manifest", source: maximum },
    decision: "reject",
    classification: "MALFORMED",
    owner: "source-manifest-validator",
    implementations: verified,
  });
  addCase({
    id: "source-manifest.exact-cbor.cap-plus-one-source",
    description: "Reject source bytes above the source cap before manifest binding.",
    ruleIds: ["source-manifest.exact-cbor"],
    object: "SourceManifest",
    wireFormat: "cbor",
    mediaType: manifestMediaType,
    variant: "cap-plus-one",
    path: "mjs-source-manifest/ordinary.cbor",
    bytes: sourceManifest(cborEncode, ordinarySource),
    context: {
      kind: "source-manifest",
      source: retainFixture("mjs-source/cap-plus-one.mjs", bytes("a".repeat(262_145))),
    },
    decision: "reject",
    classification: "DOMAIN",
    owner: "source-manifest-validator",
    implementations: verified,
  });

  addRule(
    "source-manifest.closed-shape-binding",
    "ADR-0034#canonical-source-identity",
    "Object role, version, member count/path/digest/length, aggregate length, and map order are closed and source-bound.",
    ordinaryRequirements,
  );
  const digest = sha256(ordinarySource);
  const mutations = [
    [
      "object-type",
      sourceManifest(cborEncode, ordinarySource, { objectType: "capsule.execution-plan" }),
      "UNSUPPORTED",
    ],
    [
      "object-version",
      sourceManifest(cborEncode, ordinarySource, { objectVersion: 1 }),
      "UNSUPPORTED",
    ],
    [
      "entrypoint",
      sourceManifest(cborEncode, ordinarySource, { entrypoint: "other.mjs" }),
      "DOMAIN",
    ],
    [
      "member-path",
      sourceManifest(cborEncode, ordinarySource, { memberPath: "other.mjs" }),
      "DOMAIN",
    ],
    ["member-count", sourceManifest(cborEncode, ordinarySource, { members: [] }), "MALFORMED"],
    [
      "digest",
      sourceManifest(cborEncode, ordinarySource, { digest: Buffer.alloc(32, 0x44) }),
      "DOMAIN",
    ],
    [
      "cross-role-digest",
      sourceManifest(cborEncode, ordinarySource, {
        digest: sha256(sourceManifest(cborEncode, ordinarySource)),
      }),
      "DOMAIN",
    ],
    [
      "member-length",
      sourceManifest(cborEncode, ordinarySource, { memberLength: ordinarySource.length + 1 }),
      "DOMAIN",
    ],
    [
      "aggregate-length",
      sourceManifest(cborEncode, ordinarySource, { aggregateLength: ordinarySource.length + 1 }),
      "DOMAIN",
    ],
    ["map-order", sourceManifestWrongOrder(cborEncode, ordinarySource), "MALFORMED"],
  ];
  for (const [id, value, classification] of mutations) {
    addCase({
      id: `source-manifest.closed-shape-binding.${id}`,
      description: `Reject the ${id.replaceAll("-", " ")} mutation.`,
      ruleIds: ["source-manifest.closed-shape-binding"],
      object: "SourceManifest",
      wireFormat: "cbor",
      mediaType: manifestMediaType,
      variant: "malformed",
      path: `mjs-source-manifest/reject-${id}.cbor`,
      bytes: value,
      context: { kind: "source-manifest", source: ordinary },
      decision: "reject",
      classification,
      owner: "source-manifest-validator",
      implementations: verified,
    });
  }
  addCase({
    id: "source-manifest.closed-shape-binding.ordinary",
    description: "Accept the exact ordinary single-member manifest and source binding.",
    ruleIds: ["source-manifest.closed-shape-binding"],
    object: "SourceManifest",
    wireFormat: "cbor",
    mediaType: manifestMediaType,
    variant: "ordinary",
    path: "mjs-source-manifest/ordinary.cbor",
    bytes: sourceManifest(cborEncode, ordinarySource, { digest }),
    context: { kind: "source-manifest", source: ordinary },
    owner: "source-manifest-validator",
    implementations: verified,
  });
}

function addLanguageHoldCases({ addCase, addRule }) {
  addRule(
    "mjs-source.language-validator-hold",
    "ADR-0034#closed-module-loading-policy",
    "Live module requests and CommonJS must refuse, but the separately reviewed exact ECMAScript parser boundary is unimplemented.",
    ordinaryRequirements,
  );
  const cases = [
    ["property-import-meta", "obj.import.meta;\n", "accept", null],
    ["method-import", "({ import() {} });\n", "accept", null],
    [
      "template-interpolation-import",
      "const value = `$" + "{await import('./evil.mjs')}`;\n",
      "reject",
      "SEMANTIC",
    ],
    ["eval-string-data", "eval(\"import('./evil.mjs')\");\n", "accept", null],
    [
      "division-regexp-counterexample",
      'const of = 9; of / import("evil") / divisor;\n',
      "reject",
      "SEMANTIC",
    ],
    ["static-import", "import value from './evil.mjs';\n", "reject", "SEMANTIC"],
    ["side-effect-import", "import './evil.mjs';\n", "reject", "SEMANTIC"],
    ["export-from", "export { value } from './evil.mjs';\n", "reject", "SEMANTIC"],
    ["export-star", "export * from './evil.mjs';\n", "reject", "SEMANTIC"],
    ["import-meta", "const value = import.meta;\n", "reject", "SEMANTIC"],
    ["specifier-absolute", "import('/evil.mjs');\n", "reject", "SEMANTIC"],
    ["specifier-bare", "import('evil');\n", "reject", "SEMANTIC"],
    ["specifier-node", "import('node:fs');\n", "reject", "SEMANTIC"],
    ["specifier-npm", "import('npm:evil');\n", "reject", "SEMANTIC"],
    ["specifier-http", "import('http://invalid/evil.mjs');\n", "reject", "SEMANTIC"],
    ["specifier-https", "import('https://invalid/evil.mjs');\n", "reject", "SEMANTIC"],
    ["specifier-data", "import('data:text/javascript,export{}');\n", "reject", "SEMANTIC"],
    ["specifier-blob", "import('blob:opaque');\n", "reject", "SEMANTIC"],
    ["specifier-file", "import('file:///evil.mjs');\n", "reject", "SEMANTIC"],
    ["specifier-capsule", "import('capsule:evil');\n", "reject", "SEMANTIC"],
    ["commonjs-require", "const value = require('evil');\n", "reject", "SEMANTIC"],
    ["commonjs-require-resolve", "require.resolve('evil');\n", "reject", "SEMANTIC"],
    ["commonjs-module-exports", "module.exports = {};\n", "reject", "SEMANTIC"],
    ["commonjs-exports", "exports.value = 1;\n", "reject", "SEMANTIC"],
    ["commonjs-dirname", "const value = __dirname;\n", "reject", "SEMANTIC"],
    ["commonjs-filename", "const value = __filename;\n", "reject", "SEMANTIC"],
    ["local-export", "const value = 1; export { value };\n", "accept", null],
    [
      "noncode-spellings",
      "// import('comment')\nconst text = \"require module.exports import.meta\";\n",
      "accept",
      null,
    ],
  ];
  for (const [id, value, decision, classification] of cases) {
    addCase({
      id: `mjs-source.language-validator-hold.${id}`,
      description: `${decision === "accept" ? "Accept" : "Reject"} ${id.replaceAll("-", " ")} only after the parser boundary is implemented.`,
      ruleIds: ["mjs-source.language-validator-hold"],
      object: "MjsSource",
      wireFormat: "raw-bytes",
      mediaType: memberMediaType,
      variant: decision === "accept" ? "ordinary" : "malformed",
      path: `mjs-source/language-hold-${id}.mjs`,
      bytes: value,
      decision,
      classification,
      owner: "mjs-source-language-validator",
      implementations: languageHold,
    });
  }
}

function sourceManifest(cborEncode, sourceBytes, overrides = {}) {
  const digest = overrides.digest ?? sha256(sourceBytes);
  const members = overrides.members ?? [
    [overrides.memberPath ?? mainPath, digest, overrides.memberLength ?? sourceBytes.length],
  ];
  return cborEncode(
    new Map([
      [1, overrides.objectType ?? "capsule.source-manifest"],
      [2, overrides.objectVersion ?? 0],
      [3, overrides.entrypoint ?? mainPath],
      [4, members],
      [5, overrides.aggregateLength ?? sourceBytes.length],
    ]),
  );
}

function sourceManifestWrongOrder(cborEncode, sourceBytes) {
  return Buffer.concat([
    Buffer.from([0xa5]),
    cborEncode(2),
    cborEncode(0),
    cborEncode(1),
    cborEncode("capsule.source-manifest"),
    cborEncode(3),
    cborEncode(mainPath),
    cborEncode(4),
    cborEncode([[mainPath, sha256(sourceBytes), sourceBytes.length]]),
    cborEncode(5),
    cborEncode(sourceBytes.length),
  ]);
}

function mutateDigest(manifest) {
  const mutated = Buffer.from(manifest);
  const digestOffset = mutated.indexOf(Buffer.from([0x58, 0x20])) + 2;
  mutated[digestOffset] ^= 0x01;
  return mutated;
}

function bytes(value) {
  return Buffer.from(value, "utf8");
}
