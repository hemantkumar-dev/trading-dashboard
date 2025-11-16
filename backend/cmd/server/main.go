package main

import (
	"log"
	"os"

	"trading-dashboard/backend/internal/api"
	"trading-dashboard/backend/internal/services"
	"trading-dashboard/backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// load env
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// prefer DATABASE_URL (commonly provided by managed Postgres services)
		dbPath = os.Getenv("DATABASE_URL")
	}
	if dbPath == "" {
		dbPath = "./data/trading.db"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	// DB init
	db, err := services.NewDB(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	if err := services.MigrateDB(db); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	// WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Price Engine
	priceEngine := services.NewPriceEngine(hub)
	go priceEngine.Start()

	// Order matcher background
	orderMatcher := services.NewOrderMatcher(db, priceEngine, hub)
	go orderMatcher.Start()

	// API routes
	api.RegisterAuthRoutes(r, db, jwtSecret)
	api.RegisterPriceRoutes(r, priceEngine)
	api.RegisterWebSocketRoutes(r, hub)
	api.RegisterOrderRoutes(r, db, jwtSecret)
	api.RegisterHoldingsRoutes(r, db, jwtSecret)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
