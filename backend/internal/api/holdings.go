package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HoldingRow struct {
	Symbol   string  `json:"symbol"`
	NetQty   int     `json:"net_qty"`
	AvgCost  float64 `json:"avg_cost"`
	Realized float64 `json:"realized,omitempty"`
}

func RegisterHoldingsRoutes(r *gin.Engine, db *sqlx.DB, jwtSecret string) {
	auth := JWTMiddleware(jwtSecret)
	r.GET("/holdings", auth, func(c *gin.Context) {
		username := c.GetString("username")
		// derive holdings from fills for this user
		var userID int64
		if err := db.Get(&userID, "SELECT id FROM users WHERE username = ?", username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user"})
			return
		}

		// aggregate buys and sells
		rows, err := db.Queryx(`
SELECT symbol,
 SUM(CASE WHEN side='buy' THEN quantity ELSE 0 END) as buys_qty,
 SUM(CASE WHEN side='buy' THEN quantity*price ELSE 0 END) as buys_value,
 SUM(CASE WHEN side='sell' THEN quantity ELSE 0 END) as sells_qty
FROM fills
WHERE user_id = ?
GROUP BY symbol
`, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "db"})
			return
		}
		defer rows.Close()
		var out []HoldingRow
		for rows.Next() {
			var sym string
			var buysQty int
			var buysValue float64
			var sellsQty int
			if err := rows.Scan(&sym, &buysQty, &buysValue, &sellsQty); err != nil {
				continue
			}
			net := buysQty - sellsQty
			avg := 0.0
			if buysQty > 0 {
				avg = buysValue / float64(buysQty)
			}
			out = append(out, HoldingRow{
				Symbol:  sym,
				NetQty:  net,
				AvgCost: avg,
			})
		}
		c.JSON(200, out)
	})
}
