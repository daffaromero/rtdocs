package main

import (
	"context"
	"log"
	"net/http"
	"rtdocs/config"
	"rtdocs/controller"
	"rtdocs/middleware"
	"rtdocs/repository"
	"rtdocs/service"
)

func main() {
	// Connect to the database
	dbConfig := config.NewPostgresDatabase()

	// Set up dependencies
	docsRepo := repository.NewDocumentRepository(dbConfig)
	userRepo := repository.NewUserRepository(dbConfig)

	docsService := service.NewDocumentService(docsRepo)
	userService := service.NewUserService(userRepo)

	docsController := controller.NewDocumentController(docsService)
	userController := controller.NewUserController(userService)
	wsController := controller.NewWebSocketController(docsService)

	ctx := context.Background()

	// Start the WebSocket message handler in a goroutine
	go wsController.HandleMessages(ctx)

	// Set up HTTP handler for WebSocket connections
	http.HandleFunc("/ws", wsController.HandleConnections)

	// Set up HTTP handlers for document operations
	http.HandleFunc("/api/documents", docsController.GetAllDocuments)
	http.HandleFunc("/api/document/{id}", docsController.GetDocument)
	http.HandleFunc("/api/document/create", docsController.CreateDocument)
	http.HandleFunc("/api/document/save", docsController.UpdateDocument)

	// Set up HTTP handlers for user operations
	http.HandleFunc("/api/user/{id}", userController.GetUser)
	http.HandleFunc("/api/users", userController.GetAllUsers)
	http.HandleFunc("/api/user/create", userController.CreateUser)
	http.HandleFunc("/api/user/update", userController.UpdateUser)

	// Wrap the HTTP handler with the middlewares
	corsHandler := middleware.CORSMiddleware(http.DefaultServeMux)
	authHandler := middleware.AuthMiddleware(corsHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", authHandler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
