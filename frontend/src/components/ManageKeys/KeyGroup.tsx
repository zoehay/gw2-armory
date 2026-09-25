import React from "react";
import { Account } from "../../models/Account";
import { KeyTile } from "./KeyTile";
import managekeys from "./managekeys.module.css";

interface KeyGroupProps {
  accounts?: Account[];
  handleUpdate: React.Dispatch<React.SetStateAction<Account | null>>;
}

export const KeyGroup: React.FC<KeyGroupProps> = ({
  accounts,
  handleUpdate,
}) => {
  let keyTiles;
  if (accounts) {
    keyTiles = accounts.map((account, index) => (
      <KeyTile account={account} handleUpdate={handleUpdate} key={index} />
    ));
  }

  return <div className={managekeys.group}>{keyTiles}</div>;
};
