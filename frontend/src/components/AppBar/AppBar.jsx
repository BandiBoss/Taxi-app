import CloseRoundedIcon from "@mui/icons-material/CloseRounded";
import MenuIcon from "@mui/icons-material/Menu";
import { AppBar as MUIAppBar } from "@mui/material";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import Drawer from "@mui/material/Drawer";
import IconButton from "@mui/material/IconButton";
import { alpha, styled } from "@mui/material/styles";
import Toolbar from "@mui/material/Toolbar";
import * as React from "react";
import { Link } from "react-router-dom";
import logo from "../../assets/logo.png";
import { getRole, logout } from "../../services/auth";
import ColorModeSelect from "../../shared-theme/ColorModeSelect";

const StyledToolbar = styled(Toolbar)(({ theme }) => ({
  display: "flex",
  alignItems: "center",

  flexShrink: 0,
  borderRadius: `calc(${theme.shape.borderRadius}px + 8px)`,
  backdropFilter: "blur(24px)",
  border: "1px solid",
  borderColor: (theme.vars || theme).palette.divider,
  backgroundColor: theme.vars
    ? `rgba(${theme.vars.palette.background.defaultChannel} / 0.4)`
    : alpha(theme.palette.background.default, 0.4),
  boxShadow: (theme.vars || theme).shadows[1],
  padding: "8px 12px",
}));

export default function AppBar() {
  const [open, setOpen] = React.useState(false);
  const role = getRole();

  const toggleDrawer = (newOpen) => () => setOpen(newOpen);

  return (
    <MUIAppBar
      position="fixed"
      enableColorOnDark
      sx={{
        boxShadow: 0,
        bgcolor: "transparent",
        backgroundImage: "none",
        mt: "calc(var(--template-frame-height, 0px) + 28px)",
      }}
    >
      <Container maxWidth="lg">
        <StyledToolbar variant="dense" disableGutters>
          {/* ЛЕВЫЙ блок: логотип + кнопка */}
          <Box
            sx={{
              flexGrow: 1,
              display: "flex",
              alignItems: "center",
              gap: 2, // расстояние между лого и кнопками
            }}
          >
            {/* Сам логотип */}
            <Box
              component="img"
              src={logo}
              alt="Taxi Fast & Furious"
              sx={{
                height: 60, // подкорректируйте размер по вкусу
                objectFit: "contain",
              }}
            />

            {/* Навигационная кнопка */}
            <Box sx={{ display: { xs: "none", md: "flex" } }}>
              {role === "user" && (
                <Button
                  variant="text"
                  color="info"
                  size="small"
                  component={Link}
                  to="/orders"
                >
                  My Orders
                </Button>
              )}
              {role === "admin" && (
                <Button
                  variant="text"
                  color="info"
                  size="small"
                  component={Link}
                  to="/admin/drivers"
                >
                  Drivers
                </Button>
              )}
            </Box>
          </Box>

          {/* ПРАВЫЙ блок (десктоп): переключатель темы + выход */}
          <Box
            sx={{
              flexGrow: 1,
              display: { xs: "none", md: "flex" },
              alignItems: "center",
              justifyContent: "flex-end",
              gap: 1,
            }}
          >
            <ColorModeSelect />
            <Button
              color="primary"
              variant="contained"
              size="small"
              onClick={logout}
            >
              Logout
            </Button>
          </Box>

          {/* Drawer для мобильных */}
          <Box sx={{ display: { xs: "flex", md: "none" } }}>
            <IconButton aria-label="Menu" onClick={toggleDrawer(true)}>
              <MenuIcon />
            </IconButton>
            <Drawer
              anchor="top"
              open={open}
              onClose={toggleDrawer(false)}
              PaperProps={{
                sx: { top: "var(--template-frame-height, 0px)" },
              }}
            >
              <Box sx={{ p: 2, backgroundColor: "background.default" }}>
                <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
                  <IconButton onClick={toggleDrawer(false)}>
                    <CloseRoundedIcon />
                  </IconButton>
                </Box>
                {/* Здесь ваш контент для мобильного меню */}
              </Box>
            </Drawer>
          </Box>
        </StyledToolbar>
      </Container>
    </MUIAppBar>
  );
}
