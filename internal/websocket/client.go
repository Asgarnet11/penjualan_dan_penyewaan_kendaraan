package websocket

import (
	"log"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type Client struct {
	ID   uuid.UUID
	Conn *ws.Conn
	Hub  *Hub
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			// Hub menutup channel.
			c.Conn.WriteMessage(ws.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteMessage(ws.TextMessage, message)
	}
}
