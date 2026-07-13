import { describe, expect, it } from "vitest";

import { config } from "./config";

describe("config", () => {
  it("exposes the API base URL", () => {
    expect(config.apiBaseUrl).toBe("http://localhost:8080");
  });
});
