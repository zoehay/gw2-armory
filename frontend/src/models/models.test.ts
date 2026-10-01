import { describe, expect, it } from "vitest";
import { APIAccountToAccount } from "./Account";
import { APIAccountInventoryToAccountInventory } from "./AccountInventory";
import { APIBagItem, APIBagItemToBagItem } from "./BagItem";
import { APICharacterToCharacter } from "./Character";

const upgrade = {
  id: 7,
  name: "Superior Rune",
  icon: "rune.png",
  rarity: "Exotic",
};

const apiBagItem: APIBagItem = {
  character_name: "Char",
  source: "character",
  id: 100,
  count: 2,
  charges: 3,
  infusions: [1],
  upgrades: [7],
  infusion_details: [upgrade],
  upgrade_details: [upgrade],
  skin: 5,
  stats: { id: 1 },
  dyes: [9],
  binding: "Account",
  bound_to: "Someone",
  slot: "Helm",
  location: "Equipped",
  name: "Helm",
  icon: "helm.png",
  description: "A helm",
  type: "Armor",
  rarity: "Exotic",
  vendor_value: 42,
  details: { defense: 10 },
};

describe("APIBagItemToBagItem", () => {
  it("maps every field to camelCase", () => {
    expect(APIBagItemToBagItem(apiBagItem)).toEqual({
      characterName: "Char",
      source: "character",
      id: 100,
      count: 2,
      charges: 3,
      infusions: [1],
      upgrades: [7],
      infusionDetails: [upgrade],
      upgradeDetails: [upgrade],
      skin: 5,
      stats: { id: 1 },
      dyes: [9],
      binding: "Account",
      boundTo: "Someone",
      slot: "Helm",
      location: "Equipped",
      name: "Helm",
      icon: "helm.png",
      description: "A helm",
      type: "Armor",
      rarity: "Exotic",
      vendorValue: 42,
      details: { defense: 10 },
    });
  });
});

describe("APICharacterToCharacter", () => {
  it("maps equipment and inventory items", () => {
    const character = APICharacterToCharacter({
      name: "Char",
      equipment: [apiBagItem],
      inventory: [apiBagItem, apiBagItem],
    });
    expect(character.name).toBe("Char");
    expect(character.equipment).toHaveLength(1);
    expect(character.inventory).toHaveLength(2);
    expect(character.equipment?.[0].characterName).toBe("Char");
  });

  it("leaves missing arrays undefined", () => {
    expect(APICharacterToCharacter({ name: "Empty" })).toEqual({
      name: "Empty",
      equipment: undefined,
      inventory: undefined,
    });
  });
});

describe("APIAccountInventoryToAccountInventory", () => {
  it("maps each inventory section", () => {
    const inventory = APIAccountInventoryToAccountInventory({
      id: "acct",
      shared_inventory: [apiBagItem],
      bank_inventory: [apiBagItem, apiBagItem],
      materials_inventory: [apiBagItem],
      characters: [{ name: "Char", inventory: [apiBagItem] }],
    });
    expect(inventory.accountID).toBe("acct");
    expect(inventory.sharedInventory).toHaveLength(1);
    expect(inventory.bankInventory).toHaveLength(2);
    expect(inventory.materialsInventory?.[0].vendorValue).toBe(42);
    expect(inventory.characters?.[0].inventory?.[0].boundTo).toBe("Someone");
  });

  it("leaves missing sections undefined", () => {
    expect(APIAccountInventoryToAccountInventory({ id: "acct" })).toEqual({
      accountID: "acct",
      sharedInventory: undefined,
      bankInventory: undefined,
      materialsInventory: undefined,
      characters: undefined,
    });
  });
});

describe("APIAccountToAccount", () => {
  it("maps fields to camelCase", () => {
    expect(
      APIAccountToAccount({
        id: "acct",
        account_name: "name",
        gw2_name: "Tester.1234",
        gw2_token_name: "my key",
        api_key: "KEY",
        password: "pw",
        session_id: "sess",
      }),
    ).toEqual({
      accountID: "acct",
      accountName: "name",
      gw2AccountName: "Tester.1234",
      gw2TokenName: "my key",
      apiKey: "KEY",
      password: "pw",
      sessionID: "sess",
    });
  });
});
