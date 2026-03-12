package handlers

import (
	"Taxi-app/backend/repository"
	"Taxi-app/backend/utils"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Driver struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	CarModel     string `json:"car_model"`
	LicensePlate string `json:"license_plate"`
	IsActive     bool   `json:"is_active"`
}

// GetDrivers returns a list of all drivers (admin only).
//
// @Summary      Get all drivers
// @Description  Get a list of all drivers (admin only)
// @Tags         admin
// @Produce      json
// @Success      200 {array} repository.Driver
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden"
// @Router       /api/admin/drivers [get]
// @Security     ApiKeyAuth
func GetDrivers(repo repository.DriverRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		limit := 50
		offset := (page - 1) * limit
		drivers, err := repo.GetDrivers(limit, offset)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "DB query failed", err)
			return
		}
		total, err := repo.CountDrivers()
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "DB count failed", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"drivers": drivers,
			"total":   total,
		})
	}
}

// AddDriver creates a new driver (admin only).
//
// @Summary      Add a driver
// @Description  Create a new driver (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        driver body repository.Driver true "Driver info"
// @Success      201 {object} repository.Driver
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden"
// @Router       /api/admin/drivers [post]
// @Security     ApiKeyAuth
func AddDriver(repo repository.DriverRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d repository.Driver
		if err := c.BindJSON(&d); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
			return
		}
		if err := repo.AddDriver(&d); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "DB query failed", err)
			return
		}
		c.JSON(http.StatusCreated, d)
	}
}

// UpdateDriver updates an existing driver (admin only).
//
// @Summary      Update a driver
// @Description  Update an existing driver (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path int true "Driver ID"
// @Param        driver body repository.Driver true "Driver info"
// @Success      200 {object} map[string]string "Driver updated"
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden"
// @Failure      404 {object} map[string]string "Driver not found"
// @Router       /api/admin/drivers/{id} [put]
// @Security     ApiKeyAuth
func UpdateDriver(repo repository.DriverRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid driver ID", err)
			return
		}
		var d repository.Driver
		if err := c.BindJSON(&d); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
			return
		}
		if err := repo.UpdateDriver(id, &d); err != nil {
			if err == sql.ErrNoRows {
				utils.ErrorResponse(c, http.StatusNotFound, "Driver not found", err)
			} else {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Update failed", err)
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Driver updated"})
	}
}

// DeleteDriver deletes a driver (admin only).
//
// @Summary      Delete a driver
// @Description  Delete a driver by ID (admin only)
// @Tags         admin
// @Produce      json
// @Param        id path int true "Driver ID"
// @Success      200 {object} map[string]string "Driver deleted"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden"
// @Failure      404 {object} map[string]string "Driver not found"
// @Router       /api/admin/drivers/{id} [delete]
// @Security     ApiKeyAuth
func DeleteDriver(repo repository.DriverRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid driver ID", err)
			return
		}
		if err := repo.DeleteDriver(id); err != nil {
			if err == sql.ErrNoRows {
				utils.ErrorResponse(c, http.StatusNotFound, "Driver not found", err)
			} else {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Delete failed", err)
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Driver deleted"})
	}
}

// DriverLocationHistory returns the location history for a driver (admin only).
//
// @Summary      Get driver location history
// @Description  Get location history for a driver (admin only)
// @Tags         admin
// @Produce      json
// @Param        id path int true "Driver ID"
// @Param        page query int false "Page number"
// @Param        limit query int false "Page size"
// @Success      200 {object} map[string][]repository.DriverLocation
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden"
// @Failure      404 {object} map[string]string "Driver not found"
// @Router       /api/admin/drivers/{id}/location-history [get]
// @Security     ApiKeyAuth
func DriverLocationHistory(repo repository.DriverRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		driverID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid driver ID", err)
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		limit := 50
		offset := (page - 1) * limit

		locations, err := repo.GetDriverLocationHistory(driverID, limit, offset)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Query failed", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	}
}
