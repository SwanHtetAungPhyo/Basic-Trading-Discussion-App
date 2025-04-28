package handlers

import (
	"github.com/SwanHtetAungPhyo/binance-dash/internal/services"
	"github.com/gofiber/contrib/websocket"
)

func WebsocketHandler(clientService *services.ClientService) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		client := services.NewClient(conn)
		clientService.AddClient(client)
		defer clientService.RemoveClient(client)
		
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}
