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
    this.#baseUrl = new URL(options.baseUrl ?? "http://127.0.0.1:7777/");
    this.#fetch = options.fetch ?? globalThis.fetch;
  }

  async health(signal?: AbortSignal): Promise<"ok"> {
    const response = await this.#request("healthz", signal);
    const body = (await response.json()) as { status: string };

    if (body.status !== "ok") {
      throw new Error(`unexpected daemon health status: ${body.status}`);
    }

    return "ok";
  }

  async version(signal?: AbortSignal): Promise<VersionInfo> {
    const response = await this.#request("v1/version", signal);
    return (await response.json()) as VersionInfo;
  }

  async listRuntimes(signal?: AbortSignal): Promise<RuntimeProfileDescriptor[]> {
    const response = await this.#request("v1/runtimes", signal);
    const body = (await response.json()) as { profiles: RuntimeProfileDescriptor[] };
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
