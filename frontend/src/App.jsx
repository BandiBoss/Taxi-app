import { ThemeProvider, createTheme } from '@mui/material/styles';
import React from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import Layout from "./components/Layout";
import DriverLocationHistory from "./pages/DriverLocationHistory";
import Drivers from "./pages/Drivers";
import Login from "./pages/Login/Login";
import OrderDetails from "./pages/OrderDetails";
import Orders from "./pages/Orders";
import OrderLocationHistory from "./pages/OrdersLocationHistory";
import Register from "./pages/Register";
import { getRole, isLoggedIn } from "./services/auth";

function ProtectedLayout() {
  const loggedIn = isLoggedIn();
  const role = getRole();

  if (!loggedIn) {
    return <Navigate to="/login" replace />;
  }

  return (
    <>

      <Layout>
        <Routes>
          <Route
            path="/orders"
            element={role === "user" ? <Orders /> : <Navigate to="/login" replace />}
          />
          <Route
            path="/orders/:id"
            element={role === "user" ? <OrderDetails /> : <Navigate to="/login" replace />}
          />
          <Route
            path="/orders/:id/location-history"
            element={role === "user" ? <OrderLocationHistory /> : <Navigate to="/login" replace />}
          />
          <Route
            path="/admin/drivers"
            element={role === "admin" ? <Drivers /> : <Navigate to="/login" replace />}
          />
          <Route
            path="/admin/drivers/:id/location-history"
            element={role === "admin" ? <DriverLocationHistory /> : <Navigate to="/login" replace />}
          />
          <Route
            path="*"
            element={<Navigate to={role === "admin" ? "/admin/drivers" : "/orders"} replace />}
          />
        </Routes>
      </Layout>
    </>
  );
}

function App() {
  const theme = createTheme();
  const loggedIn = isLoggedIn();
  const role = getRole();

  return (
    <ThemeProvider theme={theme}>
      <Routes>
        <Route
          path="/login"
          element={loggedIn ? <Navigate to={role === "admin" ? "/admin/drivers" : "/orders"} /> : <Login />}
        />
        <Route
          path="/register"
          element={loggedIn ? <Navigate to="/orders" /> : <Register />}
        />

        {/* Все защищённые страницы внутри ProtectedLayout */}
        <Route path="/*" element={<ProtectedLayout />} />
      </Routes>
    </ThemeProvider>
  );
}

export default App;
