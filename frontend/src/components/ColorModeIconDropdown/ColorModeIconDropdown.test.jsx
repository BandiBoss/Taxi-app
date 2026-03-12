import React from "react";
import { act, render, waitFor } from "@testing-library/react";
import { useColorScheme } from "@mui/material/styles";
import ColorModeIconDropdown from "./ColorModeIconDropdown";

jest.mock("@mui/material/styles");

describe("ColorModeIconDropdown", () => {
  beforeEach(() => {
    useColorScheme.mockReturnValue({
      mode: "light",
      systemMode: "light",
      setMode: jest.fn(),
    });
  });

  it("should render select", () => {
    const { getByTestId } = render(<ColorModeIconDropdown />);

    expect(getByTestId("toggle-mode")).toBeInTheDocument();
  });

  it("should open menu on icon button click", async () => {
    const { getByTestId, getByText } = render(<ColorModeIconDropdown />);
    const iconButton = getByTestId("toggle-mode");

    act(() => {
      iconButton.click();
    });

    await waitFor(() => {
      expect(getByText("Dark")).toBeVisible();
      expect(getByText("Light")).toBeVisible();
      expect(getByText("System")).toBeVisible();
    });
  });

  it("should change mode when menu item is clicked", async () => {
    const { getByTestId, getByText } = render(<ColorModeIconDropdown />);
    const iconButton = getByTestId("toggle-mode");

    act(() => {
      iconButton.click();
    });

    act(() => {
      getByText("Dark").click();
    });

    await waitFor(() => {
      expect(useColorScheme().setMode).toHaveBeenCalledWith("dark");
    });
  });
});
