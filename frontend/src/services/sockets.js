let connection = null;
let isConnecting = false;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectDelay = 2000;

const listeners = [];

const connect = () => {
  if (connection && connection.readyState === WebSocket.OPEN) {
    console.log("WebSocket already connected");
    return;
  }

  if (isConnecting) {
    console.log("WebSocket connection already in progress");
    return;
  }

  isConnecting = true;
  console.log("Attempting WebSocket connection...");

  const ws = new WebSocket(`ws://localhost:8080/api/ws`);

  ws.onopen = () => {
    console.log("WebSocket connected");
    connection = ws;
    isConnecting = false;
    reconnectAttempts = 0;
  };

  ws.onmessage = (event) => {
    try {
      listeners.forEach((listener) => listener(event));
    } catch (error) {
      console.error("Error in WebSocket message handler:", error);
    }
  };

  ws.onerror = (err) => {
    console.error("WebSocket error:", err);
    isConnecting = false;
  };

  ws.onclose = (event) => {
    console.log("WebSocket closed", event.code, event.reason);
    connection = null;
    isConnecting = false;

    if (event.code !== 1000 && reconnectAttempts < maxReconnectAttempts) {
      reconnectAttempts++;
      console.log(
        `Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts}) in ${reconnectDelay}ms...`
      );
      setTimeout(() => {
        connect();
      }, reconnectDelay);
    } else if (reconnectAttempts >= maxReconnectAttempts) {
      console.error("Max reconnection attempts reached");
    }
  };
};

export const initiateConnection = () => {
  connect();
};

export const subscribe = (listener) => {
  if (typeof listener === "function") {
    listeners.push(listener);
  }
};

export const unsubscribe = (listener) => {
  const index = listeners.indexOf(listener);
  if (index > -1) {
    listeners.splice(index, 1);
  }
};

export const closeConnection = () => {
  if (connection) {
    connection.close(1000, "Manual close");
    connection = null;
  }
  isConnecting = false;
  reconnectAttempts = 0;
};

export const getConnectionState = () => {
  if (!connection) return "CLOSED";
  switch (connection.readyState) {
    case WebSocket.CONNECTING:
      return "CONNECTING";
    case WebSocket.OPEN:
      return "OPEN";
    case WebSocket.CLOSING:
      return "CLOSING";
    case WebSocket.CLOSED:
      return "CLOSED";
    default:
      return "UNKNOWN";
  }
};
