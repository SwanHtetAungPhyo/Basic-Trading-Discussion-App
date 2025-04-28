package services

import (
	"context"
	"encoding/json"
	"github.com/SwanHtetAungPhyo/binance-dash/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

type SQSService struct {
	client        *sqs.Client
	clientService *ClientService
	queueURL      string
}

func NewSQSService(client *sqs.Client, clientService *ClientService) *SQSService {
	return &SQSService{
		client:        client,
		clientService: clientService,
		queueURL:      "https://sqs.eu-north-1.amazonaws.com/162047532564/Binance-ticker",
	}
}

func (s *SQSService) StartConsumer() {
	for {
		output, err := s.client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(s.queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("SQS receive error: %v", err)
			continue
		}

		for _, msg := range output.Messages {
			s.processMessage(msg)
		}
	}
}

func (s *SQSService) processMessage(msg types.Message) {
	var ticker models.TickerAndIndicator
	if err := json.Unmarshal([]byte(*msg.Body), &ticker); err != nil {
		log.Printf("message unmarshal error: %v", err)
		return
	}

	log.Printf("Received message: %v", ticker)
	s.clientService.Broadcast(&ticker)

	if _, err := s.client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	}); err != nil {
		log.Printf("message delete error: %v", err)
	}
}
