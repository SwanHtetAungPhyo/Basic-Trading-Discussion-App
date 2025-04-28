package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/markcheno/go-talib"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type TickerAndIndicator struct {
	Ticker *binance.WsKlineEvent
	Sma    string
	Rsi    string
	Ema    string
	Time   string
	Price  string
}

//type MACD struct {
//	m []float64
//	a []
//}

func main() {
	//aws := LoadAWSConfig()
	//client := SQSClient(aws)
	closingPrices := make([]float64, 0)
	period := 14
	minDataPoint := period + 1
	wsKlineHandler := func(event *binance.WsKlineEvent) {

		converted, err := strconv.ParseFloat(event.Kline.Close, 64)
		if err != nil {
			fmt.Println("Error parsing closing price:", err)
			return
		}
		closingPrices = append(closingPrices, converted)
		if len(closingPrices) > 100 {
			closingPrices = closingPrices[len(closingPrices)-100:]
		}
		fmt.Println("Closing price:", converted)
		if len(closingPrices) >= minDataPoint {
			latest := closingPrices[len(closingPrices)-period:]

			sma := talib.Sma(latest, period)
			rsi := talib.Rsi(latest, period-1)
			ema := talib.Ema(latest, period)
			//m, a, c := talib.Macd(latest, 12, 26, 9)
			var smaValue, rsiValue, emaValue float64

			if len(sma) > 0 {
				smaValue = sma[len(sma)-1]
			}
			if len(rsi) > 0 {
				rsiValue = rsi[len(rsi)-1]
			}
			if len(ema) > 0 {
				emaValue = ema[len(ema)-1]
			}

			ticker := &TickerAndIndicator{
				Ticker: event,
				Time:   strconv.FormatInt(event.Time, 10),
				Sma:    fmt.Sprintf("%.2f", smaValue),
				Rsi:    fmt.Sprintf("%.2f", rsiValue),
				Ema:    fmt.Sprintf("%.2f", emaValue),
				Price:  fmt.Sprintf("%.2f", converted),
			}
			jsonMessage, err := json.MarshalIndent(ticker, "", "  ")
			if err != nil {
				log.Fatal("Error encoding message:", err.Error())
				return
			}
			log.Println(string(jsonMessage))
			//SendingMessage(client, ticker)
		}
	}

	errHandler := func(err error) {
		fmt.Println("Error in WebSocket handler:", err)
	}

	doneC, _, err := binance.WsKlineServe("BTCUSDT", "1m", wsKlineHandler, errHandler)
	if err != nil {
		log.Fatal("Error starting WebSocket:", err)
		return
	}

	<-doneC
}

func LoadAWSConfig() *aws.Config {
	defaultConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal("Error loading AWS config:", err)
		return nil
	}
	return &defaultConfig
}

func SQSClient(cfg *aws.Config) *sqs.Client {
	return sqs.NewFromConfig(*cfg)
}

func SendingMessage(client *sqs.Client, payload *TickerAndIndicator) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error encoding payload:", err)
		return
	}
	stringPayload := string(jsonPayload)
	queueUrl := "https://sqs.eu-north-1.amazonaws.com/162047532564/Binance-ticker"
	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(stringPayload),
	}

	msg, err := client.SendMessage(context.Background(), sendMessageInput)
	if err != nil {
		log.Fatal("Error sending message to SQS:", err)
		return
	}

	err = json.NewEncoder(os.Stdout).Encode(msg)
	if err != nil {
		log.Fatal("Error encoding response:", err)
		return
	}
}

func formatSma(sma []float64) string {
	var result []string
	for _, v := range sma {
		result = append(result, fmt.Sprintf("%.2f", v))
	}
	return strings.Join(result, ", ")
}
