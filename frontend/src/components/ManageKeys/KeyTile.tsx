import React, { useContext } from "react";
import { ClientContext } from "../../util/ClientContext";
import { Account } from "../../models/Account";
import managekeys from "./managekeys.module.css";

interface KeyTileProps {
  account: Account;
  handleUpdate: React.Dispatch<React.SetStateAction<Account | null>>;
}

export const KeyTile: React.FC<KeyTileProps> = ({ account, handleUpdate }) => {
  let context = useContext(ClientContext);
  let client = context;

  const handleClick = async () => {
    setTimeout(async () => {
      let deletedAccount;
      if (account.apiKey) {
        deletedAccount = await client.deleteAPIKey(account.apiKey);
      }
      if (deletedAccount) {
        handleUpdate(null);
      }
    }, 2000);
  };

  return (
    <div className={managekeys.keytile}>
      <div className={managekeys.field}>{account.gw2TokenName}</div>
      <div className={managekeys.field}>{account.gw2AccountName}</div>
      <div className={managekeys.field}>{account.accountID}</div>
      <input type="button" onClick={handleClick} value="Delete Account"></input>
    </div>
  );
};
