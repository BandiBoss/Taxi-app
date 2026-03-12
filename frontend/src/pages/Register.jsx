import Alert from "@mui/material/Alert";
import Snackbar from "@mui/material/Snackbar";
import Typography from "@mui/material/Typography";
import * as React from "react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import AuthForm from "../components/AuthForm";
import { register } from "../services/api";

export default function Register() {
  const [username, setUserName] = useState("");
  const [password, setPassword] = useState("");
  const [isOpenSnackbar, setIsOpenSnackbar] = useState(false);
  const [snackbarMessage, setIsSnackbarMessage] = useState("");
  const navigate = useNavigate();

  const HelperComponent = () => (
    <Typography sx={{ textAlign: "center" }}>
      Already have an account?{" "}
      <Link to="/login" variant="body2" sx={{ alignSelf: "center" }}>
        Login
      </Link>
    </Typography>
  );
  const handleClose = () => {
    setIsOpenSnackbar(false);
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await register(username, password);
      alert("Registration successful. Please log in");
      navigate("/login");
    } catch (err) {
      if (err.response?.status === 409) {
        setIsOpenSnackbar(true);
        setIsSnackbarMessage(
          "That username is already taken. Please choose another"
        );
      } else {
        setIsOpenSnackbar(true);
        setIsSnackbarMessage("Registration failed. Please try again.");
      }
    }
  };

  return (
    <>
      <Snackbar
        open={isOpenSnackbar}
        autoHideDuration={6000}
        onClose={handleClose}
        anchorOrigin={{ vertical: "top", horizontal: "left" }}
      >
        <Alert severity="error">{snackbarMessage}.</Alert>
      </Snackbar>
      <AuthForm
        ButtonText={"Register"}
        TitleText={"Register"}
        onSubmit={handleRegister}
        setUserName={setUserName}
        setPassword={setPassword}
        username={username}
        password={password}
        HelperComponent={HelperComponent}
      />
    </>
  );
}
