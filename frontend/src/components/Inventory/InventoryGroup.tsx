import React from "react";
import content from "../content.module.css";
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
  if (characterInventory?.length) {
    itemTiles = characterInventory.map((item, index) => (
      <InventoryTile bagItem={item} key={index} />
    ));
  }

  let equipmentTiles;
  if (equipment?.length) {
    equipmentTiles = equipment.map((item, index) => (
      <InventoryTile bagItem={item} key={index} />
    ));
  }

  // Only label sections when a group has both, i.e. characters
  const showLabels = !!equipmentTiles && !!itemTiles;
  const itemCount =
    (characterInventory?.length ?? 0) + (equipment?.length ?? 0);

  return (
    <section className={content.card}>
      <div className={inventory.groupHeader}>
        <h2 className={inventory.groupName}>{characterName}</h2>
        <span className={inventory.groupCount}>
          {itemCount} {itemCount === 1 ? "item" : "items"}
        </span>
      </div>
      {equipmentTiles && (
        <div className={inventory.section}>
          {showLabels && <h3 className={inventory.sectionLabel}>Equipment</h3>}
          <div className={inventory.contents}>{equipmentTiles}</div>
        </div>
      )}
      {itemTiles && (
        <div className={inventory.section}>
          {showLabels && <h3 className={inventory.sectionLabel}>Inventory</h3>}
          <div className={inventory.contents}>{itemTiles}</div>
        </div>
      )}
    </section>
  );
};
