import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { useRouteError } from "react-router-dom";
import { ErrorPage } from "./ErrorPage";

// Data-router loaders trip over jsdom's AbortSignal, so stub the hook directly
vi.mock("react-router-dom", async (importOriginal) => ({
  ...(await importOriginal<typeof import("react-router-dom")>()),
  useRouteError: vi.fn(),
}));

const routeError = (status: number, statusText: string, data: unknown) => ({
  status,
  statusText,
  internal: false,
  data,
});

describe("ErrorPage", () => {
  beforeEach(() => {
    vi.mocked(useRouteError).mockReset();
  });

  it("shows the message of a thrown Error", () => {
    vi.mocked(useRouteError).mockReturnValue(new Error("Something broke"));
    render(<ErrorPage />);

    expect(screen.getByRole("heading")).toHaveTextContent("Oopsie Woopsie");
    expect(screen.getByText("Something broke")).toBeInTheDocument();
  });

  it("shows a thrown string", () => {
    vi.mocked(useRouteError).mockReturnValue("plain string");
    render(<ErrorPage />);

    expect(screen.getByText("plain string")).toBeInTheDocument();
  });

  it("prefers the message in a route error response", () => {
    vi.mocked(useRouteError).mockReturnValue(
      routeError(404, "Not Found", { message: "Missing item" }),
    );
    render(<ErrorPage />);

    expect(screen.getByText("Missing item")).toBeInTheDocument();
  });

  it("falls back to the route error's status text", () => {
    vi.mocked(useRouteError).mockReturnValue(
      routeError(404, "Not Found", "Error: No route matches URL"),
    );
    render(<ErrorPage />);

    expect(screen.getByText("Not Found")).toBeInTheDocument();
  });

  it("shows a generic message for unknown errors", () => {
    const consoleError = vi
      .spyOn(console, "error")
      .mockImplementation(() => {});
    vi.mocked(useRouteError).mockReturnValue(42);
    render(<ErrorPage />);

    expect(screen.getByText("Unknown error")).toBeInTheDocument();
    expect(consoleError).toHaveBeenCalledWith(42);
  });
});
