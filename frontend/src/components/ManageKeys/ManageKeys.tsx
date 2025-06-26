import React, { useContext, useState, useEffect } from "react";
import { ClientContext } from "../../util/ClientContext";
import content from "../content.module.css";
import { Account } from "../../models/Account";
import { KeyGroup } from "./KeyGroup";

export const ManageKeys = () => {
  // if user show UserKeys else only one account AccountKey
  return (
    <>
      <AccountKey />
    </>
  );
};

const AccountKey = () => {
  let context = useContext(ClientContext);
  let client = context;

  let [account, setAccount] = useState<Account | null>(null);

  async function fetchData() {
    let fetchAccount = await client.getAccount();
    setAccount(fetchAccount);
  }

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div className={content.page}>
      {account ? (
        <>
          <p>Keys</p>
          <KeyGroup accounts={[account]} handleUpdate={setAccount}></KeyGroup>
        </>
      ) : (
        <>
          <p>No keys</p>
          <KeyInput handleUpdate={setAccount}></KeyInput>
        </>
      )}
    </div>
  );
};

interface KeyInputProps {
  handleUpdate: React.Dispatch<React.SetStateAction<Account | null>>;
}

const KeyInput: React.FC<KeyInputProps> = ({ handleUpdate }) => {
  const fieldName = "API Key";
  const [formState, setFormState] = useState("");
  let context = useContext(ClientContext);
  let client = context;

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const input = e.currentTarget.value;
    setFormState(input);
  };

  const handleSubmit = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    const account = await client.postAPIKey(formState);
    if (!account) {
      console.log("Could not post key");
    } else {
      handleUpdate(account);
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <label htmlFor="apikey-input">{`Add ${fieldName}`}</label>
        <div>
          <input
            type="apikey"
            name="apikey-input"
            id="input"
            value={formState}
            onChange={handleChange}
          />
        </div>
        <input type="submit" value="Submit" />
      </form>
    </div>
  );
};

// A User can have multiple Accounts / Keys
// const UserKeys = () => {

//   let [accounts, setAccounts] = useState<Account | null>(null);

//   async function fetchData() {
//     let fetchAccount = await client.getAccounts();
//     setAccounts(fetchAccount);
//   }

//   useEffect(() => {
//     fetchData();
//   }, []);

//   return (
//     <div className={content.main}>
//       <p>Add A Key</p>
//       <KeyInput></KeyInput>

//       {accounts[0] != null && accounts[0].apiKey != null ? (
//         <>
//           <p>Keys</p>
//           <KeyGroup accounts={accounts}></KeyGroup>
//         </>
//       ) : (
//         <p>No keys</p>
//       )}
//     </div>
//   );
// }
