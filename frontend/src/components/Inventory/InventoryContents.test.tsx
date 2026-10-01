import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { InventoryContents } from "./InventoryContents";
import { makeBagItem } from "../../test/utils";

describe("InventoryContents", () => {
  it("shows an empty state when there are no items", () => {
    render(
      <InventoryContents
        accountInventory={{
          accountID: "acct",
          sharedInventory: [],
          characters: [{ name: "Empty", inventory: [], equipment: [] }],
        }}
      />,
    );

    expect(screen.getByText("No items found.")).toBeInTheDocument();
  });

  it("renders only the non-empty groups", () => {
    render(
      <InventoryContents
        accountInventory={{
          accountID: "acct",
          sharedInventory: [makeBagItem()],
          bankInventory: [],
          materialsInventory: [makeBagItem()],
          characters: [
            { name: "Warrior", inventory: [makeBagItem()] },
            { name: "Guardian", equipment: [makeBagItem()] },
            { name: "Empty" },
          ],
        }}
      />,
    );

    const headings = screen
      .getAllByRole("heading", { level: 2 })
      .map((h) => h.textContent);
    expect(headings).toEqual([
      "Shared Inventory",
      "Materials",
      "Warrior",
      "Guardian",
    ]);
  });
});
