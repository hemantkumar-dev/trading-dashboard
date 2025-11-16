package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) WritePump() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("ws write err:", err)
			return
		}
	}
}
