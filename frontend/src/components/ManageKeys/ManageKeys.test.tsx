import { describe, expect, it, vi } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ManageKeys } from "./ManageKeys";
import { makeAccount, renderWithClient } from "../../test/utils";

describe("ManageKeys", () => {
  it("shows loading, then the account's key", async () => {
    renderWithClient(<ManageKeys />, {
      getAccount: vi.fn().mockResolvedValue(makeAccount()),
    });

    expect(screen.getByText("Loading…")).toBeInTheDocument();
    expect(await screen.findByText("Tester.1234")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "Remove key" }),
    ).toBeInTheDocument();
    expect(screen.queryByLabelText("Add an API key")).not.toBeInTheDocument();
  });

  it("disables adding until a non-blank key is entered", async () => {
    const user = userEvent.setup();
    renderWithClient(<ManageKeys />, {
      getAccount: vi.fn().mockResolvedValue(null),
    });

    const input = await screen.findByLabelText("Add an API key");
    const button = screen.getByRole("button", { name: "Add key" });
    expect(button).toBeDisabled();

    await user.type(input, "   ");
    expect(button).toBeDisabled();

    await user.type(input, "KEY");
    expect(button).toBeEnabled();
  });

  it("adds a trimmed key and shows the new account", async () => {
    const user = userEvent.setup();
    const postAPIKey = vi.fn().mockResolvedValue(makeAccount());
    renderWithClient(<ManageKeys />, {
      getAccount: vi.fn().mockResolvedValue(null),
      postAPIKey,
    });

    await user.type(await screen.findByLabelText("Add an API key"), "  KEY  ");
    await user.click(screen.getByRole("button", { name: "Add key" }));

    expect(postAPIKey).toHaveBeenCalledWith("KEY");
    expect(await screen.findByText("Tester.1234")).toBeInTheDocument();
  });

  it("shows an error for a rejected key and clears it on edit", async () => {
    const user = userEvent.setup();
    renderWithClient(<ManageKeys />, {
      getAccount: vi.fn().mockResolvedValue(null),
      postAPIKey: vi.fn().mockResolvedValue(null),
    });

    const input = await screen.findByLabelText("Add an API key");
    await user.type(input, "BAD");
    await user.click(screen.getByRole("button", { name: "Add key" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Couldn't add that key.",
    );
    expect(input).toHaveAttribute("aria-invalid", "true");
    expect(input).toHaveAccessibleDescription(/Couldn't add that key/);

    await user.type(input, "X");
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
    expect(input).toHaveAttribute("aria-invalid", "false");
  });

  it("shows an error when posting the key throws", async () => {
    const user = userEvent.setup();
    vi.spyOn(console, "error").mockImplementation(() => {});
    renderWithClient(<ManageKeys />, {
      getAccount: vi.fn().mockResolvedValue(null),
      postAPIKey: vi.fn().mockRejectedValue(new Error("boom")),
    });

    await user.type(await screen.findByLabelText("Add an API key"), "KEY");
    await user.click(screen.getByRole("button", { name: "Add key" }));

    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Add key" })).toBeEnabled();
  });
});
