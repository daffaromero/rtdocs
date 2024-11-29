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
	"rtdocs/utils"

	"github.com/gorilla/mux"
)

var (
	secretKey       = utils.GetEnv("SECRET_KEY")
	accessDuration  = utils.GetEnv("ACCESS_TOKEN_DURATION")
	refreshDuration = utils.GetEnv("REFRESH_TOKEN_DURATION")
)

func main() {
	// Connect to the database
	dbConfig := config.NewPostgresDatabase()

	// Token generator

	// Set up dependencies
	docsRepo := repository.NewDocumentRepository(dbConfig)
	userRepo := repository.NewUserRepository(dbConfig)

	docsService := service.NewDocumentService(docsRepo)
	authService := service.NewAuthService(userRepo, utils.NewTokenGenerator(secretKey, accessDuration, refreshDuration))
	userService := service.NewUserService(userRepo)

	docsController := controller.NewDocumentController(docsService)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	wsController := controller.NewWebSocketController(docsService)

	ctx := context.Background()

	// Start the WebSocket message handler in a goroutine
	go wsController.HandleMessages(ctx)

	// Create a new router
	router := mux.NewRouter()

	// Set up HTTP handler for WebSocket connections
	router.HandleFunc("/ws", wsController.HandleConnections)

	// Set up HTTP handlers for document operations
	router.HandleFunc("/api/documents", docsController.GetAllDocuments)
	router.HandleFunc("/api/document/create", docsController.CreateDocument)
	router.HandleFunc("/api/document/save", docsController.UpdateDocument)
	router.HandleFunc("/api/document/{id}", docsController.GetDocument)

	// Set up HTTP handlers for authentication operations
	router.HandleFunc("/api/auth/register", authController.Register)
	router.HandleFunc("/api/auth/login", authController.Login)
	router.HandleFunc("/api/auth/logout", authController.Logout)
	router.HandleFunc("/api/auth/guest", authController.Guest)

	// Set up HTTP handlers for user operations
	router.HandleFunc("/api/user/{id}", userController.GetUser)
	router.HandleFunc("/api/users", userController.GetAllUsers)
	router.HandleFunc("/api/user/create", userController.CreateUser)
	router.HandleFunc("/api/user/update", userController.UpdateUser)

	// Wrap the HTTP handler with the middlewares
	corsHandler := middleware.CORSMiddleware(router)

	// Create a subrouter for the routes that require authentication
	authRouter := router.PathPrefix("/api").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	// Set up HTTP handlers for the routes that require authentication
	authRouter.HandleFunc("/documents", docsController.GetAllDocuments).Methods("GET")
	authRouter.HandleFunc("/document/{id}", docsController.GetDocument).Methods("GET")
	authRouter.HandleFunc("/document/create", docsController.CreateDocument).Methods("POST")
	authRouter.HandleFunc("/document/save", docsController.UpdateDocument).Methods("PUT")
	authRouter.HandleFunc("/auth/logout", authController.Logout).Methods("POST")
	authRouter.HandleFunc("/user/{id}", userController.GetUser).Methods("GET")
	authRouter.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	authRouter.HandleFunc("/user/create", userController.CreateUser).Methods("POST")
	authRouter.HandleFunc("/user/update", userController.UpdateUser).Methods("PUT")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe("localhost:8080", corsHandler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
