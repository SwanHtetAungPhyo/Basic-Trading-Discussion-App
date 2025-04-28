package app

import (
	"context"
	"github.com/SwanHtetAungPhyo/binance-dash/internal/config"
	"github.com/SwanHtetAungPhyo/binance-dash/internal/handlers"
	"github.com/SwanHtetAungPhyo/binance-dash/internal/services"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type App struct {
	FiberApp        *fiber.App
	SQSClient       *sqs.Client
	ClientService   *services.ClientService
	dynamoDBService services.Dynamo
	handler         handlers.UserHandlers
	log             *logrus.Logger
	port            string
}

func NewApp() *App {
	logger := logrus.New()
	cfg, err := config.LoadAWSConfig(context.Background())
	if err != nil {
		logger.Fatalf("failed to load AWS config: %v", err)
	}

	clientService := services.NewClientService(&sync.Map{})
	sqsClient := sqs.NewFromConfig(*cfg)

	// Dynamo Services
	dynamoClient, err := config.GetDynamoDbClient()
	if err != nil {
		logger.Fatalf("failed to load DynamoDB client: %v", err)
	}
	dynamoServices := services.NewDynamoService(dynamoClient, logger)
	fiberApp := fiber.New()

	app := &App{
		FiberApp:        fiberApp,
		SQSClient:       sqsClient,
		ClientService:   clientService,
		log:             logger,
		dynamoDBService: dynamoServices,
	}

	app.setupRoutes()

	go services.NewSQSService(sqsClient, clientService).StartConsumer()

	return app
}

func (a *App) Run() error {
	a.log.Debug("Starting server...")
	port := viper.GetString("server.port")
	if port == "" {
		port = "8081"
	}
	return a.FiberApp.Listen(":" + port)
}

func (a *App) setupRoutes() {
	a.FiberApp.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			if c.Path() == "/ws/ticker" {
				if c.Get("access_token") == "" {
					return c.SendStatus(fiber.StatusUnauthorized)
				}
				return c.Next()
			}
		}
		return c.SendStatus(fiber.StatusNotFound)
	})
	a.FiberApp.Get("/ws/ticker", websocket.New(handlers.WebsocketHandler(a.ClientService)))
	a.FiberApp.Get("/health", handlers.HealthCheck)
}

func (a *App) setupUserRoutes() {
	user := a.FiberApp.Group("/user")
	user.Post("/login", a.handler.LoginHandler)
	user.Post("/register", a.handler.RegisterHandler)
	user.Post("/logout", a.handler.LogoutHandler)
}

func (a *App) Close() {
	if err := a.FiberApp.Shutdown(); err != nil {
		a.log.Fatal(err.Error())
	}
}
