import { describe, expect, it } from "vitest";
import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { MobileNav } from "./MobileNav";
import { DesktopNav } from "./DesktopNav";

describe("MobileNav", () => {
  function renderMobileNav() {
    render(
      <MemoryRouter>
        <MobileNav />
      </MemoryRouter>,
    );
  }

  it("toggles the menu button state", async () => {
    const user = userEvent.setup();
    renderMobileNav();

    const button = screen.getByRole("button", { name: "Open menu" });
    expect(button).toHaveAttribute("aria-expanded", "false");

    await user.click(button);
    expect(button).toHaveAttribute("aria-expanded", "true");
    expect(button).toHaveAccessibleName("Close menu");

    await user.click(button);
    expect(button).toHaveAttribute("aria-expanded", "false");
  });

  it("closes on Escape", async () => {
    const user = userEvent.setup();
    renderMobileNav();

    await user.click(screen.getByRole("button", { name: "Open menu" }));
    await user.keyboard("{Escape}");

    expect(screen.getByRole("button", { name: "Open menu" })).toHaveAttribute(
      "aria-expanded",
      "false",
    );
  });

  it("closes when a link is chosen", async () => {
    const user = userEvent.setup();
    renderMobileNav();

    await user.click(screen.getByRole("button", { name: "Open menu" }));
    const menu = document.getElementById("mobile-menu")!;
    await user.click(within(menu).getByRole("link", { name: "Inventory" }));

    expect(screen.getByRole("button", { name: "Open menu" })).toHaveAttribute(
      "aria-expanded",
      "false",
    );
  });
});

describe("DesktopNav", () => {
  it("marks the current route's link as active", () => {
    render(
      <MemoryRouter initialEntries={["/inventory"]}>
        <DesktopNav />
      </MemoryRouter>,
    );

    expect(screen.getByRole("link", { name: "Inventory" })).toHaveAttribute(
      "aria-current",
      "page",
    );
    expect(
      screen.getByRole("link", { name: "Manage Keys" }),
    ).not.toHaveAttribute("aria-current");
    expect(screen.getByRole("link", { name: "armory" })).toHaveAttribute(
      "href",
      "/",
    );
  });
});
