import React, { useEffect, useState } from "react";
import {
  getDrivers,
  addDriver,
  updateDriver,
  deleteDriver,
} from "../services/api";
import { Link } from "react-router-dom";

import {
  Box,
  Button,
  Checkbox,
  FormControlLabel,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
  Pagination,
} from "@mui/material";

export default function Drivers() {
  const [drivers, setDrivers] = useState([]);
  const [page, setPage] = useState(1);
  const [limit] = useState(50);
  const [total, setTotal] = useState(0);
  const [editingDriver, setEditingDriver] = useState(null);
  const [form, setForm] = useState({
    name: "",
    phone: "",
    car_model: "",
    license_plate: "",
    is_active: true,
  });

  const fetchDrivers = async () => {
    try {
      const res = await getDrivers(page, limit);
      setDrivers(res.data.drivers || res.data);
      setTotal(res.data.total || 0);
    } catch {
      alert("Failed to load drivers");
    }
  };

  useEffect(() => {
    fetchDrivers();
  }, [page]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setForm({ ...form, [name]: type === "checkbox" ? checked : value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingDriver) {
        await updateDriver(editingDriver.id, form);
      } else {
        await addDriver(form);
      }
      setForm({
        name: "",
        phone: "",
        car_model: "",
        license_plate: "",
        is_active: true,
      });
      setEditingDriver(null);
      fetchDrivers();
    } catch {
      alert("Failed to save driver");
    }
  };

  const handleEdit = (driver) => {
    setEditingDriver(driver);
    setForm(driver);
  };

  const handleDelete = async (id) => {
    if (window.confirm("Are you sure?")) {
      try {
        await deleteDriver(id);
        fetchDrivers();
      } catch {
        alert("Failed to delete driver");
      }
    }
  };

  return (
    <Box p={3}>
      <Typography variant="h4" mb={2}>
        Driver Management
      </Typography>

      <Box
        component="form"
        onSubmit={handleSubmit}
        mb={4}
        sx={{ display: "flex", flexDirection: "column", maxWidth: 400, gap: 2 }}
      >
        <Typography variant="h6">
          {editingDriver ? "Edit Driver" : "Add New Driver"}
        </Typography>
        <TextField
          label="Name"
          name="name"
          value={form.name}
          onChange={handleChange}
          required
          fullWidth
        />
        <TextField
          label="Phone"
          name="phone"
          value={form.phone}
          onChange={handleChange}
          required
          fullWidth
        />
        <TextField
          label="Car Model"
          name="car_model"
          value={form.car_model}
          onChange={handleChange}
          fullWidth
        />
        <TextField
          label="License Plate"
          name="license_plate"
          value={form.license_plate}
          onChange={handleChange}
          fullWidth
        />
        <FormControlLabel
          control={
            <Checkbox
              checked={form.is_active}
              onChange={handleChange}
              name="is_active"
            />
          }
          label="Active"
        />
        <Box display="flex" gap={2}>
          <Button variant="contained" type="submit">
            {editingDriver ? "Update" : "Add"} Driver
          </Button>
          {editingDriver && (
            <Button
              variant="outlined"
              onClick={() => {
                setEditingDriver(null);
                setForm({
                  name: "",
                  phone: "",
                  car_model: "",
                  license_plate: "",
                  is_active: true,
                });
              }}
            >
              Cancel
            </Button>
          )}
        </Box>
      </Box>

      <Typography variant="h6" mb={1}>
        Driver List
      </Typography>
      <TableContainer component={Paper}>
        <Table aria-label="drivers table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Phone</TableCell>
              <TableCell>Car Model</TableCell>
              <TableCell>License Plate</TableCell>
              <TableCell>Active</TableCell>
              <TableCell align="center" colSpan={3}>
                Actions
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {Array.isArray(drivers) && drivers.length > 0 ? (
              drivers.map((driver) => {
                const { id, name, phone, car_model, license_plate, is_active } =
                  driver;
                return (
                  <TableRow key={id}>
                    <TableCell>{id}</TableCell>
                    <TableCell>{name}</TableCell>
                    <TableCell>{phone}</TableCell>
                    <TableCell>{car_model}</TableCell>
                    <TableCell>{license_plate}</TableCell>
                    <TableCell>{is_active ? "Yes" : "No"}</TableCell>
                    <TableCell>
                      <Button
                        size="small"
                        color="primary"
                        variant="contained"
                        onClick={() => handleEdit(driver)}
                      >
                        Edit
                      </Button>
                    </TableCell>
                    <TableCell>
                      <Button
                        size="small"
                        variant="contained"
                        color="primary"
                        onClick={() => handleDelete(id)}
                      >
                        Delete
                      </Button>
                    </TableCell>
                    <TableCell>
                      <Button
                        color="primary"
                        variant="contained"
                        size="small"
                        component={Link}
                        to={`/admin/drivers/${id}/location-history`}
                      >
                        History
                      </Button>
                    </TableCell>
                  </TableRow>
                );
              })
            ) : (
              <TableRow>
                <TableCell colSpan={9} align="center">
                  No drivers found
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>
      <Pagination
        count={Math.ceil(total / limit) || 1}
        page={page}
        onChange={(e, value) => setPage(value)}
        sx={{ mt: 2 }}
      />
    </Box>
  );
}
