import React from "react";
import { AccountInventory } from "../../models/AccountInventory";
import { InventoryGroup } from "./InventoryGroup";
import { TooltipProvider } from "./TooltipProvider";
import content from "../content.module.css";
import inventory from "./inventory.module.css";

interface AccountInventoryProps {
  accountInventory: AccountInventory;
}

export const InventoryContents: React.FC<AccountInventoryProps> = ({
  accountInventory,
}) => {
  const { sharedInventory, bankInventory, materialsInventory, characters } =
    accountInventory;

  const nonEmptyCharacters = (characters ?? []).filter(
    (character) =>
      (character.inventory?.length ?? 0) + (character.equipment?.length ?? 0) >
      0,
  );

  const hasItems =
    !!sharedInventory?.length ||
    !!bankInventory?.length ||
    !!materialsInventory?.length ||
    nonEmptyCharacters.length > 0;

  if (!hasItems) {
    return (
      <div className={content.card}>
        <p className={inventory.emptyState}>No items found.</p>
      </div>
    );
  }

  return (
    <TooltipProvider>
      <div className={inventory.groups}>
        {!!sharedInventory?.length && (
          <InventoryGroup
            characterName="Shared Inventory"
            characterInventory={sharedInventory}
          ></InventoryGroup>
        )}
        {!!bankInventory?.length && (
          <InventoryGroup
            characterName="Bank"
            characterInventory={bankInventory}
          ></InventoryGroup>
        )}
        {!!materialsInventory?.length && (
          <InventoryGroup
            characterName="Materials"
            characterInventory={materialsInventory}
          ></InventoryGroup>
        )}
        {nonEmptyCharacters.map((character) => {
          return (
            <InventoryGroup
              key={character.name}
              characterName={character.name}
              characterInventory={character.inventory}
              equipment={character.equipment}
            ></InventoryGroup>
          );
        })}
      </div>
    </TooltipProvider>
  );
};
