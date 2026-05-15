import React, { createContext, useEffect, useState } from "react";

export const TooltipContext = createContext<{
  activeId: string | null;
  setActiveId: (id: string | null) => void;
}>({ activeId: null, setActiveId: () => {} });

export const TooltipProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [activeId, setActiveId] = useState<string | null>(null);

  useEffect(() => {
    if (activeId === null) return;

    const handleTouchStart = (e: TouchEvent) => {
      const target = e.target as Element;
      if (!target.closest("[data-inventory-tile]")) {
        setActiveId(null);
      }
    };

    document.addEventListener("touchstart", handleTouchStart);
    return () => document.removeEventListener("touchstart", handleTouchStart);
  }, [activeId]);

  return (
    <TooltipContext.Provider value={{ activeId, setActiveId }}>
      {children}
    </TooltipContext.Provider>
  );
};
