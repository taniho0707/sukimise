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

	"sukimise-discord-bot/internal/config"
	"sukimise-discord-bot/internal/handlers"
	"sukimise-discord-bot/internal/services"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from parent directory (for local development)
	// In Docker, environment variables are already set via docker-compose.yml
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Info: .env file not found (normal in Docker environment)")
	} else {
		log.Println("Successfully loaded .env file")
	}
	
	// Log Google Maps API key status (first few characters for security)
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey != "" {
		keyPrefix := apiKey
		if len(apiKey) > 10 {
			keyPrefix = apiKey[:10]
		}
		log.Printf("Google Maps API Key loaded: %s...", keyPrefix)
	} else {
		log.Println("Google Maps API Key: NOT SET")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize services
	discordService := services.NewDiscordService(db, cfg.SukimiseAPIURL)
	
	// Create Discord session
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	// Initialize handlers
	commandHandler := handlers.NewCommandHandler(discordService)

	// Register command handlers
	dg.AddHandler(commandHandler.HandleSlashCommand)
	dg.AddHandler(commandHandler.HandleReady)

	// Set intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// Open WebSocket connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Failed to open Discord connection: %v", err)
	}
	defer dg.Close()

	// Start health check HTTP server
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"error","message":"Database connection failed","error":"%s"}`, err.Error())
			return
		}

		// Check Discord connection
		discordStatus := "connected"
		if dg.State == nil {
			discordStatus = "disconnected"
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","message":"Discord Bot is running","database":"ok","discord":"%s","timestamp":%d}`, 
			discordStatus, time.Now().Unix())
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.BotPort),
		Handler: httpMux,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Printf("Health check server starting on port %s", cfg.BotPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health check server error: %v", err)
		}
	}()

	log.Println("Sukimise Discord Bot is running. Press Ctrl+C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down Sukimise Discord Bot...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Health check server shutdown error: %v", err)
	}
}