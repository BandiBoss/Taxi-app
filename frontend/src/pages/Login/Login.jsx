import Typography from "@mui/material/Typography";
import * as React from "react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import AuthForm from "../../components/AuthForm";
import { login } from "../../services/auth";


export default function SignIn() {
  const [username, setUserName] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const HelperComponent = () => (
    <Typography sx={{ textAlign: "center" }}>
      Don&apos;t have an account?{" "} 
      <Link to="/register" variant="body2" sx={{ alignSelf: "center" }}>
        Register
      </Link>
    </Typography>
  );

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const { role } = await login(username, password);
      if (role === "admin") {
        navigate("/admin/drivers");
      } else {
        navigate("/orders");
      }
    } catch (err) {
      alert("Login failed");
    }
  };

  return (
    <AuthForm
      ButtonText={"Login"}
      TitleText={"Login"}
      onSubmit={handleLogin}
      setUserName={setUserName}
      setPassword={setPassword}
      username={username}
      password={password}
      HelperComponent={HelperComponent}
    />
  );
}
