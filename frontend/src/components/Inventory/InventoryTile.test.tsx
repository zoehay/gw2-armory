import { describe, expect, it } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { InventoryTile } from "./InventoryTile";
import { TooltipProvider } from "./TooltipProvider";
import { BagItem } from "../../models/BagItem";
import { makeBagItem } from "../../test/utils";

function renderTiles(...items: BagItem[]) {
  render(
    <TooltipProvider>
      {items.map((item, i) => (
        <InventoryTile key={i} bagItem={item} />
      ))}
      <p>Outside</p>
    </TooltipProvider>,
  );
}

// The item name appears in the tooltip as text; the tile only uses it as alt text
const tooltipFor = (name: string) => screen.queryByText(name);

describe("InventoryTile", () => {
  it("renders the icon and shows a count badge only for stacks", () => {
    renderTiles(
      makeBagItem({ name: "Single", count: 1 }),
      makeBagItem({ name: "Stack", count: 250 }),
    );

    expect(screen.getByAltText("Single")).toHaveAttribute(
      "src",
      "https://example.com/icon.png",
    );
    expect(screen.getByText("250")).toBeInTheDocument();
    expect(screen.queryByText("1")).not.toBeInTheDocument();
  });

  it("shows the tooltip on hover and hides it on leave", async () => {
    const user = userEvent.setup();
    renderTiles(makeBagItem({ name: "Sword" }));
    const tile = screen.getByAltText("Sword");

    await user.hover(tile);
    expect(tooltipFor("Sword")).toBeInTheDocument();

    await user.unhover(tile);
    expect(tooltipFor("Sword")).not.toBeInTheDocument();
  });

  it("toggles the tooltip on tap", () => {
    renderTiles(makeBagItem({ name: "Sword" }));
    const tile = screen.getByAltText("Sword");

    fireEvent.click(tile);
    expect(tooltipFor("Sword")).toBeInTheDocument();

    fireEvent.click(tile);
    expect(tooltipFor("Sword")).not.toBeInTheDocument();
  });

  it("shows only one tooltip at a time", () => {
    renderTiles(
      makeBagItem({ name: "Sword" }),
      makeBagItem({ name: "Shield" }),
    );

    fireEvent.click(screen.getByAltText("Sword"));
    fireEvent.click(screen.getByAltText("Shield"));

    expect(tooltipFor("Sword")).not.toBeInTheDocument();
    expect(tooltipFor("Shield")).toBeInTheDocument();
  });

  it("closes the tooltip when touching outside a tile", () => {
    renderTiles(makeBagItem({ name: "Sword" }));

    fireEvent.click(screen.getByAltText("Sword"));
    fireEvent.touchStart(screen.getByAltText("Sword"));
    expect(tooltipFor("Sword")).toBeInTheDocument();

    fireEvent.touchStart(screen.getByText("Outside"));
    expect(tooltipFor("Sword")).not.toBeInTheDocument();
  });

  describe("tooltip content", () => {
    function openTooltip(item: BagItem) {
      renderTiles(item);
      fireEvent.click(screen.getByAltText(item.name!));
    }

    it("shows defense, weapon strength, and labelled infix attributes", () => {
      openTooltip(
        makeBagItem({
          name: "Greatsword",
          type: "Weapon",
          details: {
            defense: 12,
            min_power: 1000,
            max_power: 1100,
            infix_upgrade: {
              attributes: [
                { attribute: "Power", modifier: 179 },
                { attribute: "CritDamage", modifier: 128 },
                { attribute: "Mystery", modifier: 5 },
              ],
            },
          },
          stats: { attributes: { Vitality: 1 } },
        }),
      );

      expect(screen.getByText("Defense 12")).toBeInTheDocument();
      expect(
        screen.getByText("Weapon Strength 1000 - 1100"),
      ).toBeInTheDocument();
      expect(screen.getByText("+179 Power")).toBeInTheDocument();
      expect(screen.getByText("+128 Ferocity")).toBeInTheDocument();
      expect(screen.getByText("+5 Mystery")).toBeInTheDocument();
      // infix attributes take precedence over stats
      expect(screen.queryByText("+1 Vitality")).not.toBeInTheDocument();
      expect(screen.getByText("Weapon")).toBeInTheDocument();
    });

    it("omits zero defense and falls back to stats attributes", () => {
      openTooltip(
        makeBagItem({
          name: "Coat",
          details: { defense: 0 },
          stats: { id: 1, attributes: { Healing: 60, BoonDuration: 40 } },
        }),
      );

      expect(screen.queryByText(/Defense/)).not.toBeInTheDocument();
      expect(screen.getByText("+60 Healing Power")).toBeInTheDocument();
      expect(screen.getByText("+40 Concentration")).toBeInTheDocument();
    });

    it("lists upgrades with their bonuses", () => {
      openTooltip(
        makeBagItem({
          name: "Helm",
          upgradeDetails: [
            {
              id: 24836,
              name: "Superior Rune of the Scholar",
              icon: "rune.png",
              rarity: "Exotic",
              details: { bonuses: ["+25 Power", "+35 Ferocity"] },
            },
          ],
        }),
      );

      expect(
        screen.getByText("Superior Rune of the Scholar"),
      ).toBeInTheDocument();
      expect(screen.getByText("+25 Power")).toBeInTheDocument();
      expect(screen.getByText("+35 Ferocity")).toBeInTheDocument();
    });

    it("renders color tags in the description as spans", () => {
      openTooltip(
        makeBagItem({
          name: "Tome",
          description: "Read me.<br><c=@flavor>Ancient words</c> end",
        }),
      );

      const flavor = screen.getByText("Ancient words");
      expect(flavor.tagName).toBe("SPAN");
      expect(flavor.parentElement).toHaveTextContent(
        "Read me. Ancient words end",
      );
      expect(flavor.parentElement?.innerHTML).not.toContain("<br");
    });
  });
});
