import type { RuntimeProfileDescriptor } from "@capsule-corp/protocol";

export interface CapsuleClientOptions {
  baseUrl?: string;
  fetch?: typeof globalThis.fetch;
}

export interface VersionInfo {
  version: string;
  commit: string;
  buildDate: string;
}

export class CapsuleClient {
  readonly #baseUrl: URL;
  readonly #fetch: typeof globalThis.fetch;

  constructor(options: CapsuleClientOptions = {}) {
    const baseUrl = options.baseUrl ?? "http://127.0.0.1:7777/";
    // URL's relative resolution drops the base's final path segment unless
    // it ends in "/", silently mangling any baseUrl with a path prefix.
    this.#baseUrl = new URL(baseUrl.endsWith("/") ? baseUrl : `${baseUrl}/`);
    this.#fetch = options.fetch ?? globalThis.fetch;
  }

  async health(signal?: AbortSignal): Promise<"ok"> {
    const response = await this.#request("healthz", signal);
    const body: unknown = await response.json();
    assertHealthResponse(body);

    if (body.status !== "ok") {
      throw new Error(`unexpected daemon health status: ${body.status}`);
    }

    return "ok";
  }

  async version(signal?: AbortSignal): Promise<VersionInfo> {
    const response = await this.#request("v1/version", signal);
    const body: unknown = await response.json();
    assertVersionInfo(body);
    return body;
  }

  async listRuntimes(signal?: AbortSignal): Promise<RuntimeProfileDescriptor[]> {
    const response = await this.#request("v1/runtimes", signal);
    const body: unknown = await response.json();
    assertRuntimesResponse(body);
    return body.profiles;
  }

  async #request(path: string, signal?: AbortSignal): Promise<Response> {
    const response = await this.#fetch(new URL(path, this.#baseUrl), {
      headers: { Accept: "application/json" },
      method: "GET",
      signal,
    });

    if (!response.ok) {
      throw new Error(`capsule daemon returned HTTP ${response.status}`);
    }

    return response;
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

/**
 * Validates the daemon's `/healthz` body shape before the caller trusts it.
 * The daemon is a locally spawned process, not a fully trusted peer — a
 * malformed or unexpected-shape response should fail loudly here rather
 * than propagate an `undefined`/garbage value into caller code.
 */
function assertHealthResponse(value: unknown): asserts value is { status: string } {
  if (!isRecord(value) || typeof value.status !== "string") {
    throw new Error("capsule daemon returned a malformed health response");
  }
}

/** Validates the daemon's `/v1/version` body shape before the caller trusts it. */
function assertVersionInfo(value: unknown): asserts value is VersionInfo {
  if (
    !isRecord(value) ||
    typeof value.version !== "string" ||
    typeof value.commit !== "string" ||
    typeof value.buildDate !== "string"
  ) {
    throw new Error("capsule daemon returned a malformed version response");
  }
}

function isRuntimeProfileDescriptor(value: unknown): value is RuntimeProfileDescriptor {
  return (
    isRecord(value) &&
    typeof value.name === "string" &&
    (value.runtime === "bun" || value.runtime === "node" || value.runtime === "deno") &&
    (value.status === "draft" || value.status === "active" || value.status === "retired") &&
    typeof value.available === "boolean"
  );
}

/** Validates the daemon's `/v1/runtimes` body shape, including every array element, before the caller trusts it. */
function assertRuntimesResponse(
  value: unknown,
): asserts value is { profiles: RuntimeProfileDescriptor[] } {
  if (
    !isRecord(value) ||
    !Array.isArray(value.profiles) ||
    !value.profiles.every(isRuntimeProfileDescriptor)
  ) {
    throw new Error("capsule daemon returned a malformed runtimes response");
  }
}
