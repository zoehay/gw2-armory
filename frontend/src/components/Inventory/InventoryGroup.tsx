import React from "react";
import inventory from "./inventory.module.css";
import { BagItem } from "../../models/BagItem";
import { InventoryTile } from "./InventoryTile";

interface InventoryGroupProps {
  characterName: string;
  characterInventory?: BagItem[];
  equipment?: BagItem[];
}

export const InventoryGroup: React.FC<InventoryGroupProps> = ({
  characterName,
  equipment,
  characterInventory,
}) => {
  let itemTiles;
  if (characterInventory) {
    itemTiles = characterInventory.map((item, index) => (
      <InventoryTile bagItem={item} key={index} />
    ));
  }

  let equipmentTiles;
  if (equipment) {
    equipmentTiles = equipment.map((item, index) => (
      <InventoryTile bagItem={item} key={index} />
    ));
  }

  return (
    <div className={inventory.group}>
      <div className={inventory.name}>{characterName}</div>
      <div className={inventory.contents}>{equipmentTiles}</div>
      <div className={inventory.contents}>{itemTiles}</div>
    </div>
  );
};
