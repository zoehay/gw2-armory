import "./App.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { ManageKeys } from "./components/ManageKeys/ManageKeys";
import { ClientProvider } from "./util/ClientContext";
import { Inventory } from "./components/Inventory/Inventory";
import { Root } from "./components/Root";
import { ErrorPage } from "./components/ErrorPage/ErrorPage";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "/manageKeys",
        element: <ManageKeys />,
      },
      {
        path: "/inventory",
        element: <Inventory />,
      },
    ],
  },
]);

export function App() {
  return (
    <>
      <ClientProvider>
        <RouterProvider router={router} />
      </ClientProvider>
    </>
  );
}
