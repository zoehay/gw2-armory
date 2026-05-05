import { BagItem } from "../../models/BagItem";
import React, { useRef, useId, useContext } from "react";
import { createPortal } from "react-dom";
import inventory from "./inventory.module.css";
import { TooltipContext } from "./TooltipContext";

export interface InventoryTileProps {
  bagItem: BagItem;
}

export const InventoryTile: React.FC<InventoryTileProps> = ({ bagItem }) => {
  const myId = useId();
  const { activeId, setActiveId } = useContext(TooltipContext);
  const displayDetails = activeId === myId;
  const tileRef = useRef<HTMLDivElement>(null);

  const style = {
    borderColor: `var(--outline-${bagItem.rarity})`,
  } as React.CSSProperties;

  const handleMouseEnter = () => setActiveId(myId);
  const handleMouseLeave = () => setActiveId(null);
  const handleTapToggle = () => setActiveId(displayDetails ? null : myId);

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
        <ToolTip
          bagItem={bagItem}
          rect={tileRef.current.getBoundingClientRect()}
        />
      )}
    </div>
  );
};

const AttributeLabels: Record<string, string> = {
  AgonyResistance: "Agony Resistance",
  BoonDuration: "Concentration",
  ConditionDamage: "Condition Damage",
  ConditionDuration: "Expertise",
  CritDamage: "Ferocity",
  Healing: "Healing Power",
  Power: "Power",
  Precision: "Precision",
  Toughness: "Toughness",
  Vitality: "Vitality",
};

interface ToolTipProps {
  bagItem: BagItem;
  rect: DOMRect;
}

const ToolTip: React.FC<ToolTipProps> = ({ bagItem, rect }) => {
  const viewportWidth = window.innerWidth;
  const TOOLTIP_MAX_W = Math.min(256, viewportWidth * 0.8); // mirrors CSS min(16rem, 80vw)
  const overflowsRight = rect.left + TOOLTIP_MAX_W > viewportWidth;

  const position = overflowsRight
    ? { top: rect.bottom, right: viewportWidth - rect.right }
    : { top: rect.bottom, left: rect.left };

  const details = bagItem.details;
  const defense =
    typeof details?.defense === "number" && details?.defense != 0
      ? details.defense
      : null;
  const minPower =
    typeof details?.min_power === "number" ? details.min_power : null;
  const maxPower =
    typeof details?.max_power === "number" ? details.max_power : null;

  const infix = details?.infix_upgrade;
  const attributes =
    infix != null &&
    typeof infix === "object" &&
    Array.isArray((infix as Record<string, unknown>).attributes)
      ? ((infix as Record<string, unknown>).attributes as {
          attribute: string;
          modifier: number;
        }[])
      : [];

  return createPortal(
    <div className={inventory.tooltip} style={position}>
      <p className={inventory.name}>{bagItem.name}</p>
      <p className={inventory.description}>{bagItem.description}</p>
      {defense !== null && <p>Defense: {defense}</p>}
      {minPower !== null && maxPower !== null && (
        <p>
          Weapon Strength: {minPower} - {maxPower}
        </p>
      )}
      {attributes.map(({ attribute, modifier }) => (
        <p key={attribute}>
          +{modifier} {AttributeLabels[attribute] ?? attribute}
        </p>
      ))}
    </div>,
    document.body,
  );
};
