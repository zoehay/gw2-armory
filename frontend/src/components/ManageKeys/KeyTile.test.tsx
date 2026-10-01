import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { act, fireEvent, screen } from "@testing-library/react";
import { KeyTile } from "./KeyTile";
import { makeAccount, renderWithClient } from "../../test/utils";

beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

describe("KeyTile", () => {
  it("renders the account fields with fallbacks", () => {
    renderWithClient(
      <KeyTile
        account={makeAccount({ gw2TokenName: undefined })}
        handleUpdate={vi.fn()}
      />,
    );

    expect(screen.getByText("—")).toBeInTheDocument();
    expect(screen.getByText("Tester.1234")).toBeInTheDocument();
    expect(screen.getByText("account-id")).toHaveAttribute(
      "title",
      "account-id",
    );
  });

  it("deletes the key after a delay and clears the account", async () => {
    const deleteAPIKey = vi.fn().mockResolvedValue("API-KEY");
    const handleUpdate = vi.fn();
    renderWithClient(
      <KeyTile account={makeAccount()} handleUpdate={handleUpdate} />,
      { deleteAPIKey },
    );

    fireEvent.click(screen.getByRole("button", { name: "Remove key" }));
    expect(screen.getByRole("button", { name: "Removing…" })).toBeDisabled();
    expect(deleteAPIKey).not.toHaveBeenCalled();

    await act(() => vi.advanceTimersByTimeAsync(2000));

    expect(deleteAPIKey).toHaveBeenCalledWith("API-KEY");
    expect(handleUpdate).toHaveBeenCalledWith(null);
  });

  it("re-enables the button when the delete fails", async () => {
    const handleUpdate = vi.fn();
    renderWithClient(
      <KeyTile account={makeAccount()} handleUpdate={handleUpdate} />,
      { deleteAPIKey: vi.fn().mockResolvedValue(null) },
    );

    fireEvent.click(screen.getByRole("button", { name: "Remove key" }));
    await act(() => vi.advanceTimersByTimeAsync(2000));

    expect(handleUpdate).not.toHaveBeenCalled();
    expect(screen.getByRole("button", { name: "Remove key" })).toBeEnabled();
  });
});
