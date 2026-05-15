import { Outlet } from "react-router-dom";
import { Navbar } from "./Navbar/Navbar";

export const Root = () => {
  return (
    <>
      <Navbar />
      <Outlet />
    </>
  );
};
