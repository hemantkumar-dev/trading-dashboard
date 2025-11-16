package services

import (
	"math/rand"
	"time"

	"trading-dashboard/backend/internal/models"
	"trading-dashboard/backend/internal/websocket"
)

type PriceEngine struct {
	Prices map[string]float64
	Hub    *websocket.Hub
}

func NewPriceEngine(hub *websocket.Hub) *PriceEngine {
	return &PriceEngine{
		Hub: hub,
		Prices: map[string]float64{
			"AAPL": 180.0,
			"TSLA": 250.0,
			"AMZN": 140.0,
			"INFY": 1550.0,
			"TCS":  4120.0,
		},
	}
}

func (p *PriceEngine) Start() {
	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(2 * time.Second)

		for symbol, price := range p.Prices {
			// +/- 0.5% - 2%
			pct := 0.005 + rand.Float64()*0.015
			if rand.Intn(2) == 0 {
				pct = -pct
			}
			newPrice := price * (1 + pct)
			p.Prices[symbol] = newPrice

			msg := models.Price{Symbol: symbol, Price: newPrice}
			p.Hub.BroadcastPrice(msg)
		}
	}
}

func (p *PriceEngine) GetPrice(symbol string) float64 {
	if v, ok := p.Prices[symbol]; ok {
		return v
	}
	return 0
}
