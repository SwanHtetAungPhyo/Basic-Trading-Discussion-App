package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"sync"
)

// Client represents a connected WebSocket client
type Client struct {
	Conn *websocket.Conn
	IP   string
}

// TickerAndIndicator contains market data and technical indicators
type TickerAndIndicator struct {
	Ticker string `json:"ticker"`
	SMA    string `json:"sma"`
	RSI    string `json:"rsi"`
	EMA    string `json:"ema"`
	Time   string `json:"time"`
	Price  string `json:"price"`
}

// AppState holds a shared application state
type AppState struct {
	Clients sync.Map
	Config  *aws.Config
	SQS     *sqs.Client
}

func main() {
	appState := &AppState{
		Clients: sync.Map{},
	}

	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	appState.Config = &cfg
	appState.SQS = sqs.NewFromConfig(cfg)

	app := fiber.New()

	// WebSocket route handler
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		client := &Client{
			Conn: conn,
			IP:   conn.IP(),
		}

		// Store the client in the shared map
		appState.Clients.Store(client.IP, client)
		defer func() {
			// Clean up on disconnect
			appState.Clients.Delete(client.IP)
			if err := conn.Close(); err != nil {
				log.Printf("error closing WebSocket connection: %v", err)
			}
		}()

		// Keep reading from the WebSocket connection
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("error reading message: %v", err)
				break
			}
		}
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, from ec2 instance!",
		})
	})
	// Start SQS consumer in a goroutine
	go appState.startSQSConsumer()

	// Start the Fiber server on port 8081
	log.Fatal(app.Listen(":8081"))
}

// startSQSConsumer initiates the SQS consumption pipeline
func (a *AppState) startSQSConsumer() {
	queueURL := "https://sqs.eu-north-1.amazonaws.com/162047532564/Binance-ticker"

	// Loop for receiving messages from SQS
	for {
		output, err := a.SQS.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("SQS receive error: %v", err)
			continue
		}

		// Process each message from the queue
		for _, msg := range output.Messages {
			var ticker TickerAndIndicator
			if err := json.Unmarshal([]byte(*msg.Body), &ticker); err != nil {
				log.Printf("message unmarshal error: %v", err)
				continue
			}
			log.Printf("Received message: %v", ticker)
			a.broadcastToClients(&ticker)

			// Delete the message from the queue after processing
			if _, err := a.SQS.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			}); err != nil {
				log.Printf("message delete error: %v", err)
			}
		}
	}
}

// broadcastToClients sends the SQS message to all connected WebSocket clients
func (a *AppState) broadcastToClients(ticker *TickerAndIndicator) {
	a.Clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		if err := client.Conn.WriteJSON(ticker); err != nil {
			log.Printf("broadcast error to %s: %v", client.IP, err)
			if err := client.Conn.Close(); err != nil {
				log.Printf("error closing connection: %v", err)
			}
			a.Clients.Delete(key)
		}
		return true
	})
}
