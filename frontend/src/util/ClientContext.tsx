import React, { createContext } from "react";
import { Client } from "./Client";

const client = new Client();

export const ClientContext = createContext(client);

export const ClientProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <ClientContext.Provider value={client}>{children}</ClientContext.Provider>
  );
};
