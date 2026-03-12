import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import Select from "@mui/material/Select";
import TextField from "@mui/material/TextField";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Card } from "../components/OrderCard";
import StickyHeadTable from "../components/Table";
import { createOrder, getOrders } from "../services/api";

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [page, setPage] = useState(1);
  const limit = 10;
  const [sortField, setSortField] = useState("created_at");
  const [sortDirection, setSortDirection] = useState("desc");
  const [origin, setOrigin] = useState("");
  const [destination, setDestination] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    fetchOrders();
  }, [page, sortField, sortDirection]);

  async function fetchOrders() {
    try {
      const res = await getOrders(page, limit, sortField, sortDirection);
      setOrders(res.data.orders || []);
    } catch (err) {
      console.log(err);
    }
  }

  const toggleOrder = () => {
    setSortDirection((prev) => (prev === "asc" ? "desc" : "asc"));
  };

  const handleSortChange = (e) => {
    setSortField(e.target.value);
    toggleOrder();
  };

  const handleCreateOrder = async () => {
    try {
      const res = await createOrder({
        customer_id: 2,
        origin,
        destination,
      });
      const newOrderId = res.data.order_id;
      setOrigin("");
      setDestination("");
      navigate(`/orders/${newOrderId}`);
    } catch (err) {
      console.error(err);
    }
  };

  const nexPage = () => {
    setPage((p) => p + 1);
  };
  const prevPage = () => {
    setPage((p) => p - 1);
  };

  return (
    <div>
      <Card variant="outlined">
        <h4>Order a Taxi</h4>
        <Box
          component="form"
          sx={{
            "& > :not(style)": { m: 1, width: "25ch" },
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
          }}
          noValidate
          autoComplete="off"
        >
          <TextField
            id="Origin"
            label="Origin"
            value={origin}
            variant="standard"
            onChange={(e) => setOrigin(e.target.value)}
          />
          <TextField
            id="Destination"
            label="Destination"
            value={destination}
            variant="standard"
            onChange={(e) => setDestination(e.target.value)}
          />
        </Box>
        <Box sx={{ display: "flex", justifyContent: "center" }}>
          {origin && destination && (
            <Button onClick={handleCreateOrder} variant="contained">
              Order a Taxi
            </Button>
          )}
        </Box>
      </Card>
      <h2>Order List</h2>
      <div>
        <label>Sort by: </label>
        <Select
          labelId="demo-simple-select-standard-label"
          id="demo-simple-select-standard"
          value={sortField}
          onChange={handleSortChange}
          label="Sort by"
          sx={{ minWidth: 150, ml: 1, mr: 2 }}
        >
          <MenuItem value="created_at">Created At</MenuItem>
          <MenuItem value="status">Status</MenuItem>
        </Select>
      </div>

      <Box>
        <StickyHeadTable
          orders={orders}
          nexPage={nexPage}
          page={page}
          limit={limit}
          prevPage={prevPage}
        />
      </Box>
    </div>
  );
}
