import axios from "axios";
import { logout } from "./auth";


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

const TOKEN_KEY = "taxi_token";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  const token = getCookie(TOKEN_KEY);
  if (token) {
    config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const original = err.config;
    if (err.response?.status === 401 && !original._retry) {
      if (isRefreshing) {
        return new Promise(function (resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            original.headers["Authorization"] = `Bearer ${token}`;
            return api(original);
          })
          .catch((error) => Promise.reject(error));
      }
      original._retry = true;
      isRefreshing = true;
      try {
        const refreshRes = await api.post(
          "/refresh",
          {},
          { withCredentials: true }
        );
        const newToken = refreshRes.data.access_token;
        setCookie(TOKEN_KEY, newToken, 1); 
        api.defaults.headers["Authorization"] = `Bearer ${newToken}`;
        original.headers["Authorization"] = `Bearer ${newToken}`;
        processQueue(null, newToken);
        return api(original);
      } catch (refreshErr) {
        processQueue(refreshErr, null);
        eraseCookie(TOKEN_KEY);
        logout();
        return Promise.reject(refreshErr);
      } finally {
        isRefreshing = false;
      }
    }
    return Promise.reject(err);
  }
);

// Driver API
export const getDrivers = (page = 1, limit = 50) =>
  api.get("/admin/drivers", { params: { page, limit } });
export const addDriver = (data) => api.post("/admin/drivers", data);
export const updateDriver = (id, data) => api.put(`/admin/drivers/${id}`, data);
export const deleteDriver = (id) => api.delete(`/admin/drivers/${id}`);
export const getDriverLocationHistory = (id, page, limit) =>
  api.get(`/admin/drivers/${id}/location-history`, { params: { page, limit } });

// Order API
export const getOrders = (page, limit, sortField, sortDirection) =>
  api.get("/orders", {
    params: { page, limit, sort: sortField, order: sortDirection },
  });
export const createOrder = (data) => api.post("/orders", data);
export const getOrderDetails = (id) => api.get(`/orders/${id}`);
export const simulateOrder = (id) => api.post(`/simulate/order/${id}`);
export const getOrderLocationHistory = (orderId) =>
  api.get(`/orders/${orderId}/location-history`);

//Auth API
export const register = (username, password) =>
  api.post("/register", { username, password }, { withCredentials: true });
export const login = (username, password) =>
  api.post("/login", { username, password });
export const logoutApi = () => api.post("/logout");

export default api;
