import { createContext } from "react";
import { Client } from "./Client";

export const client = new Client();

export const ClientContext = createContext(client);
