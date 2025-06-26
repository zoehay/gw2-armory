import navbar from "./navbar.module.css";
import { useState } from "react";
import { Link } from "react-router-dom";

export const MobileNav = () => {
  let [showMenu, setShowMenu] = useState(false);
  let position;

  if (!showMenu) {
    position = "-100%";
  } else {
    position = "0%";
  }
  let style = { left: position } as React.CSSProperties;

  const handleClose = () => {
    if (showMenu) {
      setShowMenu(false);
    }
  };

  return (
    <>
      <button onClick={() => setShowMenu(!showMenu)}> </button>
      <div
        className={navbar.mobileDropdownWrapper}
        onClick={handleClose}
        style={style}
      >
        <div className={navbar.mobileDropdown}>
          <Link to={`manageKeys`} className={navbar.link}>
            Manage Keys
          </Link>
          <Link to={`inventory`} className={navbar.link}>
            Inventory
          </Link>
        </div>
      </div>
    </>
  );
};
