import React from "react";
import { AccountInventory } from "../../models/AccountInventory";
import { InventoryGroup } from "./InventoryGroup";
import inventory from "./inventory.module.css";

interface AccountInventoryProps {
  accountInventory: AccountInventory;
}

export const InventoryContents: React.FC<AccountInventoryProps> = ({
  accountInventory,
}) => {
  let { sharedInventory, characters } = accountInventory;

  return (
    <>
      {sharedInventory && (
        <InventoryGroup
          characterName="Shared Inventory"
          characterInventory={sharedInventory}
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
    </>
  );
};
