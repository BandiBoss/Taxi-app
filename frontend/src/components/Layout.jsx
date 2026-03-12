import Container from "@mui/material/Container";
import CssBaseline from "@mui/material/CssBaseline";
import * as React from "react";
import AppTheme from "../shared-theme/AppTheme";
import AppBar from "./AppBar/AppBar";
import Footer from "./Footer";

export default function Layout({ children, ...props }) {
  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />

      <AppBar />
      <Container
        maxWidth="lg"
        component="main"
        sx={{ display: "flex", flexDirection: "column", my: 16, gap: 4 }}
      >
        {children}
      </Container>
      <Footer />
    </AppTheme>
  );
}
