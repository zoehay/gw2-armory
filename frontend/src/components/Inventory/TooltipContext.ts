import { createContext } from "react";

export const TooltipContext = createContext<{
  activeId: string | null;
  setActiveId: (id: string | null) => void;
}>({ activeId: null, setActiveId: () => {} });
