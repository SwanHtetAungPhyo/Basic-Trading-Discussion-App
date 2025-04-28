package models

// TickerAndIndicator contains market data and technical indicators
type TickerAndIndicator struct {
	Ticker string `json:"ticker"`
	SMA    string `json:"sma"`
	RSI    string `json:"rsi"`
	EMA    string `json:"ema"`
	Time   string `json:"time"`
	Price  string `json:"price"`
}
