package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shihab369/Tapx-lite/internal/handler"
	"github.com/Shihab369/Tapx-lite/internal/repository"
	"github.com/Shihab369/Tapx-lite/internal/service"

	_ "github.com/lib/pq"
)

// Config holds all environment-driven configuration for the application.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
}

// loadConfig reads configuration from environment variables and validates required fields.
// Returns an error if any required variable is missing.
func loadConfig() (Config, error) {
	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	missing := []string{}
	if cfg.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	// Apply defaults for optional fields.
	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable" // safe default for local dev; use "require" in production
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg, nil
}

// connectDB opens a PostgreSQL connection and configures the connection pool.
// Verifies connectivity with a context-bound ping before returning.
func connectDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool tuning.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to reach database: %w", err)
	}

	return db, nil
}

// loggingMiddleware logs method, path, duration, and remote IP for every request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s ip=%s",
			r.Method, r.URL.Path, time.Since(start), r.RemoteAddr,
		)
	})
}

// jsonMiddleware sets the Content-Type header to application/json for all responses.
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// chainMiddleware applies a list of middleware to a handler in declaration order.
func chainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func main() {
	// Load and validate configuration from environment.
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	log.Println("configuration loaded")

	// Establish database connection with pool.
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer db.Close()
	log.Println("database connected")

	// Wire up dependency layers.
	repo := repository.NewTransactionRepository(db)
	svc := service.NewTransactionService(repo)
	h := handler.NewTransactionHandler(svc)

	// Register routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/sync", h.Sync)
	mux.HandleFunc("/balance", h.GetBalance)

	// Apply middleware chain.
	chain := chainMiddleware(mux,
		loggingMiddleware,
		jsonMiddleware,
	)

	// Configure HTTP server with production timeouts.
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      chain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine to allow graceful shutdown handling.
	go func() {
		log.Printf("tapx-lite listening on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until OS termination signal is received (SIGINT or SIGTERM).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("signal received: %v — initiating graceful shutdown", sig)

	// Allow up to 10 seconds for in-flight requests to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("tapx-lite stopped cleanly")
}