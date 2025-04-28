package services

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

type Dynamo interface {
	PutUser(table string, user map[string]types.AttributeValue) error
	DeleteUser(table string, user map[string]types.AttributeValue) error
	GetUser(table string, user map[string]types.AttributeValue) (map[string]types.AttributeValue, error)
}
type DynamoService struct {
	log    *logrus.Logger
	client *dynamodb.Client
}

func NewDynamoService(client *dynamodb.Client, log *logrus.Logger) *DynamoService {
	return &DynamoService{
		log:    log,
		client: client,
	}
}

func (d DynamoService) PutUser(table string, user map[string]types.AttributeValue) error {
	//TODO implement me
	panic("implement me")
}

func (d DynamoService) DeleteUser(table string, user map[string]types.AttributeValue) error {
	//TODO implement me
	panic("implement me")
}

func (d DynamoService) GetUser(table string, user map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	//TODO implement me
	panic("implement me")
}
