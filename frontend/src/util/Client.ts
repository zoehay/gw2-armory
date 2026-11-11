import { Account, APIAccount, APIAccountToAccount } from "../models/Account";
import {
  AccountInventory,
  APIAccountInventory,
  APIAccountInventoryToAccountInventory,
} from "../models/AccountInventory";
import { BagItem, APIBagItem, APIBagItemToBagItem } from "../models/BagItem";

export type APIKeyError =
  | "invalid_key"
  | "missing_permissions"
  | "rate_limited"
  | "unavailable"
  | "unknown";

export type PostAPIKeyResult =
  { account: Account } | { error: APIKeyError; missingPermissions?: string[] };

interface ClientResponse {
  status: number;
  body: unknown;
}

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

  // like clientPost, but keeps the status and body of error responses so callers can
  // tell failures apart. Returns null when the request itself fails
  async clientPostResult(
    endpoint: string,
    body: string,
  ): Promise<ClientResponse | null> {
    try {
      const response = await fetch(endpoint, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });

      let responseBody: unknown = null;
      try {
        responseBody = await response.json();
      } catch {
        // non-JSON bodies, like nginx's rate limit page
      }
      return { status: response.status, body: responseBody };
    } catch (error) {
      console.error(error);
      return null;
    }
  }

  async clientDelete(endpoint: string): Promise<unknown> {
    try {
      const response = await fetch(endpoint, {
        credentials: "include",
        method: "DELETE",
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

  async postAPIKey(key: string): Promise<PostAPIKeyResult> {
    const body = JSON.stringify({
      APIKey: key,
    });
    const endpoint: string = `${this.baseURL}/apikeys`;

    const response = await this.clientPostResult(endpoint, body);
    if (!response) {
      return { error: "unavailable" };
    }

    switch (response.status) {
      case 200:
        return { account: APIAccountToAccount(response.body as APIAccount) };
      case 400: {
        // the client always sends valid JSON, so a 400 is about the key itself
        const missing = (response.body as { missing_permissions?: unknown })
          ?.missing_permissions;
        if (Array.isArray(missing)) {
          return {
            error: "missing_permissions",
            missingPermissions: missing as string[],
          };
        }
        return { error: "invalid_key" };
      }
      case 429:
        return { error: "rate_limited" };
      case 502:
        return { error: "unavailable" };
      default:
        return { error: "unknown" };
    }
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

  // the session cookie identifies the account, so no key is sent
  async deleteAPIKey(): Promise<boolean> {
    const endpoint: string = `${this.baseURL}/account/delete`;

    const response: unknown = await this.clientDelete(endpoint);
    return Boolean(response);
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
