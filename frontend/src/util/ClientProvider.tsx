import React from "react";
import { client, ClientContext } from "./ClientContext";

export const ClientProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <ClientContext.Provider value={client}>{children}</ClientContext.Provider>
  );
};
