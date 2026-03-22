import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getDriverLocationHistory } from "../services/api";

import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Typography,
} from "@mui/material";

const columns = [
  { id: "index", label: "#", align: "center" },
  { id: "order_id", label: "Order ID", align: "center" },
  { id: "latitude", label: "Latitude", align: "center" },
  { id: "longitude", label: "Longitude", align: "center" },
  {
    id: "generated_time",
    label: "Generated Time",
    align: "center",
    format: (value) => new Date(value).toLocaleString(),
  },
];

export default function DriverLocationHistory() {
  const { id } = useParams();
  const [locations, setLocations] = useState([]);
  const [page, setPage] = useState(0);
  const rowsPerPage = 50; 

  useEffect(() => {
    const fetchLocations = async () => {
      try {
        const res = await getDriverLocationHistory(id, page + 1, rowsPerPage);
        setLocations(res.data.locations || []);
      } catch (err) {
        console.error(err);
      }
    };
    fetchLocations();
  }, [id, page]);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  return (
    <Box sx={{ maxWidth: 900, mx: "auto", p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Driver #{id} - Location History
      </Typography>

      <Paper sx={{ width: "100%", overflow: "hidden" }}>
        <TableContainer sx={{ maxHeight: 440 }}>
          <Table stickyHeader aria-label="driver location history" size="small">
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
              {locations.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={columns.length} align="center">
                    No location records found.
                  </TableCell>
                </TableRow>
              ) : (
                locations.map((row, index) => {
                  const displayIndex = page * rowsPerPage + index + 1;
                  return (
                    <TableRow key={index} hover>
                      {columns.map((column) => {
                        let value;
                        if (column.id === "index") {
                          value = displayIndex;
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
                  );
                })
              )}
            </TableBody>
          </Table>
        </TableContainer>

        <TablePagination
          rowsPerPageOptions={[rowsPerPage]}
          component="div"
          count={-1}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          labelRowsPerPage=""
          nextIconButtonProps={{
            disabled: locations.length < rowsPerPage,
            "aria-label": "Next page",
          }}
          backIconButtonProps={{
            disabled: page === 0,
            "aria-label": "Previous page",
          }}
        />
      </Paper>
    </Box>
  );
}
