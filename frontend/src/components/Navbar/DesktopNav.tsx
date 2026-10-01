import { Link, NavLink } from "react-router-dom";
import navbar from "./navbar.module.css";

const linkClass = ({ isActive }: { isActive: boolean }) =>
  isActive ? `${navbar.link} ${navbar.active}` : navbar.link;

export const DesktopNav = () => {
  return (
    <div className={navbar.desktopContent}>
      <Link to="/" className={navbar.brand}>
        armory
      </Link>
      <div className={navbar.desktopLinks}>
        <NavLink to={`manageKeys`} className={linkClass}>
          Manage Keys
        </NavLink>
        <NavLink to={`inventory`} className={linkClass}>
          Inventory
        </NavLink>
      </div>
    </div>
  );
};
