package simulator

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"Taxi-app/backend/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimulatedMessage struct {
	DriverID   string `json:"driver_id"`
	Cordinates struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	} `json:"cordinates"`
	Order struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"order"`
}

func SimulateMovement(driverID, orderID int, ch *amqp.Channel, queue string, repo repository.OrderRepository, r *rand.Rand, points int, delay time.Duration) {
	for i := 0; i < points; i++ {
		lat := 47.84 + r.Float64()*0.01
		lon := 34.85 + r.Float64()*0.01

		msg := SimulatedMessage{
			DriverID: fmt.Sprintf("%d", driverID),
			Order: struct {
				ID     string `json:"id"`
				Status string `json:"status"`
			}{
				ID:     fmt.Sprintf("%d", orderID),
				Status: "in_progress",
			},
		}
		msg.Cordinates.Lat = fmt.Sprintf("%.6f", lat)
		msg.Cordinates.Lon = fmt.Sprintf("%.6f", lon)

		body, _ := json.Marshal(msg)
		_ = ch.Publish("", queue, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		log.Printf("Sent location for driver %d (order %d): %s,%s", driverID, orderID, msg.Cordinates.Lat, msg.Cordinates.Lon)
		time.Sleep(delay)
	}

	doneMsg := SimulatedMessage{
		DriverID: fmt.Sprintf("%d", driverID),
		Order: struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}{
			ID:     fmt.Sprintf("%d", orderID),
			Status: "done",
		},
	}

	doneMsg.Cordinates.Lat = "47.850000"
	doneMsg.Cordinates.Lon = "34.860000"

	body, _ := json.Marshal(doneMsg)
	_ = ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	_ = repo.UpdateOrderStatus(orderID, "done")
	log.Printf("Simulation finished for driver %d (order %d)", driverID, orderID)
}
