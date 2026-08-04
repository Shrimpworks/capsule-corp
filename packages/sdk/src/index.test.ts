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

test("health rejects a malformed daemon response", async () => {
  const client = new CapsuleClient({
    fetch: async () => Response.json({ status: 42 }),
  });

  await assert.rejects(() => client.health(), /malformed health response/);
});

test("version accepts a well-formed daemon response", async () => {
  const client = new CapsuleClient({
    fetch: async () => Response.json({ version: "0.0.0", commit: "abc123", buildDate: "unknown" }),
  });

  assert.deepEqual(await client.version(), {
    version: "0.0.0",
    commit: "abc123",
    buildDate: "unknown",
  });
});

test("version rejects a malformed daemon response", async () => {
  const client = new CapsuleClient({
    fetch: async () => Response.json({ version: "0.0.0" }),
  });

  await assert.rejects(() => client.version(), /malformed version response/);
});

test("listRuntimes accepts a well-formed daemon response", async () => {
  const profile = { name: "bun-data@1", runtime: "bun", status: "draft", available: false };
  const client = new CapsuleClient({
    fetch: async () => Response.json({ profiles: [profile] }),
  });

  assert.deepEqual(await client.listRuntimes(), [profile]);
});

test("listRuntimes rejects a response with a malformed profile entry", async () => {
  const client = new CapsuleClient({
    fetch: async () =>
      Response.json({ profiles: [{ name: "bun-data@1", runtime: "unsupported-runtime" }] }),
  });

  await assert.rejects(() => client.listRuntimes(), /malformed runtimes response/);
});

test("listRuntimes rejects a response whose profiles field is not an array", async () => {
  const client = new CapsuleClient({
    fetch: async () => Response.json({ profiles: "not-an-array" }),
  });

  await assert.rejects(() => client.listRuntimes(), /malformed runtimes response/);
});

test("a baseUrl path prefix without a trailing slash is preserved, not dropped", async () => {
  let requestedUrl: string | undefined;
  const client = new CapsuleClient({
    baseUrl: "http://localhost:7777/v2",
    fetch: async (input) => {
      requestedUrl = input.toString();
      return Response.json({ status: "ok" });
    },
  });

  await client.health();

  assert.equal(requestedUrl, "http://localhost:7777/v2/healthz");
});
