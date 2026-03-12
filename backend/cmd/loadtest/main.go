package main

import (
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"sync"
	"time"

	"Taxi-app/backend/repository"
	"Taxi-app/backend/simulator"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	numDrivers = flag.Int("drivers", 1000, "how many drivers to simulate concurrently")
	numUpdates = flag.Int("updates", 50, "number of location updates per order")
	interval   = flag.Duration("interval", 5*time.Second, "delay between updates")
)

func main() {
	flag.Parse()

	start := time.Now()
	// 1) Open DB & AMQP
	db, err := sql.Open("postgres", "user=postgres password=12345 dbname=taxidb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing AMQP connection: %v", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("Error closing AMQP channel: %v", err)
		}
	}()

	repo := repository.NewOrderRepository(db)
	if err := repo.SeedCreatedOrders(*numDrivers); err != nil {
		log.Fatalf("could not seed created orders: %v", err)
	}

	// 2) Grab DRIVERS and ORDERS
	drivers, err := repo.ListActiveDrivers(*numDrivers)
	if err != nil {
		log.Fatal(err)
	}
	orders, err := repo.ListCreatedOrders(*numDrivers)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	launched := 0

	for i, driverID := range drivers {
		if i >= len(orders) {
			break
		}
		orderID := orders[i]

		// assign and flip to in_progress
		if err := repo.AssignDriverAndStart(orderID, driverID); err != nil {
			log.Printf("assign error for order %d driver %d: %v", orderID, driverID, err)
			continue
		}

		launched++
		wg.Add(1)

		// each goroutine runs the simulation then marks done
		go func(d, o int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			simulator.SimulateMovement(
				d,
				o,
				ch,
				"driver_updates",
				repo,
				r,
				*numUpdates,
				*interval,
			)
		}(driverID, orderID)
	}

	// 6) Wait for all sims to finish
	wg.Wait()

	// 7) Log the summary
	elapsed := time.Since(start).Seconds()
	log.Printf(
		"Load test passed correctly. Number of drivers simulated: %d. Number of orders finished: %d. Duration: %.2f seconds.",
		launched, launched, elapsed,
	)
}
