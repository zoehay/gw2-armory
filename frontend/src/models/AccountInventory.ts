import { APIBagItem, APIBagItemToBagItem, BagItem } from "./BagItem";
import { APICharacter, APICharacterToCharacter, Character } from "./Character";

export interface AccountInventory {
  accountID: string;
  sharedInventory?: BagItem[];
  bankInventory?: BagItem[];
  characters?: Character[];
}

export interface APIAccountInventory {
  id: string;
  shared_inventory?: APIBagItem[];
  bank_inventory?: APIBagItem[];
  characters?: APICharacter[];
}

export function APIAccountInventoryToAccountInventory(
  apiInventory: APIAccountInventory
): AccountInventory {
  return {
    accountID: apiInventory.id,
    sharedInventory: apiInventory.shared_inventory?.map(APIBagItemToBagItem),
    bankInventory: apiInventory.bank_inventory?.map(APIBagItemToBagItem),
    characters: apiInventory.characters?.map(APICharacterToCharacter),
  };
}
