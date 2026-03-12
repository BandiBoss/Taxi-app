import ArrowBackIosIcon from "@mui/icons-material/ArrowBackIosNew";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import * as React from "react";
import { useNavigate } from "react-router-dom";

const columns = [
  { id: "id", label: "ID", align: "center" },
  { id: "status", label: "Status", align: "center" },
  { id: "driver_id", label: "Driver", align: "center" },
  { id: "origin", label: "Origin", align: "center" },
  { id: "destination", label: "Destination", align: "center" },
  {
    id: "created_at",
    label: "Created At",
    align: "center",
    format: (value) => new Date(value).toLocaleString(),
  },
];

export default function StickyHeadTable({
  orders,
  nexPage,
  page,
  limit,
  prevPage,
}) {
  const navigate = useNavigate();

  const handleRowClick = (row) => {
    navigate(`/orders/${row.id}`);
  };

  return (
    <Paper sx={{ width: "100%", overflow: "hidden" }}>
      <TableContainer sx={{ maxHeight: 440 }}>
        <Table stickyHeader aria-label="sticky table" size="small">
          <TableHead>
            <TableRow>
              {columns.map((column) => (
                <TableCell
                  key={column.id}
                  align={column.align}
                  style={{ minWidth: 100 }}
                >
                  {column.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {orders.map((row) => (
              <TableRow
                hover
                role="checkbox"
                tabIndex={-1}
                key={row.id}
                onClick={() => {
                  handleRowClick(row);
                }}
              >
                {columns.map((column) => {
                  const value = row[column.id];
                  return (
                    <TableCell key={column.id} align={column.align}>
                      {column.format && value ? column.format(value) : value}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Box sx={{ display: "flex", justifyContent: "end", padding: "1rem" }}>
        <IconButton
          aria-label="prev"
          onClick={prevPage}
          disabled={page === 1}
          size="small"
        >
          <ArrowBackIosIcon />
        </IconButton>

        <span style={{ margin: "0.5rem 1rem" }}>Page {page}</span>
        <IconButton
          aria-label="next"
          size="small"
          onClick={nexPage}
          disabled={orders.length < limit}
        >
          <ArrowForwardIosIcon />
        </IconButton>
      </Box>
    </Paper>
  );
}
