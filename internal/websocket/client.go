package websocket

import (
	"log"

	"sultra-otomotif-api/internal/model"

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
		var msg model.Message

		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading json: %v", err)
			break
		}

		msg.SenderID = c.ID

		c.Hub.Broadcast <- &msg
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
