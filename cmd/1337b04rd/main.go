package main

import (
	"1337b04rd/internal/adapters/fileUtils"
	"1337b04rd/internal/adapters/handlers"
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/adapters/rickMorty"
	"1337b04rd/internal/adapters/triples"
	"1337b04rd/internal/services"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	file_utils := fileUtils.NewFileUtils()
	imageStorage := triples.NewTriples("1337b04rd", 1414)
	userOutlook := rickMorty.NewRickMortyAPI()

	userRepo := postgres.NewUserRepository(db)
	postRepo := postgres.NewPostRepository(db, "1337b04rd")

	userService := services.NewUserService(userRepo, userOutlook)
	postServices := services.NewPostService(postRepo, imageStorage, file_utils, *userService, "1337b04rd")

	router := handlers.NewRouter(*userService, *postServices)

	port := ":8080"

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func initDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to database")
	return db, nil
}
