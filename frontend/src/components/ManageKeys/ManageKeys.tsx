import React, { useContext, useState, useEffect } from "react";
import { ClientContext } from "../../util/ClientContext";
import content from "../content.module.css";
import { Account } from "../../models/Account";
import { KeyGroup } from "./KeyGroup";
import managekeys from "./managekeys.module.css";

export const ManageKeys = () => {
  // if user show UserKeys else only one account AccountKey
  return (
    <>
      <AccountKey />
    </>
  );
};

const AccountKey = () => {
  const client = useContext(ClientContext);

  const [account, setAccount] = useState<Account | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const fetchAccount = await client.getAccount();
        setAccount(fetchAccount);
      } finally {
        setLoading(false);
      }
    };
    void fetchData();
  }, [client]);

  let body;
  if (loading) {
    body = <p className={content.status}>Loading…</p>;
  } else if (account) {
    body = <KeyGroup accounts={[account]} handleUpdate={setAccount}></KeyGroup>;
  } else {
    body = <KeyInput handleUpdate={setAccount}></KeyInput>;
  }

  return (
    <div className={content.page}>
      <div className={managekeys.container}>
        <header className={content.header}>
          <h1 className={content.title}>API Keys</h1>
          <p className={content.subtitle}>
            Your Guild Wars 2 API key lets armory read your account and
            inventory.
          </p>
        </header>
        {body}
      </div>
    </div>
  );
};

interface KeyInputProps {
  handleUpdate: React.Dispatch<React.SetStateAction<Account | null>>;
}

const KeyInput: React.FC<KeyInputProps> = ({ handleUpdate }) => {
  const [formState, setFormState] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const client = useContext(ClientContext);

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const input = e.currentTarget.value;
    setFormState(input);
    setError(null);
  };

  const handleSubmit = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    setSubmitting(true);
    let account = null;
    try {
      account = await client.postAPIKey(formState.trim());
    } catch (err) {
      console.error(err);
    }
    setSubmitting(false);
    if (!account) {
      setError("Couldn't add that key. Check that it's valid and try again.");
    } else {
      handleUpdate(account);
    }
  };

  return (
    <form
      className={`${content.card} ${managekeys.form}`}
      onSubmit={(e) => void handleSubmit(e)}
    >
      <label htmlFor="apikey-input" className={managekeys.label}>
        Add an API key
      </label>
      <p className={managekeys.hint}>
        Create a key at{" "}
        <a
          href="https://account.arena.net/applications"
          target="_blank"
          rel="noreferrer"
        >
          account.arena.net/applications
        </a>{" "}
        with the account, inventories, and characters permissions.
      </p>
      <div className={managekeys.inputRow}>
        <input
          type="text"
          name="apikey-input"
          id="apikey-input"
          className={`${content.input} ${content.mono}`}
          placeholder="Paste your API key"
          autoComplete="off"
          spellCheck={false}
          value={formState}
          onChange={handleChange}
          aria-invalid={error !== null}
          aria-describedby={error ? "apikey-error" : undefined}
        />
        <button
          type="submit"
          className={`${content.button} ${content.primary}`}
          disabled={submitting || formState.trim() === ""}
        >
          {submitting ? "Adding…" : "Add key"}
        </button>
      </div>
      {error && (
        <p id="apikey-error" className={content.error} role="alert">
          {error}
        </p>
      )}
    </form>
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
