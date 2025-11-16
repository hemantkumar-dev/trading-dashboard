package api

import (
	"net/http"

	"trading-dashboard/backend/internal/models"
	"trading-dashboard/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type placeOrderReq struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func RegisterOrderRoutes(r *gin.Engine, db *sqlx.DB, jwtSecret string) {
	auth := JWTMiddleware(jwtSecret)
	r.POST("/orders", auth, func(c *gin.Context) {
		var req placeOrderReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		username := c.GetString("username")
		order := models.Order{
			Symbol:   req.Symbol,
			Side:     req.Side,
			Quantity: req.Quantity,
			Price:    req.Price,
		}
		id, err := services.PlaceOrder(db, username, order)
		if err != nil {
			c.JSON(500, gin.H{"error": "db"})
			return
		}
		c.JSON(200, gin.H{"order_id": id})
	})

	// list orders (admin/demo, no auth required here but you can secure)
	r.GET("/orders", func(c *gin.Context) {
		orders, err := services.GetAllOrders(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "db"})
			return
		}
		c.JSON(200, orders)
	})

	// get order by id
	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		var o services.DBOrder
		if err := db.Get(&o, "SELECT * FROM orders WHERE id = ?", id); err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.JSON(200, o)
	})

	// quick debug: cancel order
	r.POST("/orders/:id/cancel", auth, func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("UPDATE orders SET status='cancelled' WHERE id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db"})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})
}
