// internal/api/server.go
package api

import (
	"net/http"

	db "devtrack/db/sqlc" // Import generated db package
	"devtrack/internal/config"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our service.
type Server struct {
	config config.Config
	store  *db.Queries // Use the generated Queries struct
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(cfg config.Config, store *db.Queries) *Server {
	server := &Server{
		config: cfg,
		store:  store,
	}
	router := gin.Default() // Use default Gin router with logger and recovery middleware

	// Add routes here
	router.GET("/health", server.healthCheck) // Simple health check endpoint

	// User routes
	router.POST("/users", server.createUser)

	// TODO: Add routes for auth, repositories, tasks, time logs etc.
	// e.g., router.POST("/login", server.loginUser)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// healthCheck handler for the /health endpoint
func (server *Server) healthCheck(ctx *gin.Context) {
	// In a real app, you might check DB connection or other dependencies
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// TODO: Implement handler functions for other routes (e.g., loginUser)
