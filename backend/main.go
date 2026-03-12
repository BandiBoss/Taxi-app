// @title           Taxi App API
// @version         1.0
// @description     API documentation for the Taxi App backend.
// @contact.name    Bohdan Zghonnyk
// @contact.email   bogdanzn2002@gmail.com
// @host            localhost:8081
// @BasePath /api

package main

import (
	"Taxi-app/backend/consumer"
	"Taxi-app/backend/handlers"
	"Taxi-app/backend/middleware"

	_ "Taxi-app/backend/docs"
	"Taxi-app/backend/repository"
	"database/sql"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading configuration from environment")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	rabbitConn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Fatal("RabbitMQ connection failed:", err)
	}
	defer func() {
		if err := rabbitConn.Close(); err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}
	defer func() {
		if err := rabbitCh.Close(); err != nil {
			log.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}()

	r := gin.Default()

	r.OPTIONS("/*any", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	orderRepo := repository.NewOrderRepository(db)
	userRepo := repository.NewUserRepository(db)
	driverRepo := repository.NewDriverRepository(db)

	api := r.Group("/api")
	user := api.Group("/")
	user.Use(middleware.AuthMiddleware(), middleware.RequireRole("user"))
	{
		api.POST("/register", handlers.Register(userRepo))
		api.POST("/login", handlers.Login(userRepo))
		api.POST("/refresh", handlers.RefreshToken(userRepo))
		api.POST("/logout", handlers.Logout(userRepo))

		api.GET("/ws", handlers.WebSocketHandler)
		handlers.StartBroadcastLoop()
		consumer.StartConsumer(driverRepo, rabbitCh)

		user.POST("/orders", handlers.CreateOrder(orderRepo))
		user.GET("/orders", handlers.GetOrders(orderRepo))
		user.GET("/orders/:id", handlers.GetOrderByID(orderRepo))
		user.GET("/orders/:id/location-history", handlers.GetOrderLocationHistory(orderRepo))
		user.POST("/simulate/order/:orderId", handlers.StartOrderSimulation(orderRepo, rabbitCh))
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		admin.GET("/drivers", handlers.GetDrivers(driverRepo))
		admin.POST("/drivers", handlers.AddDriver(driverRepo))
		admin.PUT("/drivers/:id", handlers.UpdateDriver(driverRepo))
		admin.DELETE("/drivers/:id", handlers.DeleteDriver(driverRepo))
		admin.GET("/drivers/:id/location-history", handlers.DriverLocationHistory(driverRepo))

	}

	auth := api.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/profile", func(c *gin.Context) {
		userID := c.MustGet("userID").(int)
		c.JSON(http.StatusOK, gin.H{"message": "Hello User", "user_id": userID})
	})

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/build/index.html")
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to run server:", err)
	}

}
