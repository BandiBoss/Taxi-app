import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getOrderLocationHistory } from "../services/api";
import ArrowBackIosIcon from "@mui/icons-material/ArrowBackIosNew";

import {
  Box,
  Link,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
  Typography,
} from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

const columns = [
  { id: "index", label: "#", align: "center" },
  { id: "latitude", label: "Latitude", align: "center" },
  { id: "longitude", label: "Longitude", align: "center" },
  {
    id: "generated_time",
    label: "Time",
    align: "center",
    format: (value) => new Date(value).toLocaleString(),
  },
];

export default function OrderLocationHistory() {
  const { id } = useParams();
  const [history, setHistory] = useState([]);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  useEffect(() => {
    getOrderLocationHistory(id)
      .then((res) => setHistory(res.data.history || []))
      .catch(console.error);
  }, [id]);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  return (
    <Box sx={{ maxWidth: 900, mx: "auto", p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Location History for Order #{id}
      </Typography>
      <Box mb={2}>
        <Link
          component={RouterLink}
          to={`/orders/${id}`}
          underline="hover"
          sx={{ display: "flex", alignItems: "center" }}
        >
          <ArrowBackIosIcon size="small" />
          Back to Order Details
        </Link>
      </Box>

      <Paper sx={{ width: "100%", overflow: "hidden" }}>
        <TableContainer sx={{ maxHeight: 440 }}>
          <Table stickyHeader aria-label="location history table" size="small">
            <TableHead>
              <TableRow>
                {columns.map((column) => (
                  <TableCell
                    key={column.id}
                    align={column.align}
                    style={{ minWidth: 80 }}
                  >
                    {column.label}
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {history.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={columns.length} align="center">
                    No history found.
                  </TableCell>
                </TableRow>
              ) : (
                history
                  .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                  .map((row, index) => (
                    <TableRow key={index} hover>
                      {columns.map((column) => {
                        let value;
                        if (column.id === "index") {
                          value = page * rowsPerPage + index + 1;
                        } else {
                          value = row[column.id];
                        }
                        return (
                          <TableCell key={column.id} align={column.align}>
                            {column.format && value
                              ? column.format(value)
                              : value}
                          </TableCell>
                        );
                      })}
                    </TableRow>
                  ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        {history.length > 0 && (
          <TablePagination
            rowsPerPageOptions={[5, 10, 25]}
            component="div"
            count={history.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
          />
        )}
      </Paper>
    </Box>
  );
}
