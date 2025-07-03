import { BagItem } from "../../models/BagItem";
import React, { useState } from "react";
import inventory from "./inventory.module.css";

export interface InventoryTileProps {
  bagItem: BagItem;
}

export const InventoryTile: React.FC<InventoryTileProps> = ({ bagItem }) => {
  let [displayDetails, setDisplayDetails] = useState(false);
  let style = {
    borderColor: `var(--outline-${bagItem.rarity})`,
  } as React.CSSProperties;

  const handleMouseEnter = () => {
    setDisplayDetails(true);
  };

  const handleMouseLeave = () => {
    setDisplayDetails(false);
  };

  const handleTapToggle = () => {
    setDisplayDetails(!displayDetails);
  };

  return (
    <div
      className={inventory.tile}
      onMouseEnter={handleMouseEnter}
      onClick={handleTapToggle}
      onMouseLeave={handleMouseLeave}
      style={style}
    >
      {bagItem.count > 1 && (
        <div className={inventory.count}>{bagItem.count}</div>
      )}
      <img className={inventory.icon} src={bagItem.icon} alt={bagItem.name} />
      {displayDetails && <ToolTip bagItem={bagItem}></ToolTip>}
    </div>
  );
};

const ToolTip: React.FC<InventoryTileProps> = ({ bagItem }) => {
  return (
    <div className={inventory.tooltip}>
      <p className={inventory.name}>{bagItem.name}</p>
      <p className={inventory.description}>{bagItem.description}</p>
    </div>
  );
};
