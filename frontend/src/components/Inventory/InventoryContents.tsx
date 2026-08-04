import React from "react";
import { AccountInventory } from "../../models/AccountInventory";
import { InventoryGroup } from "./InventoryGroup";
import { TooltipProvider } from "./TooltipContext";
import inventory from "./inventory.module.css";

interface AccountInventoryProps {
  accountInventory: AccountInventory;
}

export const InventoryContents: React.FC<AccountInventoryProps> = ({
  accountInventory,
}) => {
  const { sharedInventory, bankInventory, materialsInventory, characters } =
    accountInventory;

  return (
    <TooltipProvider>
      {sharedInventory && (
        <InventoryGroup
          characterName="Shared Inventory"
          characterInventory={sharedInventory}
        ></InventoryGroup>
      )}
      {bankInventory && (
        <InventoryGroup
          characterName="Bank"
          characterInventory={bankInventory}
        ></InventoryGroup>
      )}
      {materialsInventory && (
        <InventoryGroup
          characterName="Materials"
          characterInventory={materialsInventory}
        ></InventoryGroup>
      )}
      <div className={inventory.inventoryGroups}>
        {characters &&
          characters.map((character) => {
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
