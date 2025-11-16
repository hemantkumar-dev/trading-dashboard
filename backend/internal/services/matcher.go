package services

import (
	"context"
	"math/rand"
	"time"

	"trading-dashboard/backend/internal/websocket"

	"github.com/jmoiron/sqlx"
)

// OrderMatcher periodically scans open orders and attempts fills
type OrderMatcher struct {
	db    *sqlx.DB
	price *PriceEngine
	hub   *websocket.Hub
	stop  context.CancelFunc
}

func NewOrderMatcher(db *sqlx.DB, priceEngine *PriceEngine, hub *websocket.Hub) *OrderMatcher {
	return &OrderMatcher{db: db, price: priceEngine, hub: hub}
}

func (m *OrderMatcher) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.stop = cancel
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.tryMatch()
		}
	}
}

func (m *OrderMatcher) tryMatch() {
	// fetch open orders
	var orders []DBOrder
	if err := m.db.Select(&orders, "SELECT * FROM orders WHERE status = 'open'"); err != nil {
		return
	}

	for _, o := range orders {
		cur := m.price.GetPrice(o.Symbol)
		if cur == 0 {
			continue
		}
		// buyer wants price >= current? For buys, if limit >= current -> fill.
		match := false
		if o.Side == "buy" && o.Price >= cur {
			match = true
		}
		if o.Side == "sell" && o.Price <= cur {
			match = true
		}
		if !match {
			continue
		}

		// simulate fill fraction (30-100%)
		fraction := 0.3 + rand.Float64()*0.7
		quantityToFill := int(float64(o.Remaining) * fraction)
		if quantityToFill <= 0 {
			quantityToFill = 1
		}
		if quantityToFill > o.Remaining {
			quantityToFill = o.Remaining
		}

		// update order remaining
		newRemaining := o.Remaining - quantityToFill
		newStatus := "open"
		if newRemaining <= 0 {
			newStatus = "filled"
			newRemaining = 0
		}
		_, err := m.db.Exec("UPDATE orders SET remaining = ?, status = ? WHERE id = ?", newRemaining, newStatus, o.ID)
		if err != nil {
			continue
		}

		// record fill
		_ = recordFill(m.db, o.ID, o.UserID, o.Symbol, o.Side, quantityToFill, cur)

		// notify via WebSocket about fill
		m.hub.BroadcastFill(websocket.FillEvent{
			OrderID:  o.ID,
			Symbol:   o.Symbol,
			Side:     o.Side,
			Quantity: quantityToFill,
			Price:    cur,
		})
	}
}
