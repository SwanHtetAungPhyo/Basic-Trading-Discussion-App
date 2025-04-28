package services

import (
	"github.com/SwanHtetAungPhyo/binance-dash/internal/models"
	"github.com/gofiber/contrib/websocket"

	"sync"
)

type Client struct {
	Conn *websocket.Conn
	IP   string
}

type ClientService struct {
	clients *sync.Map
}

func NewClientService(clients *sync.Map) *ClientService {
	return &ClientService{clients: clients}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		IP:   conn.IP(),
	}
}

func (cs *ClientService) AddClient(client *Client) {
	cs.clients.Store(client.IP, client)
}

func (cs *ClientService) RemoveClient(client *Client) {
	cs.clients.Delete(client.IP)
	_ = client.Conn.Close()
}

func (cs *ClientService) Broadcast(ticker *models.TickerAndIndicator) {
	cs.clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		if err := client.Conn.WriteJSON(ticker); err != nil {
			_ = client.Conn.Close()
			cs.clients.Delete(key)
		}
		return true
	})
}
