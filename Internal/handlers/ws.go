package handlers

import (
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
    Conn     *websocket.Conn
    StudioID string
    UserID   string
}

var (
    clients      = make(map[string][]*Client) // studioID -> clients
    clientsMutex sync.Mutex
)

func WebSocketHandler(c *gin.Context) {
    studioID := c.Query("studioId")
    userID := c.Query("userId") // You should validate this!

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    client := &Client{Conn: conn, StudioID: studioID, UserID: userID}

    clientsMutex.Lock()
    clients[studioID] = append(clients[studioID], client)
    clientsMutex.Unlock()

    defer func() {
        conn.Close()
        clientsMutex.Lock()
        // Remove client from list
        for i, cl := range clients[studioID] {
            if cl == client {
                clients[studioID] = append(clients[studioID][:i], clients[studioID][i+1:]...)
                break
            }
        }
        clientsMutex.Unlock()
    }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        // Broadcast to all clients in the same studio except sender
        clientsMutex.Lock()
        for _, cl := range clients[studioID] {
            if cl != client {
                cl.Conn.WriteMessage(websocket.TextMessage, msg)
            }
        }
        clientsMutex.Unlock()
    }
}