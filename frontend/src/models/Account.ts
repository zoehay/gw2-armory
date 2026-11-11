export interface Account {
  accountID: string;
  accountName?: string;
  gw2AccountName?: string;
  gw2TokenName?: string;
  permissions: string[];
}

export interface APIAccount {
  id: string;
  account_name?: string;
  name?: string;
  gw2_name?: string;
  gw2_token_name?: string;
  // null for accounts stored before permissions were tracked
  permissions?: string[] | null;
}

export function APIAccountToAccount(apiAccount: APIAccount): Account {
  return {
    accountID: apiAccount.id,
    accountName: apiAccount.account_name,
    gw2AccountName: apiAccount.gw2_name,
    gw2TokenName: apiAccount.gw2_token_name,
    permissions: apiAccount.permissions ?? [],
  };
}
