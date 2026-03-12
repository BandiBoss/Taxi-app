package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, status int, userMsg string, err error) {
	if err != nil {
		log.Printf("HTTP %d: %s | details: %v", status, userMsg, err)
	} else {
		log.Printf("HTTP %d: %s", status, userMsg)
	}
	c.JSON(status, gin.H{"error": userMsg})
}
