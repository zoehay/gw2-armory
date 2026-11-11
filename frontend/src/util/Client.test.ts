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
      okResponse({ id: "acct", gw2_name: "Tester.1234" }),
    );

    const account = await client.getAccount();

    expect(lastRequest().url).toBe("http://api/account/info");
    expect(lastRequest().init).toEqual({ credentials: "include" });
    expect(account).toMatchObject({
      accountID: "acct",
      gw2AccountName: "Tester.1234",
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
    fetchMock.mockReturnValue(okResponse({ id: "acct" }));

    const account = await client.postAPIKey("KEY");

    const { url, init } = lastRequest();
    expect(url).toBe("http://api/apikeys");
    expect(init).toMatchObject({
      credentials: "include",
      method: "POST",
      headers: { "Content-Type": "application/json" },
    });
    expect(JSON.parse(init.body as string)).toEqual({ APIKey: "KEY" });
    expect(account).toMatchObject({ account: { accountID: "acct" } });
  });

  it.each([
    [
      "missing permissions",
      400,
      {
        error: "API key is missing required permissions",
        missing_permissions: ["wallet"],
      },
      { error: "missing_permissions", missingPermissions: ["wallet"] },
    ],
    [
      "an invalid key",
      400,
      { error: "invalid API key" },
      { error: "invalid_key" },
    ],
    [
      "rate limiting",
      429,
      "<html>429 Too Many Requests</html>",
      { error: "rate_limited" },
    ],
    [
      "a GW2 API failure",
      502,
      { error: "could not verify API key with the GW2 API" },
      { error: "unavailable" },
    ],
    [
      "an unexpected status",
      500,
      { error: "could not create or update account" },
      { error: "unknown" },
    ],
  ])("maps %s to an error result", async (_, status, body, expected) => {
    const text = typeof body === "string" ? body : JSON.stringify(body);
    fetchMock.mockReturnValue(Promise.resolve(new Response(text, { status })));

    expect(await client.postAPIKey("KEY")).toEqual(expected);
  });

  it("reports the GW2 API as unavailable when the request fails", async () => {
    vi.spyOn(console, "error").mockImplementation(() => {});
    fetchMock.mockRejectedValue(new TypeError("Failed to fetch"));

    expect(await client.postAPIKey("KEY")).toEqual({ error: "unavailable" });
  });
});

describe("deleteAPIKey", () => {
  it("DELETEs the session's account without sending the key", async () => {
    fetchMock.mockReturnValue(okResponse({ deleted: "acct" }));

    const deleted = await client.deleteAPIKey();

    const { url, init } = lastRequest();
    expect(url).toBe("http://api/account/delete");
    expect(init).toMatchObject({ credentials: "include", method: "DELETE" });
    expect(init.body).toBeUndefined();
    expect(deleted).toBe(true);
  });

  it("returns false on an error response", async () => {
    fetchMock.mockReturnValue(errorResponse());
    expect(await client.deleteAPIKey()).toBe(false);
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
