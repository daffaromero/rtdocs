package main

import (
	"context"
	"log"
	"net/http"
	"rtdocs/config"
	"rtdocs/controller"
	"rtdocs/repository"
	"rtdocs/repository/query"
	"rtdocs/service"
)

func main() {
	// Connect to the database
	dbConfig := config.NewPostgresDatabase()

	// Set up dependencies
	docsQuery := query.NewDocumentQuery(dbConfig)
	docsRepo := repository.NewDocumentRepository(docsQuery)
	docsService := service.NewDocumentService(docsRepo)
	docsController := controller.NewDocumentController(docsService)
	wsController := controller.NewWebSocketController(docsService)

	ctx := context.Background()

	// Start the WebSocket message handler in a goroutine
	go wsController.HandleMessages(ctx)

	// Set up HTTP handler for WebSocket connections
	http.HandleFunc("/ws", wsController.HandleConnections)

	// Set up HTTP handlers for document operations
	http.HandleFunc("/api/document", docsController.GetDocument)
	http.HandleFunc("/api/documents", docsController.GetAllDocuments)
	http.HandleFunc("/api/document/save", docsController.UpdateDocumentContent)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}