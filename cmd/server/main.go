package main

import (
	"github.com/SwanHtetAungPhyo/binance-dash/internal/app"
	"log"
)

func main() {
	if err := app.NewApp().Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
