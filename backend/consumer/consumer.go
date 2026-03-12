package consumer

import (
	"Taxi-app/backend/handlers"
	"Taxi-app/backend/repository"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type LocationUpdate struct {
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

func StartConsumer(repo repository.DriverRepository, ch *amqp091.Channel) {
	q, err := ch.QueueDeclare("driver_updates", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Queue declare failed: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Consumer registration failed: %v", err)
	}

	go func() {
		for d := range msgs {
			var update LocationUpdate
			if err := json.Unmarshal(d.Body, &update); err != nil {
				log.Printf("JSON parse error: %v", err)
				continue
			}

			var lat, lon float64
			if _, err := fmt.Sscanf(update.Cordinates.Lat, "%f", &lat); err != nil {
				log.Printf("Failed to parse latitude: %v", err)
				continue
			}
			if _, err := fmt.Sscanf(update.Cordinates.Lon, "%f", &lon); err != nil {
				log.Printf("Failed to parse longitude: %v", err)
				continue
			}

			driverID, err1 := strconv.Atoi(update.DriverID)
			orderID, err2 := strconv.Atoi(update.Order.ID)
			if err1 != nil || err2 != nil {
				log.Printf("Failed to parse driverID or orderID: %v, %v", err1, err2)
				continue
			}
			if err := repo.InsertDriverLocation(driverID, orderID, lat, lon, time.Now()); err != nil {
				log.Printf("DB insert error: %v", err)
				continue
			}

			handlers.BroadcastLocationUpdate(map[string]interface{}{
				"order_id":  update.Order.ID,
				"driver_id": update.DriverID,
				"lat":       lat,
				"lon":       lon,
				"status":    update.Order.Status,
			})
		}
	}()
}
