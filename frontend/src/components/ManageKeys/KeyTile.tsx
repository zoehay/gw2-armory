import React, { useContext, useState } from "react";
import { ClientContext } from "../../util/ClientContext";
import { Account } from "../../models/Account";
import content from "../content.module.css";
import managekeys from "./managekeys.module.css";

interface KeyTileProps {
  account: Account;
  handleUpdate: React.Dispatch<React.SetStateAction<Account | null>>;
}

export const KeyTile: React.FC<KeyTileProps> = ({ account, handleUpdate }) => {
  const client = useContext(ClientContext);
  const [deleting, setDeleting] = useState(false);

  const handleClick = () => {
    setDeleting(true);
    setTimeout(async () => {
      let deletedAccount;
      if (account.apiKey) {
        deletedAccount = await client.deleteAPIKey(account.apiKey);
      }
      if (deletedAccount) {
        handleUpdate(null);
      } else {
        setDeleting(false);
      }
    }, 2000);
  };

  return (
    <div className={`${content.card} ${managekeys.keytile}`}>
      <dl className={managekeys.fields}>
        <dt>Key name</dt>
        <dd>{account.gw2TokenName || "—"}</dd>
        <dt>Account</dt>
        <dd>{account.gw2AccountName || "—"}</dd>
        <dt>Account ID</dt>
        <dd className={content.mono} title={account.accountID}>
          {account.accountID}
        </dd>
      </dl>
      <button
        type="button"
        className={`${content.button} ${content.danger}`}
        onClick={handleClick}
        disabled={deleting}
      >
        {deleting ? "Removing…" : "Remove key"}
      </button>
    </div>
  );
};
