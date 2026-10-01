import { useContext, useEffect, useState } from "react";
import { ClientContext } from "../../util/ClientContext";
import content from "../content.module.css";
import { NoKeyPage } from "../ErrorPage/NoKeyPage";
import { AccountInventory } from "../../models/AccountInventory";
import { InventoryContents } from "./InventoryContents";
import inventory from "./inventory.module.css";

export const Inventory = () => {
  const client = useContext(ClientContext);

  const [accountInventory, setAccountInventory] =
    useState<AccountInventory | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const inventory = await client.getAccountInventory();
        setAccountInventory(inventory);
      } finally {
        setLoading(false);
      }
    };
    void fetchData();
  }, [client]);

  let body;
  if (loading) {
    body = <p className={content.status}>Loading…</p>;
  } else if (accountInventory) {
    body = (
      <InventoryContents
        accountInventory={accountInventory}
      ></InventoryContents>
    );
  } else {
    body = <NoKeyPage />;
  }

  return (
    <div className={content.page}>
      <header className={inventory.pageHeader}>
        <div>
          <h1 className={content.title}>Inventory</h1>
          <p className={content.subtitle}>
            Items across your bank, materials, shared slots, and characters.
          </p>
        </div>
        {accountInventory && (
          <SearchInput handleUpdate={setAccountInventory}></SearchInput>
        )}
      </header>
      {body}
    </div>
  );
};

interface SearchInputProps {
  handleUpdate: React.Dispatch<React.SetStateAction<AccountInventory | null>>;
}

const SearchInput: React.FC<SearchInputProps> = ({ handleUpdate }) => {
  const [formState, setFormState] = useState("");
  const [searching, setSearching] = useState(false);
  const client = useContext(ClientContext);

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const input = e.currentTarget.value;
    setFormState(input);
  };

  const handleSubmit = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    setSearching(true);
    try {
      const accountInventory = await client.postInventorySearch(formState);
      if (!accountInventory) {
        console.log("Could not post search");
      } else {
        handleUpdate(accountInventory);
      }
    } finally {
      setSearching(false);
    }
  };

  return (
    <form
      className={inventory.search}
      role="search"
      onSubmit={(e) => void handleSubmit(e)}
    >
      <label htmlFor="search-input" className={content.srOnly}>
        Search items
      </label>
      <input
        type="search"
        name="search-input"
        id="search-input"
        className={content.input}
        placeholder="Search items"
        value={formState}
        onChange={handleChange}
      />
      <button
        type="submit"
        className={`${content.button} ${content.primary}`}
        disabled={searching}
      >
        {searching ? "Searching…" : "Search"}
      </button>
    </form>
  );
};
