package main

import (
	"github.com/SwanHtetAungPhyo/binance-dash/internal/config"
	"log"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err.Error())
		return
	}
}
