package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler upgrades the connection to a WebSocket for real-time updates.
//
// @Summary      WebSocket for real-time updates
// @Description  Upgrade to a WebSocket connection for receiving real-time order/driver updates.
// @Tags         websocket
// @Produce      json
// @Success      101 {string} string "Switching Protocols"
// @Router       /api/ws [get]
func WebSocketHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func() { _ = ws.Close() }()

	clients[ws] = true
	log.Println("WebSocket connected")

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			log.Println("WebSocket disconnected")
			break
		}
	}
}

func StartBroadcastLoop() {
	go func() {
		for msg := range broadcast {
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("WebSocket write failed:", err)
					_ = client.Close()
					delete(clients, client)
				}
			}
		}
	}()
}

func BroadcastLocationUpdate(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal WebSocket message:", err)
		return
	}
	broadcast <- jsonData
}
