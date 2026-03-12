import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import * as React from "react";
import { useNavigate } from "react-router-dom";

export default function OrderDetailsCard({ order, location }) {
  const navigate = useNavigate();

  return (
    <>
      <Box
        sx={{
          maxWidth: 800,
          margin: "0 auto",
          p: 2,
          display: "flex",
          flexDirection: "column",
          gap: 3,
        }}
      >
        <Card variant="outlined" sx={{ width: "100%", p: 2 }}>
          <CardContent sx={{ textAlign: "center" }}>
            <Typography
              variant="h5"
              gutterBottom
            >{`Order #${order.id} Driver Info:`}</Typography>
            <Typography sx={{ color: "text.secondary", mb: 1.5 }}>
              <strong>ID:</strong> {order.driver_id ?? "—"}
            </Typography>
            <Typography sx={{ color: "text.secondary", mb: 1.5 }}>
              <strong>Name:</strong> {order.driver_name ?? "TBD"}
            </Typography>
            <Typography sx={{ color: "text.secondary", mb: 1.5 }}>
              <strong>Car License:</strong> {order.license_plate ?? "—"}
            </Typography>
          </CardContent>
        </Card>

        <Card sx={{ width: "100%", p: 2, textAlign: "center" }}>
          <Typography
            variant="h5"
            gutterBottom
          >{`Order status and location:`}</Typography>
          {order.status !== "done" ? (
            <>
              <Typography sx={{ mb: 1 }}>
                <strong>Status:</strong> {order.status}
              </Typography>
              {location ? (
                <Typography>
                  Latitude: <b>{location.lat}</b>
                  <br />
                  Longitude: <b>{location.lon}</b>
                </Typography>
              ) : (
                <Typography>Waiting for location...</Typography>
              )}
            </>
          ) : (
            <>
              <Typography sx={{ mb: 1 }}>
                <strong>Status:</strong> Done!
              </Typography>
              <Button
                size="small"
                onClick={() => navigate(`/orders/${order.id}/location-history`)}
              >
                View Ride Location History
              </Button>
            </>
          )}
        </Card>
      </Box>
    </>
  );
}
