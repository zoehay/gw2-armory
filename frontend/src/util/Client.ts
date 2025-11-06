import { Account, APIAccount, APIAccountToAccount } from "../models/Account";
import {
  AccountInventory,
  APIAccountInventory,
  APIAccountInventoryToAccountInventory,
} from "../models/AccountInventory";
import { BagItem, APIBagItem, APIBagItemToBagItem } from "../models/BagItem";

export interface ClientInterface {
  getBagItems(): BagItem[];
}

export class Client {
  baseURL: string;

  constructor(mode: string) {
    this.baseURL = "/api";
    if (mode === "dev" || "test") {
      this.baseURL = "http://localhost:8000";
    }
  }

  async clientGet(endpoint: string): Promise<any> {
    try {
      let response = await fetch(endpoint, {
        credentials: "include",
      });
      if (response.ok) {
        let responseJSON = await response.json();
        return responseJSON;
      }
    } catch (error) {
      console.error(error);
    }
  }

  async clientPost(endpoint: string, body: any): Promise<any> {
    try {
      let response = await fetch(endpoint, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });

      if (response.ok) {
        let responseJSON = await response.json();
        return responseJSON;
      }
    } catch (error) {
      console.error(error);
    }
  }

  async clientDelete(endpoint: string, body: any): Promise<any> {
    try {
      let response = await fetch(endpoint, {
        credentials: "include",
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });

      if (response.ok) {
        let responseJSON = await response.json();
        return responseJSON;
      }
    } catch (error) {
      console.error(error);
    }
  }

  async getBagItems(): Promise<BagItem[]> {
    let endpoint: string = `${this.baseURL}/account/inventory`;

    let response: unknown = await this.clientGet(endpoint);
    let apiBagItems = response as APIBagItem[];
    let bagItems: BagItem[] = apiBagItems.map(APIBagItemToBagItem);
    if (bagItems) {
      return bagItems;
    } else {
      return [];
    }
  }

  async getAccountInventory(): Promise<AccountInventory> {
    let endpoint: string = `${this.baseURL}/account/accountinventory`;

    let response: unknown = await this.clientGet(endpoint);
    let apiAccountInventory = response as APIAccountInventory;
    let accountInventory: Account =
      APIAccountInventoryToAccountInventory(apiAccountInventory);
    return accountInventory;
  }

  async postAPIKey(key: string): Promise<Account | null> {
    let body = JSON.stringify({
      APIKey: key,
    });
    let endpoint: string = `${this.baseURL}/apikeys`;

    let response: unknown = await this.clientPost(endpoint, body);
    let apiAccount = response as APIAccount;
    let account: Account = APIAccountToAccount(apiAccount);
    if (account) {
      return account;
    } else {
      return null;
    }
  }

  async getAccount(): Promise<Account | null> {
    let endpoint: string = `${this.baseURL}/account/info`;

    let response: unknown = await this.clientGet(endpoint);
    let apiAccount = response as APIAccount;
    let account: Account = APIAccountToAccount(apiAccount);
    if (account) {
      return account;
    } else {
      return null;
    }
  }

  async deleteAPIKey(key: string): Promise<string | null> {
    let body = JSON.stringify({
      APIKey: key,
    });
    let endpoint: string = `${this.baseURL}/account/delete`;

    let response: unknown = await this.clientDelete(endpoint, body);
    let deletedKey = response as string;
    if (deletedKey) {
      return deletedKey;
    } else {
      return null;
    }
  }

  async postInventorySearch(
    searchTerm: string
  ): Promise<AccountInventory | null> {
    let body = JSON.stringify({
      SearchTerm: searchTerm,
    });
    let endpoint: string = `${this.baseURL}/account/searchinventory`;

    let response: unknown = await this.clientPost(endpoint, body);
    let apiAccountInventory = response as APIAccountInventory;
    let accountInventory: Account =
      APIAccountInventoryToAccountInventory(apiAccountInventory);
    if (accountInventory) {
      return accountInventory;
    } else {
      return null;
    }
  }
}
