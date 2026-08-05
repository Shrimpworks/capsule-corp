import { createHash } from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";

const root = new URL("../schemas/conformance/authority-plane-v0/", import.meta.url);
const corpus = new URL("../schemas/conformance/v0/", import.meta.url);
const plan = await readFile(new URL("execution-plan/ordinary.cbor", corpus));
const sourceManifest = await readFile(new URL("source-manifest/ordinary.cbor", corpus));
const proposal = JSON.parse(await readFile(new URL("job-proposal/ordinary.json", corpus), "utf8"));
const source = Buffer.from(proposal.source.files["main.mjs"], "utf8");

const fields = [
  Buffer.alloc(16, 0x11),
  Buffer.alloc(32, 0x22),
  Buffer.from("c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14", "hex"),
  Buffer.from("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72", "hex"),
  Buffer.alloc(32, 0x55),
];
const reviews = [Buffer.alloc(32, 0x66), Buffer.alloc(32, 0x67)];
const roleBindings = Buffer.concat([
  Buffer.from([0]),
  ...fields,
  Buffer.from([reviews.length]),
  ...reviews,
  ...Array.from({ length: 8 - reviews.length }, () => Buffer.alloc(32)),
  Buffer.alloc(32, 0x77),
  Buffer.alloc(32, 0x88),
  Buffer.alloc(32, 0x99),
  Buffer.alloc(32, 0xaa),
  Buffer.alloc(32, 0xbb),
]);
if (roleBindings.length !== 562) throw new Error("role-binding layout drift");

const caps = Object.freeze({
  executionPlan: 65_536,
  roleBindings: roleBindings.length,
  sourceManifest: 95,
  source: 262_144,
  planRegistration: 4_096,
});
const registerPlanV0 = caps.executionPlan + caps.roleBindings + caps.sourceManifest + caps.source;
const getRegisteredPlanV0 = registerPlanV0 + caps.planRegistration;
if (registerPlanV0 !== 328_337 || getRegisteredPlanV0 !== 332_433) {
  throw new Error("authority-plane aggregate cap drift");
}
const manifest = Buffer.from(
  `${JSON.stringify(
    {
      objectType: "capsule.authority-plane-conformance",
      objectVersion: 0,
      roleBindingRecordVersion: 0,
      caps: { ...caps, registerPlanV0, getRegisteredPlanV0 },
      fixtures: Object.fromEntries(
        [
          ["execution-plan.cbor", plan],
          ["role-bindings.bin", roleBindings],
          ["source-manifest.cbor", sourceManifest],
          ["main.mjs", source],
        ].map(([path, bytes]) => [path, { byteLength: bytes.length, sha256: sha256(bytes) }]),
      ),
    },
    null,
    2,
  )}\n`,
);
const expected = new Map([
  ["execution-plan.cbor", plan],
  ["role-bindings.bin", roleBindings],
  ["source-manifest.cbor", sourceManifest],
  ["main.mjs", source],
  ["manifest.json", manifest],
]);

if (process.argv.includes("--check")) {
  for (const [path, bytes] of expected) {
    const actual = await readFile(new URL(path, root));
    if (!actual.equals(bytes)) throw new Error(`stale authority-plane fixture: ${path}`);
  }
  process.stdout.write("verified authority-plane cross-language known answers\n");
} else {
  await mkdir(root, { recursive: true });
  for (const [path, bytes] of expected) await writeFile(new URL(path, root), bytes);
  process.stdout.write("generated authority-plane cross-language known answers\n");
}

function sha256(value) {
  return createHash("sha256").update(value).digest("hex");
}
