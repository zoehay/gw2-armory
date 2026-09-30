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
    borderColor: `var(--outline-${bagItem.rarity}, var(--border))`,
  } as React.CSSProperties;

  const handleMouseEnter = () => setActiveId(myId);
  const handleMouseLeave = () => setActiveId(null);
  const handleTapToggle = () => setActiveId(displayDetails ? null : myId);

  return (
    <div
      ref={tileRef}
      className={inventory.tile}
      data-inventory-tile
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
  const tooltipMaxWidth = Math.min(256, viewportWidth * 0.8);
  const left = Math.max(
    4,
    Math.min(rect.left, viewportWidth - tooltipMaxWidth - 4),
  );
  // Open upward when the tile is in the lower part of the viewport
  const gap = 4;
  const openUp = rect.bottom > window.innerHeight * 0.6;
  const position = openUp
    ? { bottom: window.innerHeight - rect.top + gap, left }
    : { top: rect.bottom + gap, left };
  const nameStyle = {
    color: `var(--outline-${bagItem.rarity}, var(--text))`,
  } as React.CSSProperties;

  // Stats
  const details = bagItem.details;
  const defense =
    typeof details?.defense === "number" && details?.defense != 0
      ? details.defense
      : null;
  const minPower =
    typeof details?.min_power === "number" ? details.min_power : null;
  const maxPower =
    typeof details?.max_power === "number" ? details.max_power : null;

  // Attributes
  const infix = details?.infix_upgrade;
  const infixAttributes =
    infix != null &&
    typeof infix === "object" &&
    Array.isArray((infix as Record<string, unknown>).attributes)
      ? ((infix as Record<string, unknown>).attributes as {
          attribute: string;
          modifier: number;
        }[])
      : null;

  const topLevelAttributes =
    bagItem.stats != null &&
    typeof bagItem.stats.attributes === "object" &&
    bagItem.stats.attributes != null
      ? Object.entries(bagItem.stats.attributes as Record<string, number>).map(
          ([attribute, modifier]) => ({ attribute, modifier }),
        )
      : null;

  const attributes = infixAttributes ?? topLevelAttributes ?? [];
  const hasStats = defense !== null || (minPower !== null && maxPower !== null);

  return createPortal(
    <div className={inventory.tooltip} style={position}>
      <div className={inventory.name} style={nameStyle}>
        {bagItem.name}
      </div>
      {hasStats && (
        <div className={inventory.stats}>
          <ul>
            {defense !== null && <li>Defense {defense}</li>}
            {minPower !== null && maxPower !== null && (
              <li>
                Weapon Strength {minPower} - {maxPower}{" "}
              </li>
            )}
          </ul>
        </div>
      )}
      {attributes.length > 0 && (
        <div className={inventory.attributes}>
          <ul>
            {attributes.map(({ attribute, modifier }) => (
              <li key={attribute}>
                +{modifier} {AttributeLabels[attribute] ?? attribute}
              </li>
            ))}
          </ul>
        </div>
      )}
      {bagItem.upgradeDetails?.map((upgrade) => {
        const bonuses = Array.isArray(upgrade.details?.bonuses)
          ? (upgrade.details.bonuses as string[])
          : [];
        return (
          <div key={upgrade.id}>
            <div className={inventory.upgradeName}>{upgrade.name}</div>
            <ul>
              {bonuses.map((bonus, i) => (
                <li key={i}>{bonus}</li>
              ))}
            </ul>
          </div>
        );
      })}
      {bagItem.description && (
        <div className={inventory.description}>
          {parseDescription(bagItem.description)}
        </div>
      )}
      {bagItem.type && <div className={inventory.type}>{bagItem.type}</div>}
    </div>,
    document.body,
  );
};

function parseDescription(description: string): React.ReactNode {
  const parts = description
    .replace(/<br\s*\/?>/gi, "\n")
    .split(/(<c=@\w+>.*?<\/c>)/g);
  return parts.map((part, i) => {
    const match = part.match(/^<c=@(\w+)>(.*?)<\/c>$/s);
    if (match) {
      const [, tag, text] = match;
      return (
        <span key={i} className={inventory[tag]}>
          {text}
        </span>
      );
    }
    return part;
  });
}
