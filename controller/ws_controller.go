package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rtdocs/model/domain"
	"rtdocs/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketController interface {
	HandleConnections(ctx context.Context, w http.ResponseWriter, r *http.Request)
	HandleMessages()
}

type webSocketController struct {
	docService service.DocumentService
	clients    map[*websocket.Conn]bool
	broadcast  chan map[string]string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketController(docService service.DocumentService) *webSocketController {
	return &webSocketController{
		docService: docService,
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan map[string]string),
	}
}

// HandleConnections upgrades HTTP requests to WebSocket and registers clients
func (c *webSocketController) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer ws.Close()

	c.clients[ws] = true

	vars := mux.Vars(r)
	documentID := vars["id"]
	if documentID == "" {
		log.Printf("Document ID is required")
		return
	}

	// Create a new context with the request context
	ctx := r.Context()

	// Retrieve the current document state and send it to the new client
	document, err := c.docService.GetDocument(ctx, documentID)
	log.Printf("Document ID: %s", documentID)
	if err != nil {
		log.Printf("Failed to load document: %v", err)
		document = &domain.Document{
			ID:      documentID,
			Title:   "New Document",
			Content: "",
		}
	}

	initialState := map[string]string{
		"type":    "initial",
		"title":   document.Title,
		"content": document.Content,
	}

	ws.WriteJSON(initialState)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			delete(c.clients, ws)
			break
		}

		var update map[string]string
		if err := json.Unmarshal(msg, &update); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		if update["type"] == "title" {
			document.Title = update["title"]
		} else if update["type"] == "content" {
			document.Content = update["content"]
		}

		var updatedDoc domain.Document
		updatedDoc.ID = document.ID
		updatedDoc.Title = document.Title
		updatedDoc.Content = document.Content

		// Update the document and broadcast the new content
		if _, err = c.docService.UpdateDocument(ctx, &updatedDoc); err != nil {
			log.Printf("Failed to update document: %v", err)
		}

		c.broadcast <- update
	}
}

// HandleMessages listens for new messages on the broadcast channel and sends them to all clients
func (c *webSocketController) HandleMessages(ctx context.Context) {
	for {
		select {
		case update := <-c.broadcast:
			for client := range c.clients {
				if err := client.WriteJSON(update); err != nil {
					log.Printf("Write error: %v", err)
					client.Close()
					delete(c.clients, client)
				}
			}
		case <-ctx.Done():
			log.Println("Stopping message broadcasting due to context cancellation")
			return
		}
	}
}
