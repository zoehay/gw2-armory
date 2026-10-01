import { beforeEach, describe, expect, it, vi } from "vitest";
import { Client } from "./Client";

const fetchMock = vi.fn<typeof fetch>();

function okResponse(body: unknown) {
  return Promise.resolve(new Response(JSON.stringify(body), { status: 200 }));
}

function errorResponse() {
  return Promise.resolve(new Response("nope", { status: 401 }));
}

function lastRequest() {
  const [url, init] = fetchMock.mock.lastCall!;
  return { url, init: init! };
}

let client: Client;

beforeEach(() => {
  fetchMock.mockReset();
  vi.stubGlobal("fetch", fetchMock);
  vi.stubEnv("VITE_APP_API_URL", "http://api");
  client = new Client();
});

describe("getAccount", () => {
  it("GETs account info and converts it", async () => {
    fetchMock.mockReturnValue(
      okResponse({ id: "acct", gw2_name: "Tester.1234", api_key: "KEY" }),
    );

    const account = await client.getAccount();

    expect(lastRequest().url).toBe("http://api/account/info");
    expect(lastRequest().init).toEqual({ credentials: "include" });
    expect(account).toMatchObject({
      accountID: "acct",
      gw2AccountName: "Tester.1234",
      apiKey: "KEY",
    });
  });

  it("returns null on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.getAccount()).toBeNull();
  });

  it("returns null when fetch rejects", async () => {
    const consoleError = vi
      .spyOn(console, "error")
      .mockImplementation(() => {});
    fetchMock.mockRejectedValue(new TypeError("network down"));

    expect(await client.getAccount()).toBeNull();
    expect(consoleError).toHaveBeenCalled();
  });
});

describe("getAccountInventory", () => {
  it("GETs the account inventory and converts it", async () => {
    fetchMock.mockReturnValue(
      okResponse({
        id: "acct",
        bank_inventory: [
          { character_name: "", source: "bank", id: 1, count: 3 },
        ],
      }),
    );

    const inventory = await client.getAccountInventory();

    expect(lastRequest().url).toBe("http://api/account/accountinventory");
    expect(inventory?.accountID).toBe("acct");
    expect(inventory?.bankInventory?.[0]).toMatchObject({ id: 1, count: 3 });
  });

  it("returns null on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.getAccountInventory()).toBeNull();
  });
});

describe("getBagItems", () => {
  it("GETs bag items and converts them", async () => {
    fetchMock.mockReturnValue(
      okResponse([
        { character_name: "Char", source: "character", id: 2, count: 1 },
      ]),
    );

    const items = await client.getBagItems();

    expect(lastRequest().url).toBe("http://api/account/inventory");
    expect(items).toEqual([
      expect.objectContaining({ characterName: "Char", id: 2 }),
    ]);
  });

  it("returns an empty list on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.getBagItems()).toEqual([]);
  });
});

describe("postAPIKey", () => {
  it("POSTs the key as JSON and converts the account", async () => {
    fetchMock.mockReturnValue(okResponse({ id: "acct", api_key: "KEY" }));

    const account = await client.postAPIKey("KEY");

    const { url, init } = lastRequest();
    expect(url).toBe("http://api/apikeys");
    expect(init).toMatchObject({
      credentials: "include",
      method: "POST",
      headers: { "Content-Type": "application/json" },
    });
    expect(JSON.parse(init.body as string)).toEqual({ APIKey: "KEY" });
    expect(account).toMatchObject({ accountID: "acct", apiKey: "KEY" });
  });

  it("returns null on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.postAPIKey("bad")).toBeNull();
  });
});

describe("deleteAPIKey", () => {
  it("DELETEs the key and returns the response", async () => {
    fetchMock.mockReturnValue(okResponse("KEY"));

    const deleted = await client.deleteAPIKey("KEY");

    const { url, init } = lastRequest();
    expect(url).toBe("http://api/account/delete");
    expect(init).toMatchObject({ credentials: "include", method: "DELETE" });
    expect(JSON.parse(init.body as string)).toEqual({ APIKey: "KEY" });
    expect(deleted).toBe("KEY");
  });

  it("returns null on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.deleteAPIKey("KEY")).toBeNull();
  });
});

describe("postInventorySearch", () => {
  it("POSTs the search term and converts the result", async () => {
    fetchMock.mockReturnValue(okResponse({ id: "acct" }));

    const inventory = await client.postInventorySearch("sword");

    const { url, init } = lastRequest();
    expect(url).toBe("http://api/account/searchinventory");
    expect(init).toMatchObject({ method: "POST" });
    expect(JSON.parse(init.body as string)).toEqual({ SearchTerm: "sword" });
    expect(inventory?.accountID).toBe("acct");
  });

  it("returns null on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.postInventorySearch("sword")).toBeNull();
  });
});
