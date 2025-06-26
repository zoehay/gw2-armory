import { useContext, useEffect, useState } from "react";
import { ClientContext } from "../../util/ClientContext";
import content from "../content.module.css";
import { NoKeyPage } from "../ErrorPage/NoKeyPage";
import { AccountInventory } from "../../models/AccountInventory";
import { InventoryContents } from "./InventoryContents";

export const Inventory = () => {
  let context = useContext(ClientContext);
  let client = context;

  let [accountInventory, setAccountInventory] =
    useState<AccountInventory | null>(null);

  const fetchData = async () => {
    let inventory: AccountInventory = await client.getAccountInventory();
    setAccountInventory(inventory);
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div className={content.page}>
      {accountInventory ? (
        <>
          <SearchInput handleUpdate={setAccountInventory}></SearchInput>
          <InventoryContents
            accountInventory={accountInventory}
          ></InventoryContents>
        </>
      ) : (
        <NoKeyPage />
      )}
    </div>
  );
};

interface SearchInputProps {
  handleUpdate: React.Dispatch<React.SetStateAction<AccountInventory | null>>;
}

const SearchInput: React.FC<SearchInputProps> = ({ handleUpdate }) => {
  const [formState, setFormState] = useState("");
  let context = useContext(ClientContext);
  let client = context;

  const handleChange = (e: React.FormEvent<HTMLInputElement>) => {
    const input = e.currentTarget.value;
    setFormState(input);
  };

  const handleSubmit = async (e: React.SyntheticEvent) => {
    e.preventDefault();
    const accountInventory = await client.postInventorySearch(formState);
    if (!accountInventory) {
      console.log("Could not post search");
    } else {
      handleUpdate(accountInventory);
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <label htmlFor="search-input">{"Search"}</label>
        <div>
          <input
            type="search"
            name="search-input"
            id="search"
            value={formState}
            onChange={handleChange}
          />
        </div>
        <input type="submit" value="Submit" />
      </form>
    </div>
  );
};
