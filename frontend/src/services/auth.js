import { jwtDecode } from "jwt-decode";
import { login as loginApi, logoutApi } from "./api";

const TOKEN_KEY = "taxi_token";
const ROLE_KEY = "taxi_role";
const USERID_KEY = "taxi_user_id";

// Cookie helpers
function setCookie(name, value, days) {
  let expires = "";
  if (days) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

function getCookie(name) {
  const nameEQ = name + "=";
  const ca = document.cookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === " ") c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

function eraseCookie(name) {
  document.cookie = name + "=; Max-Age=-99999999; path=/";
}

export async function login(username, password) {
  const res = await loginApi(username, password);
  const { access_token: token } = res.data;
  setCookie(TOKEN_KEY, token, 1); // 1 day expiry, adjust as needed
  const { user_id, role } = jwtDecode(token);
  localStorage.setItem(USERID_KEY, user_id);
  localStorage.setItem(ROLE_KEY, role);
  return { user_id, role, token };
}

export async function logout() {
  eraseCookie(TOKEN_KEY);
  localStorage.removeItem(USERID_KEY);
  localStorage.removeItem(ROLE_KEY);
  localStorage.setItem("auth", false);
  try {
    await logoutApi();
  } catch (e) {
    // Intentionally ignore logout errors
  }
  window.location.href = "/login";
}

export function getToken() {
  return getCookie(TOKEN_KEY);
}
export function getRole() {
  return localStorage.getItem(ROLE_KEY);
}
export function getUserID() {
  const v = localStorage.getItem(USERID_KEY);
  return v ? parseInt(v, 10) : null;
}
export function isLoggedIn() {
  return !!getToken();
}
