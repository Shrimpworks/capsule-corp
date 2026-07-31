import assert from "node:assert/strict";
import test from "node:test";

import { CapsuleClient } from "./index.js";

test("health accepts an ok daemon response", async () => {
  const client = new CapsuleClient({
    fetch: async () => Response.json({ status: "ok" }),
  });

  assert.equal(await client.health(), "ok");
});

test("health rejects an unexpected daemon response", async () => {
  const client = new CapsuleClient({
    fetch: async () => Response.json({ status: "degraded" }),
  });

  await assert.rejects(() => client.health(), /unexpected daemon health status/);
});
