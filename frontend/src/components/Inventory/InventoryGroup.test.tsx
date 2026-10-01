import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { InventoryGroup } from "./InventoryGroup";
import { makeBagItem } from "../../test/utils";

describe("InventoryGroup", () => {
  it("uses the singular for one item and omits section labels", () => {
    render(
      <InventoryGroup
        characterName="Bank"
        characterInventory={[makeBagItem()]}
      />,
    );

    expect(screen.getByRole("heading", { name: "Bank" })).toBeInTheDocument();
    expect(screen.getByText("1 item")).toBeInTheDocument();
    expect(screen.queryByRole("heading", { level: 3 })).not.toBeInTheDocument();
  });

  it("labels equipment and inventory when a group has both", () => {
    render(
      <InventoryGroup
        characterName="Warrior"
        characterInventory={[makeBagItem(), makeBagItem()]}
        equipment={[makeBagItem()]}
      />,
    );

    expect(screen.getByText("3 items")).toBeInTheDocument();
    expect(
      screen.getAllByRole("heading", { level: 3 }).map((h) => h.textContent),
    ).toEqual(["Equipment", "Inventory"]);
  });
});
