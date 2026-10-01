import React from "react";
import { render } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { Client } from "../util/Client";
import { ClientContext } from "../util/ClientContext";
import { Account } from "../models/Account";
import { BagItem } from "../models/BagItem";

export function renderWithClient(
  ui: React.ReactElement,
  client: Partial<Client> = {},
  { route = "/" }: { route?: string } = {},
) {
  return render(
    <ClientContext.Provider value={client as Client}>
      <MemoryRouter initialEntries={[route]}>{ui}</MemoryRouter>
    </ClientContext.Provider>,
  );
}

export function makeBagItem(overrides: Partial<BagItem> = {}): BagItem {
  return {
    characterName: "",
    source: "bank",
    id: 1,
    count: 1,
    name: "Test Item",
    icon: "https://example.com/icon.png",
    rarity: "Rare",
    ...overrides,
  };
}

export function makeAccount(overrides: Partial<Account> = {}): Account {
  return {
    accountID: "account-id",
    gw2AccountName: "Tester.1234",
    gw2TokenName: "my key",
    apiKey: "API-KEY",
    ...overrides,
  };
}
