package repository

import (
	"database/sql"
	"fmt"
)

type Order struct {
	ID          int    `json:"id"`
	CustomerID  int    `json:"customer_id"`
	DriverID    *int   `json:"driver_id"`
	Status      string `json:"status"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	CreatedAt   string `json:"created_at"`
}

type OrderRepository interface {
	CreateOrder(customerID int, origin, destination string) (int, error)
	GetOrders(customerID, limit, offset int, sortField, sortDirection string) ([]Order, error)
	GetOrderByID(orderID, customerID int) (*Order, error)
	GetOrderDetailsByID(orderID, customerID int) (*OrderDetails, error)
	GetOrderLocationHistory(orderID, customerID int, limit int) ([]Location, error)
	UpdateOrderStatus(orderID int, status string) error
	GetOrderStatus(orderID, customerID int) (string, error)
	GetRandomActiveDriver() (int, error)
	AssignDriverAndStart(orderID, driverID int) error
	ListActiveDrivers(limit int) ([]int, error)
	ListCreatedOrders(limit int) ([]int, error)
	SeedCreatedOrders(limit int) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(customerID int, origin, destination string) (int, error) {
	var orderID int
	err := r.db.QueryRow(`
        INSERT INTO orders (customer_id, status, origin, destination, created_at)
        VALUES ($1, 'created', $2, $3, NOW())
        RETURNING id
    `, customerID, origin, destination).Scan(&orderID)
	return orderID, err
}

func (r *orderRepository) GetOrders(customerID, limit, offset int, sortField, sortDirection string) ([]Order, error) {
	query := fmt.Sprintf(`
        SELECT id, customer_id, driver_id, status, origin, destination, created_at
        FROM orders
        WHERE customer_id = $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortField, sortDirection)

	rows, err := r.db.Query(query, customerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.DriverID, &o.Status, &o.Origin, &o.Destination, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

type Location struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	GeneratedTime string  `json:"generated_time"`
}

type OrderDetails struct {
	ID           int     `json:"id"`
	CustomerID   int     `json:"customer_id"`
	DriverID     *int    `json:"driver_id"`
	Status       string  `json:"status"`
	Origin       string  `json:"origin"`
	Destination  string  `json:"destination"`
	CreatedAt    string  `json:"created_at"`
	DriverName   *string `json:"driver_name"`
	LicensePlate *string `json:"license_plate"`
}

func (r *orderRepository) GetOrderByID(orderID, customerID int) (*Order, error) {
	row := r.db.QueryRow(`
        SELECT id, customer_id, driver_id, status, origin, destination, created_at
        FROM orders
        WHERE id = $1 AND customer_id = $2
    `, orderID, customerID)
	var o Order
	err := row.Scan(&o.ID, &o.CustomerID, &o.DriverID, &o.Status, &o.Origin, &o.Destination, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderRepository) GetOrderDetailsByID(orderID, customerID int) (*OrderDetails, error) {
	row := r.db.QueryRow(`
        SELECT
            o.id,
            o.customer_id,
            o.driver_id,
            o.status,
            o.origin,
            o.destination,
            o.created_at,
            d.name AS driver_name,
            d.license_plate
        FROM orders o
        LEFT JOIN drivers d ON o.driver_id = d.id
        WHERE o.id = $1 AND o.customer_id = $2
    `, orderID, customerID)
	var o OrderDetails
	err := row.Scan(
		&o.ID, &o.CustomerID, &o.DriverID, &o.Status, &o.Origin, &o.Destination, &o.CreatedAt,
		&o.DriverName, &o.LicensePlate,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderRepository) GetOrderLocationHistory(orderID, customerID int, limit int) ([]Location, error) {
	rows, err := r.db.Query(`
        SELECT dl.latitude, dl.longitude, dl.generated_time
        FROM driver_locations dl
        JOIN orders o ON dl.order_id = o.id
        WHERE dl.order_id = $1 AND o.customer_id = $2
        ORDER BY dl.generated_time DESC
        LIMIT $3
    `, orderID, customerID, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()
	var hist []Location
	for rows.Next() {
		var l Location
		if err := rows.Scan(&l.Latitude, &l.Longitude, &l.GeneratedTime); err != nil {
			return nil, err
		}
		hist = append(hist, l)
	}
	return hist, nil
}

func (r *orderRepository) UpdateOrderStatus(orderID int, status string) error {
	_, err := r.db.Exec(
		"UPDATE orders SET status = $1 WHERE id = $2",
		status, orderID,
	)
	return err
}

func (r *orderRepository) GetOrderStatus(orderID, customerID int) (string, error) {
	var status string
	err := r.db.QueryRow(
		"SELECT status FROM orders WHERE id=$1 AND customer_id=$2",
		orderID, customerID,
	).Scan(&status)
	return status, err
}

func (r *orderRepository) GetRandomActiveDriver() (int, error) {
	var driverID int
	err := r.db.QueryRow(
		"SELECT id FROM drivers WHERE is_active=true ORDER BY RANDOM() LIMIT 1",
	).Scan(&driverID)
	return driverID, err
}

func (r *orderRepository) AssignDriverAndStart(orderID, driverID int) error {
	_, err := r.db.Exec(
		"UPDATE orders SET status='in_progress', driver_id=$1 WHERE id=$2",
		driverID, orderID,
	)
	return err
}

func (r *orderRepository) ListActiveDrivers(limit int) ([]int, error) {
	rows, err := r.db.Query(
		`SELECT id FROM drivers WHERE is_active = true LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *orderRepository) ListCreatedOrders(limit int) ([]int, error) {
	rows, err := r.db.Query(
		`SELECT id FROM orders WHERE status = 'created' LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *orderRepository) SeedCreatedOrders(limit int) error {
	_, err := r.db.Exec(`
	INSERT INTO orders (customer_id, status, origin, destination, created_at)
        SELECT
          u.id,
          'created',
          'LoadTest Origin ' || u.id,
          'LoadTest Dest '   || u.id,
          NOW() - ((u.id % 30) * INTERVAL '1 day')
        FROM users u
        WHERE u.role = 'user'
        ORDER BY u.id
        LIMIT $1
	`, limit)
	return err
}
