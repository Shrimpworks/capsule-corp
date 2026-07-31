import assert from "node:assert/strict";
import test from "node:test";

import { createDenyByDefaultCapabilities } from "./index.js";

test("deny-by-default capabilities grant no ambient authority", () => {
  assert.deepEqual(createDenyByDefaultCapabilities(), {
    network: [],
    subprocesses: [],
    environment: [],
    nativeAddons: false,
    ffi: false,
    packageInstallation: false,
  });
});
