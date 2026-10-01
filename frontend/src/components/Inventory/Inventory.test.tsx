import { describe, expect, it, vi } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Inventory } from "./Inventory";
import { AccountInventory } from "../../models/AccountInventory";
import { makeBagItem, renderWithClient } from "../../test/utils";

const inventory: AccountInventory = {
  accountID: "acct",
  bankInventory: [makeBagItem({ name: "Bank Sword" })],
};

const searchResult: AccountInventory = {
  accountID: "acct",
  materialsInventory: [makeBagItem({ name: "Mithril Ore" })],
};

function deferred<T>() {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((r) => (resolve = r));
  return { promise, resolve };
}

describe("Inventory", () => {
  it("shows loading, then the inventory groups", async () => {
    renderWithClient(<Inventory />, {
      getAccountInventory: vi.fn().mockResolvedValue(inventory),
    });

    expect(screen.getByText("Loading…")).toBeInTheDocument();
    expect(
      await screen.findByRole("heading", { name: "Bank" }),
    ).toBeInTheDocument();
    expect(screen.getByAltText("Bank Sword")).toBeInTheDocument();
    expect(screen.queryByText("Loading…")).not.toBeInTheDocument();
  });

  it("shows the no-key page and no search when there is no inventory", async () => {
    renderWithClient(<Inventory />, {
      getAccountInventory: vi.fn().mockResolvedValue(null),
    });

    expect(
      await screen.findByRole("link", { name: "Manage keys" }),
    ).toHaveAttribute("href", "/manageKeys");
    expect(screen.queryByRole("search")).not.toBeInTheDocument();
  });

  it("searches and replaces the groups with the results", async () => {
    const user = userEvent.setup();
    const search = deferred<AccountInventory | null>();
    const postInventorySearch = vi.fn().mockReturnValue(search.promise);
    renderWithClient(<Inventory />, {
      getAccountInventory: vi.fn().mockResolvedValue(inventory),
      postInventorySearch,
    });

    await user.type(await screen.findByLabelText("Search items"), "ore");
    await user.click(screen.getByRole("button", { name: "Search" }));

    expect(postInventorySearch).toHaveBeenCalledWith("ore");
    expect(screen.getByRole("button", { name: "Searching…" })).toBeDisabled();

    search.resolve(searchResult);

    expect(
      await screen.findByRole("heading", { name: "Materials" }),
    ).toBeInTheDocument();
    expect(
      screen.queryByRole("heading", { name: "Bank" }),
    ).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Search" })).toBeEnabled();
  });

  it("keeps the current groups when the search fails", async () => {
    const user = userEvent.setup();
    vi.spyOn(console, "log").mockImplementation(() => {});
    renderWithClient(<Inventory />, {
      getAccountInventory: vi.fn().mockResolvedValue(inventory),
      postInventorySearch: vi.fn().mockResolvedValue(null),
    });

    await user.type(await screen.findByLabelText("Search items"), "ore");
    await user.click(screen.getByRole("button", { name: "Search" }));

    expect(await screen.findByRole("button", { name: "Search" })).toBeEnabled();
    expect(screen.getByRole("heading", { name: "Bank" })).toBeInTheDocument();
  });
});
