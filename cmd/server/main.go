package main

import (
	"context"
	"log"

	db "devtrack/db/sqlc"      // Import generated db package
	"devtrack/internal/api"    // Import api package
	"devtrack/internal/config" // Import config package

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	// Assuming the executable runs from the root directory or Dockerfile WORKDIR
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	log.Println("Configuration loaded successfully")
	log.Printf("Attempting to connect to database: %s", cfg.DBSource) // Be careful logging sensitive info

	// Establish database connection pool
	connPool, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer connPool.Close() // Ensure pool is closed when main exits

	// Ping the database to verify connection
	err = connPool.Ping(context.Background())
	if err != nil {
		log.Fatal("cannot ping db:", err)
	}

	log.Println("Database connection established successfully")

	// Create a store (querier) instance using the connection pool
	store := db.New(connPool) // Use db.New which returns *Queries
	// --- Server Setup and Start Logic ---
	server := api.NewServer(cfg, store) // Create a new server instance

	log.Printf("Starting devtrack server on 0.0.0.0:%s...", cfg.APIPort)
	err = server.Start("0.0.0.0:" + cfg.APIPort) // Start the server
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
