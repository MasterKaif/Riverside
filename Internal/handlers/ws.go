package handlers

import (
	"encoding/json"
	"log"
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

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	UserID  string          `json:"userId,omitempty"`
}

func WebSocketHandler(c *gin.Context) {
    log.Println(("WebSocketHandler called"))
	studioID := c.Query("studioId")
	userID := c.Query("userId") // You should validate this!

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &Client{Conn: conn, StudioID: studioID, UserID: userID}

	clientsMutex.Lock()
	isFirst := len(clients[studioID]) == 0
	clients[studioID] = append(clients[studioID], client)
	clientsMutex.Unlock()

	// Notify the client that they joined
	joinMsg := WSMessage{Type: "join", UserID: userID}
	msgBytes, _ := json.Marshal(joinMsg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)

	if isFirst {
		// Optionally notify that the room is initialized
		initMsg := WSMessage{Type: "room-initialized", UserID: userID}
		msgBytes, _ := json.Marshal(initMsg)
		conn.WriteMessage(websocket.TextMessage, msgBytes)
	} else {
		// Notify others that a new user joined
		clientsMutex.Lock()
		for _, cl := range clients[studioID] {
			if cl != client {
				newJoinMsg := WSMessage{Type: "user-joined", UserID: userID}
				msgBytes, _ := json.Marshal(newJoinMsg)
				cl.Conn.WriteMessage(websocket.TextMessage, msgBytes)
			}
		}
		clientsMutex.Unlock()
	}

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
		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			continue
		}
		switch wsMsg.Type {
		case "offer", "answer", "ice-candidate":
            log.Println("Received signaling message:", wsMsg.Type)
            log.Println("Payload:", string(wsMsg.Payload))
			// Relay signaling messages to all other clients in the room
			clientsMutex.Lock()
			for _, cl := range clients[studioID] {
				if cl != client {
					cl.Conn.WriteMessage(websocket.TextMessage, msg)
				}
			}
			clientsMutex.Unlock()
		case "join":

			// Already handled on connect, but you could add logic here if needed
		}
	}
}
