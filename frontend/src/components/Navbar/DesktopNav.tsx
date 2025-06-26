import { Link } from "react-router-dom";
import navbar from "./navbar.module.css";

export const DesktopNav = () => {
  return (
    <div className={navbar.desktopContent}>
      <p>Logo</p>
      <div className={navbar.desktopLinks}>
        <Link to={`manageKeys`} className={navbar.link}>
          Manage Keys
        </Link>
        <Link to={`inventory`} className={navbar.link}>
          Inventory
        </Link>
      </div>
    </div>
  );
};
