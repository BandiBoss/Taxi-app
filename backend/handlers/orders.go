package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Taxi-app/backend/repository"
	"Taxi-app/backend/utils"

	"github.com/gin-gonic/gin"
)

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

// @Summary Create a new order
// @Description Create a new order with origin and destination
// @Tags orders
// @Accept json
// @Produce json
// @Param data body CreateOrderRequest true "Order details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/orders [post]
func CreateOrder(repo repository.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateOrderRequest
		if err := c.BindJSON(&req); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
			return
		}

		userID := c.MustGet("userID").(int)

		if req.Origin == "" {
			req.Origin = "Default Origin"
		}
		if req.Destination == "" {
			req.Destination = "Default Destination"
		}

		orderID, err := repo.CreateOrder(userID, req.Origin, req.Destination)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create order", err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Order created", "order_id": orderID})
	}
}

// @Summary Get user's orders
// @Description Get a paginated list of user's orders
// @Tags orders
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/orders [get]
func GetOrders(repo repository.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(int)

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		sortField := c.DefaultQuery("sort", "created_at")
		sortDirection := strings.ToUpper(c.DefaultQuery("order", "DESC"))

		switch sortField {
		case "created_at", "status":
		default:
			sortField = "created_at"
		}

		if sortDirection != "ASC" && sortDirection != "DESC" {
			sortDirection = "DESC"
		}

		// Logging for debugging
		logMsg := "[GetOrders] userID=%d, page=%d, limit=%d, offset=%d, sortField=%s, sortDirection=%s"
		fmt.Printf(logMsg+"\n", userID, page, limit, offset, sortField, sortDirection)

		orders, err := repo.GetOrders(userID, limit, offset, sortField, sortDirection)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Query failed", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":   page,
			"limit":  limit,
			"orders": orders,
		})
	}
}

// @Summary Get order details
// @Description Get details of a specific order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} repository.OrderDetails
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/orders/{id} [get]
func GetOrderByID(repo repository.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(int)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err)
			return
		}
		order, err := repo.GetOrderDetailsByID(id, userID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Order not found", err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// @Summary Get order location history
// @Description Get location history of a specific order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/orders/{id}/location-history [get]
func GetOrderLocationHistory(repo repository.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(int)
		orderIDStr := c.Param("id")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err)
			return
		}
		hist, err := repo.GetOrderLocationHistory(orderID, userID, 50)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Query failed", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"history": hist})
	}
}
