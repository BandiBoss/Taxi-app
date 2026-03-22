package handlers

import (
	"Taxi-app/backend/repository"
	"Taxi-app/backend/simulator"
	"Taxi-app/backend/utils"

	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

// @Summary Start order simulation
// @Description Start the simulation of a specific order
// @Tags simulator
// @Produce json
// @Param orderId path int true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/simulate/order/{orderId} [post]
func StartOrderSimulation(repo repository.OrderRepository, ch *amqp091.Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(int)
		
		orderID, err := strconv.Atoi(c.Param("orderId"))
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err)
			return
		}

		
		status, err := repo.GetOrderStatus(orderID, userID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Order not found", err)
			return
		}
		if status != "created" {
			utils.ErrorResponse(c, http.StatusConflict, "Order not in created state", nil)
			return
		}

		
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		driverID, err := repo.GetRandomActiveDriver()
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "No available drivers", err)
			return
		}

		
		if err := repo.AssignDriverAndStart(orderID, driverID); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update order status", err)
			return
		}

		
		go simulator.SimulateMovement(driverID, orderID, ch, "driver_updates", repo, r, 30, 1*time.Second)

		c.JSON(http.StatusOK, gin.H{
			"message":   "Simulation started",
			"order_id":  orderID,
			"driver_id": driverID,
		})
	}
}
