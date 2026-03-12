import TextField from "@mui/material/TextField";
import { fireEvent, render } from "@testing-library/react";
import React from "react";
import SignIn from "./Login";

jest.mock("./Login", () => () => <div>Mocked Login</div>);

describe("SignIn", () => {
  it("renders without crashing", () => {
    const { getByText } = render(<SignIn />);
    expect(getByText("Mocked Login")).toBeInTheDocument();
  });
});

describe("MUI TextField", () => {
  it("should render and accept input", () => {
    const { getByLabelText } = render(
      <TextField label="Test Input" variant="outlined" />
    );

    const input = getByLabelText("Test Input");
    expect(input).toBeInTheDocument();

    fireEvent.change(input, { target: { value: "Hello, World!" } });
    expect(input.value).toBe("Hello, World!");
  });

  it("should display error message when error prop is set", () => {
    const { getByText } = render(
      <TextField
        label="Test Input"
        variant="outlined"
        error
        helperText="Error message"
      />
    );

    expect(getByText("Error message")).toBeInTheDocument();
  });
});
