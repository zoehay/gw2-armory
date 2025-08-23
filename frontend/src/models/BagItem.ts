export interface BagItem {
  characterName: string;
  id: number;
  count: number;
  charges?: number;
  infusions?: number[];
  upgrades?: number[];
  skin?: number;
  stats?: { [key: string]: unknown };
  dyes?: number[];
  binding?: string;
  boundTo?: string;
  slot?: string;
  location?: string;

  name?: string;
  icon?: string;
  description?: string;
  type?: string;
  rarity?: string;
  vendorValue?: number;
  details?: { [key: string]: unknown };
}

export interface APIBagItem {
  character_name: string;
  id: number;
  count: number;
  charges?: number;
  infusions?: number[];
  upgrades?: number[];
  skin?: number;
  stats?: { [key: string]: unknown };
  dyes?: number[];
  binding?: string;
  bound_to?: string;
  slot?: string;
  location?: string;

  name?: string;
  icon?: string;
  description?: string;
  type?: string;
  rarity?: string;
  vendor_value?: number;
  details?: { [key: string]: unknown };
}

export function APIBagItemToBagItem(apiBagItem: APIBagItem): BagItem {
  return {
    characterName: apiBagItem.character_name,
    id: apiBagItem.id,
    count: apiBagItem.count,
    charges: apiBagItem.charges,
    infusions: apiBagItem.infusions,
    upgrades: apiBagItem.upgrades,
    skin: apiBagItem.skin,
    stats: apiBagItem.stats,
    dyes: apiBagItem.dyes,
    binding: apiBagItem.binding,
    boundTo: apiBagItem.bound_to,
    slot: apiBagItem.slot,
    location: apiBagItem.location,

    name: apiBagItem.name,
    icon: apiBagItem.icon,
    description: apiBagItem.description,
    type: apiBagItem.type,
    rarity: apiBagItem.rarity,
    vendorValue: apiBagItem.vendor_value,
    details: apiBagItem.details,
  };
}
