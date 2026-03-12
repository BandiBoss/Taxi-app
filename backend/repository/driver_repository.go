package repository

import (
	"database/sql"
	"time"
)

type DriverLocation struct {
	ID            int     `json:"id"`
	DriverID      int     `json:"driver_id"`
	OrderID       int     `json:"order_id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	GeneratedTime string  `json:"generated_time"`
}

type Driver struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	CarModel     string `json:"car_model"`
	LicensePlate string `json:"license_plate"`
	IsActive     bool   `json:"is_active"`
}

type DriverRepository interface {
	GetDriverLocationHistory(driverID, limit, offset int) ([]DriverLocation, error)
	GetDrivers(limit, offset int) ([]Driver, error)
	AddDriver(d *Driver) error
	UpdateDriver(id int, d *Driver) error
	DeleteDriver(id int) error
	InsertDriverLocation(driverID, orderID int, lat, lon float64, generatedTime time.Time) error
	CountDrivers() (int, error)
}

type driverRepository struct {
	db *sql.DB
}

func NewDriverRepository(db *sql.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) GetDriverLocationHistory(driverID, limit, offset int) ([]DriverLocation, error) {
	rows, err := r.db.Query(`
        SELECT id, driver_id, order_id, latitude, longitude, generated_time
        FROM driver_locations
        WHERE driver_id = $1
        ORDER BY generated_time DESC
        LIMIT $2 OFFSET $3
    `, driverID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var locations []DriverLocation
	for rows.Next() {
		var loc DriverLocation
		if err := rows.Scan(&loc.ID, &loc.DriverID, &loc.OrderID, &loc.Latitude, &loc.Longitude, &loc.GeneratedTime); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func (r *driverRepository) CountDrivers() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM drivers").Scan(&count)
	return count, err
}

func (r *driverRepository) GetDrivers(limit, offset int) ([]Driver, error) {
	rows, err := r.db.Query("SELECT id, name, phone, car_model, license_plate, is_active FROM drivers ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var drivers []Driver
	for rows.Next() {
		var d Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.Phone, &d.CarModel, &d.LicensePlate, &d.IsActive); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (r *driverRepository) AddDriver(d *Driver) error {
	return r.db.QueryRow(
		"INSERT INTO drivers (name, phone, car_model, license_plate, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		d.Name, d.Phone, d.CarModel, d.LicensePlate, d.IsActive,
	).Scan(&d.ID)
}

func (r *driverRepository) UpdateDriver(id int, d *Driver) error {
	res, err := r.db.Exec(
		`UPDATE drivers SET name=$1, phone=$2, car_model=$3, license_plate=$4, is_active=$5 WHERE id=$6`,
		d.Name, d.Phone, d.CarModel, d.LicensePlate, d.IsActive, id,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *driverRepository) DeleteDriver(id int) error {
	res, err := r.db.Exec("DELETE FROM drivers WHERE id = $1", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *driverRepository) InsertDriverLocation(driverID, orderID int, lat, lon float64, generatedTime time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO driver_locations (driver_id, order_id, latitude, longitude, generated_time)
         VALUES ($1, $2, $3, $4, $5)`,
		driverID, orderID, lat, lon, generatedTime,
	)
	return err
}
