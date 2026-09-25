import navbar from "./navbar.module.css";
import { useEffect, useState } from "react";
import { Link, NavLink } from "react-router-dom";

const linkClass = ({ isActive }: { isActive: boolean }) =>
  isActive ? `${navbar.link} ${navbar.active}` : navbar.link;

export const MobileNav = () => {
  const [showMenu, setShowMenu] = useState(false);

  useEffect(() => {
    if (!showMenu) return;
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setShowMenu(false);
    };
    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [showMenu]);

  const handleClose = () => {
    if (showMenu) {
      setShowMenu(false);
    }
  };

  return (
    <>
      <button
        type="button"
        className={navbar.menuButton}
        onClick={() => setShowMenu(!showMenu)}
        aria-label={showMenu ? "Close menu" : "Open menu"}
        aria-expanded={showMenu}
        aria-controls="mobile-menu"
      >
        <svg
          width="22"
          height="22"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          aria-hidden="true"
        >
          {showMenu ? (
            <path d="M6 6l12 12M18 6L6 18" />
          ) : (
            <path d="M4 7h16M4 12h16M4 17h16" />
          )}
        </svg>
      </button>
      <Link to="/" className={navbar.brand} onClick={handleClose}>
        armory
      </Link>
      <div className={showMenu ? navbar.open : undefined}>
        <div className={navbar.backdrop} onClick={handleClose} />
        <div
          id="mobile-menu"
          className={navbar.mobileDropdown}
          onClick={handleClose}
        >
          <NavLink to={`manageKeys`} className={linkClass}>
            Manage Keys
          </NavLink>
          <NavLink to={`inventory`} className={linkClass}>
            Inventory
          </NavLink>
        </div>
      </div>
    </>
  );
};
