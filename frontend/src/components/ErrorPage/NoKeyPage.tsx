import { Link } from "react-router-dom";
import content from "../content.module.css";

export const NoKeyPage = () => {
  return (
    <div className={content.card}>
      <p className={content.subtitle}>
        Add a Guild Wars 2 API key to see your inventory.{" "}
        <Link to="/manageKeys">Manage keys</Link>
      </p>
    </div>
  );
};
