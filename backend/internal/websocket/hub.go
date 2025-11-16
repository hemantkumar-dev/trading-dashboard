package websocket

import (
	"encoding/json"
	"log"

	"trading-dashboard/backend/internal/models"
)

type FillEvent struct {
	OrderID  int64   `json:"order_id"`
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client

	// separate channels for price and fills
	priceChan chan models.Price
	fillChan  chan FillEvent
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		priceChan:  make(chan models.Price, 256),
		fillChan:   make(chan FillEvent, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c] = true
		case c := <-h.Unregister:
			delete(h.Clients, c)
		case p := <-h.priceChan:
			msg := map[string]interface{}{"type": "price", "symbol": p.Symbol, "price": p.Price}
			data, _ := json.Marshal(msg)
			for client := range h.Clients {
				select {
				case client.Send <- data:
				default:
				}
			}
		case f := <-h.fillChan:
			msg := map[string]interface{}{"type": "fill", "order_id": f.OrderID, "symbol": f.Symbol, "side": f.Side, "quantity": f.Quantity, "price": f.Price}
			data, _ := json.Marshal(msg)
			for client := range h.Clients {
				select {
				case client.Send <- data:
				default:
				}
			}
		}
	}
}

func (h *Hub) BroadcastPrice(p models.Price) {
	select {
	case h.priceChan <- p:
	default:
		log.Println("price chan full, dropping")
	}
}

func (h *Hub) BroadcastFill(f FillEvent) {
	select {
	case h.fillChan <- f:
	default:
		log.Println("fill chan full, dropping")
	}
}
