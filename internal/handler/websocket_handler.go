package handler

import (
	"log"
	"net/http"
	"sultra-otomotif-api/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *websocket.Hub, ctx *gin.Context) {
	userID, exists := ctx.Get("currentUserID")
	if !exists {
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &websocket.Client{
		ID:   userID.(uuid.UUID),
		Conn: conn,
		Hub:  hub,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
