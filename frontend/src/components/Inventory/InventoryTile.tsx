import { BagItem } from "../../models/BagItem";
import React, { useState, useRef } from "react";
import { createPortal } from "react-dom";
import inventory from "./inventory.module.css";

export interface InventoryTileProps {
  bagItem: BagItem;
}

export const InventoryTile: React.FC<InventoryTileProps> = ({ bagItem }) => {
  const [displayDetails, setDisplayDetails] = useState(false);
  const tileRef = useRef<HTMLDivElement>(null);

  const style = {
    borderColor: `var(--outline-${bagItem.rarity})`,
  } as React.CSSProperties;

  const handleMouseEnter = () => setDisplayDetails(true);
  const handleMouseLeave = () => setDisplayDetails(false);
  const handleTapToggle = () => setDisplayDetails(!displayDetails);

  return (
    <div
      ref={tileRef}
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
      {displayDetails && tileRef.current && (
        <ToolTip bagItem={bagItem} rect={tileRef.current.getBoundingClientRect()} />
      )}
    </div>
  );
};

interface ToolTipProps {
  bagItem: BagItem;
  rect: DOMRect;
}

const ToolTip: React.FC<ToolTipProps> = ({ bagItem, rect }) => {
  return createPortal(
    <div
      className={inventory.tooltip}
      style={{ top: rect.bottom, left: rect.left }}
    >
      <p className={inventory.name}>{bagItem.name}</p>
      <p className={inventory.description}>{bagItem.description}</p>
    </div>,
    document.body
  );
};
