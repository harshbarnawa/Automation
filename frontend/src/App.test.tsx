import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { App } from "./App";

describe("App", () => {
  it("renders the Mintok product shell", () => {
    render(<App />);

    expect(screen.getByText("Mintok")).toBeInTheDocument();
    expect(
      screen.getByText(/AI Token Optimization Gateway/i),
    ).toBeInTheDocument();
  });
});
