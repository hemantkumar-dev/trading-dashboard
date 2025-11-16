package api

import (
	"trading-dashboard/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterPriceRoutes(r *gin.Engine, engine *services.PriceEngine) {
	r.GET("/prices", func(c *gin.Context) {
		c.JSON(200, engine.Prices)
	})
}
