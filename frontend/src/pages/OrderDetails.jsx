import React, { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import {
  initiateConnection,
  closeConnection,
  subscribe,
  unsubscribe,
} from "../services/sockets";
import { getOrderDetails, simulateOrder } from "../services/api";
import OrderDetailsCard from "../components/OrderDetailsCard";

export default function OrderDetails() {
  const { id } = useParams();
  const [order, setOrder] = useState(null);
  const [location, setLocation] = useState(null);
  const didStart = useRef(false);

  useEffect(() => {
    getOrderDetails(id)
      .then((res) => setOrder(res.data))
      .catch((err) => console.error("Failed to load order", err));
  }, [id]);

  useEffect(() => {
    if (order?.status === "created" && !didStart.current) {
      didStart.current = true;

      simulateOrder(id)
        .then(() => {
          return getOrderDetails(id);
        })
        .then((res) => setOrder(res.data))
        .catch((err) => {
          if (err.response?.status !== 409) console.error(err);
        });
    }
  }, [order?.status, id]);

  useEffect(() => {
    const listener = (evt) => {
      try {
        const data = JSON.parse(evt.data);

        if (parseInt(data.order_id, 10) === parseInt(id, 10)) {
          setLocation({ lat: data.lat, lon: data.lon });
        }
        if (data.status) {
          setOrder((o) => ({ ...o, status: data.status }));
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    subscribe(listener);

    return () => {
      unsubscribe(listener);
    };
  }, [id]);

  useEffect(() => {
    if (order && order.status !== "done") {
      initiateConnection();
    }
  }, [order?.status]);

  useEffect(() => {
    return () => {
      closeConnection();
    };
  }, []);

  if (!order) return <div>Loading order...</div>;

  return (
    <>
      <OrderDetailsCard order={order} location={location} />
    </>
  );
}
