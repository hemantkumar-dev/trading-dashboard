package api

import (
	"net/http"

	wsinternal "trading-dashboard/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterWebSocketRoutes(r *gin.Engine, hub *wsinternal.Hub) {
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		client := &wsinternal.Client{
			Conn: conn,
			Send: make(chan []byte, 256),
		}
		hub.Register <- client

		// writer
		go client.WritePump()

		// reader - just listens for close
		go func() {
			defer func() {
				hub.Unregister <- client
				_ = conn.Close()
			}()
			for {
				if _, _, err := conn.NextReader(); err != nil {
					break
				}
			}
		}()
	})
}
