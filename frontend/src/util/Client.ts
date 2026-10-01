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

  constructor() {
    this.baseURL = import.meta.env.VITE_APP_API_URL;
  }

  async clientGet(endpoint: string): Promise<unknown> {
    try {
      const response = await fetch(endpoint, {
        credentials: "include",
      });
      if (response.ok) {
        const responseJSON: unknown = await response.json();
        return responseJSON;
      }
    } catch (error) {
      console.error(error);
    }
  }

  async clientPost(endpoint: string, body: string): Promise<unknown> {
    try {
      const response = await fetch(endpoint, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });

      if (response.ok) {
        const responseJSON: unknown = await response.json();
        return responseJSON;
      }
    } catch (error) {
      console.error(error);
    }
  }

  async clientDelete(endpoint: string, body: string): Promise<unknown> {
    try {
      const response = await fetch(endpoint, {
        credentials: "include",
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });

      if (response.ok) {
        const responseJSON: unknown = await response.json();
        return responseJSON;
      }
    } catch (error) {
      console.error(error);
    }
  }

  async getBagItems(): Promise<BagItem[]> {
    const endpoint: string = `${this.baseURL}/account/inventory`;

    const response: unknown = await this.clientGet(endpoint);
    if (!response) {
      return [];
    }
    const apiBagItems = response as APIBagItem[];
    return apiBagItems.map(APIBagItemToBagItem);
  }

  async getAccountInventory(): Promise<AccountInventory | null> {
    const endpoint: string = `${this.baseURL}/account/accountinventory`;

    const response: unknown = await this.clientGet(endpoint);
    if (!response) {
      return null;
    }
    const apiAccountInventory = response as APIAccountInventory;
    return APIAccountInventoryToAccountInventory(apiAccountInventory);
  }

  async postAPIKey(key: string): Promise<Account | null> {
    const body = JSON.stringify({
      APIKey: key,
    });
    const endpoint: string = `${this.baseURL}/apikeys`;

    const response: unknown = await this.clientPost(endpoint, body);
    if (!response) {
      return null;
    }
    const apiAccount = response as APIAccount;
    return APIAccountToAccount(apiAccount);
  }

  async getAccount(): Promise<Account | null> {
    const endpoint: string = `${this.baseURL}/account/info`;

    const response: unknown = await this.clientGet(endpoint);
    if (!response) {
      return null;
    }
    const apiAccount = response as APIAccount;
    return APIAccountToAccount(apiAccount);
  }

  async deleteAPIKey(key: string): Promise<string | null> {
    const body = JSON.stringify({
      APIKey: key,
    });
    const endpoint: string = `${this.baseURL}/account/delete`;

    const response: unknown = await this.clientDelete(endpoint, body);
    const deletedKey = response as string;
    if (deletedKey) {
      return deletedKey;
    } else {
      return null;
    }
  }

  async postInventorySearch(
    searchTerm: string,
  ): Promise<AccountInventory | null> {
    const body = JSON.stringify({
      SearchTerm: searchTerm,
    });
    const endpoint: string = `${this.baseURL}/account/searchinventory`;

    const response: unknown = await this.clientPost(endpoint, body);
    if (!response) {
      return null;
    }
    const apiAccountInventory = response as APIAccountInventory;
    return APIAccountInventoryToAccountInventory(apiAccountInventory);
  }
}
