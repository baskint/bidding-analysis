package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/baskint/bidding-analysis/api/grpc/services"
	"github.com/baskint/bidding-analysis/api/trpc"
	"github.com/baskint/bidding-analysis/internal/config"
	"github.com/baskint/bidding-analysis/internal/ml"
	"github.com/baskint/bidding-analysis/internal/store"
)

// corsMiddleware handles CORS for all requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Define allowed origins
		allowedOrigins := []string{
			"http://localhost:3000", // Next.js default
			"http://localhost:3006", // Your current frontend
		}

		// Add production origins from environment variable
		if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
			prodOrigins := strings.Split(envOrigins, ",")
			for _, prodOrigin := range prodOrigins {
				allowedOrigins = append(allowedOrigins, strings.TrimSpace(prodOrigin))
			}
		}

		// Check if origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := store.NewPostgresDB(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize stores
	bidStore := store.NewBidStore(db)
	campaignStore := store.NewCampaignStore(db)

	// Initialize ML predictor
	predictor := ml.NewPredictor(cfg.OpenAI.APIKey, bidStore)

	// Check if running in Cloud Run (simpler deployment)
	if os.Getenv("K_SERVICE") != "" {
		// Running in Cloud Run - only start tRPC server
		log.Printf("Running in Cloud Run")
		if err := startTRPCServer(cfg, bidStore, campaignStore, predictor); err != nil {
			log.Fatalf("Failed to start tRPC server: %v", err)
		}
		return
	}

	// Local development - start both servers
	// Initialize services
	biddingService := services.NewBiddingService(bidStore, cfg)
	analyticsService := services.NewAnalyticsService(campaignStore, bidStore)

	// Start gRPC server in a goroutine
	go func() {
		if err := startGRPCServer(cfg, biddingService, analyticsService); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start tRPC server in a goroutine
	go func() {
		if err := startTRPCServer(cfg, bidStore, campaignStore, predictor); err != nil {
			log.Fatalf("Failed to start tRPC server: %v", err)
		}
	}()

	log.Printf("Server started successfully")
	log.Printf("gRPC server listening on port %d", cfg.Server.GRPCPort)
	log.Printf("tRPC server listening on port %d", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
}

// startGRPCServer starts the gRPC server
func startGRPCServer(cfg *config.Config, biddingService *services.BiddingService, analyticsService *services.AnalyticsService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()

	// Register services
	// pb.RegisterBiddingServiceServer(s, biddingService)
	// pb.RegisterAnalyticsServiceServer(s, analyticsService)

	// Enable reflection for development
	reflection.Register(s)

	log.Printf("gRPC server starting on port %d", cfg.Server.GRPCPort)
	return s.Serve(lis)
}

// startTRPCServer starts the tRPC server with CORS support
func startTRPCServer(cfg *config.Config, bidStore *store.BidStore, campaignStore *store.CampaignStore, predictor *ml.Predictor) error {
	// Initialize tRPC handler
	trpcHandler := trpc.NewHandler(bidStore, campaignStore, predictor)

	// Use Cloud Run's PORT environment variable, fallback to config
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", cfg.Server.Port)
	}

	// Setup HTTP server with tRPC routes and CORS middleware
	handler := corsMiddleware(trpcHandler.SetupRoutes())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	log.Printf("tRPC server starting on port %s with CORS enabled", port)

	// Log allowed origins for debugging
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:3006"}
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		prodOrigins := strings.Split(envOrigins, ",")
		allowedOrigins = append(allowedOrigins, prodOrigins...)
	}
	log.Printf("CORS allowed origins: %v", allowedOrigins)

	return server.ListenAndServe()
}
