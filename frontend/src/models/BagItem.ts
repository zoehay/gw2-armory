export interface BagItem {
  characterName: string;
  name?: string;
  description?: string;
  id: number;
  icon: string;
  count: number;
  charges?: number;
  infusions?: number[];
  upgrades?: number[];
  skin?: number;
  stats?: { [key: string]: unknown };
  dyes?: number[];
  binding?: string;
  boundTo?: string;
  rarity?: string;
  details?: { [key: string]: unknown };
}

export interface APIBagItem {
  character_name: string;
  name?: string;
  description?: string;
  id: number;
  icon: string;
  count: number;
  charges?: number;
  infusions?: number[];
  upgrades?: number[];
  skin?: number;
  stats?: { [key: string]: unknown };
  dyes?: number[];
  binding?: string;
  bound_to?: string;
  rarity?: string;
  details?: { [key: string]: unknown };
}

export function APIBagItemToBagItem(apiBagItem: APIBagItem): BagItem {
  return {
    characterName: apiBagItem.character_name,
    name: apiBagItem.name,
    description: apiBagItem.description,
    id: apiBagItem.id,
    icon: apiBagItem.icon,
    count: apiBagItem.count,
    charges: apiBagItem.charges,
    infusions: apiBagItem.infusions,
    upgrades: apiBagItem.upgrades,
    skin: apiBagItem.skin,
    stats: apiBagItem.stats,
    dyes: apiBagItem.dyes,
    binding: apiBagItem.binding,
    boundTo: apiBagItem.bound_to,
    rarity: apiBagItem.rarity as string,
    details: apiBagItem.details,
  };
}
