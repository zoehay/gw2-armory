export interface Account {
  accountID: string;
  accountName?: string;
  gw2AccountName?: string;
  apiKey?: string;
  password?: string;
  sessionID?: string;
  session?: Session;
}

export interface APIAccount {
  id: string;
  account_name?: string;
  name?: string;
  gw2_name?: string;
  api_key?: string;
  password?: string;
  session_id?: string;
}

export function APIAccountToAccount(apiAccount: APIAccount): Account {
  return {
    accountID: apiAccount.id,
    accountName: apiAccount.account_name,
    gw2AccountName: apiAccount.gw2_name,
    apiKey: apiAccount.api_key,
    password: apiAccount.password,
    sessionID: apiAccount.session_id,
  };
}

interface Session {
  sessionID: string;
  expires: Date;
}
