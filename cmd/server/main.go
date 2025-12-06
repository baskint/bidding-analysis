// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/baskint/bidding-analysis/api/trpc"
	"github.com/baskint/bidding-analysis/internal/config"
	"github.com/baskint/bidding-analysis/internal/ml"
	"github.com/baskint/bidding-analysis/internal/store"
)

func main() {
	log.Printf("🚀 Starting bidding-analysis server...")
	log.Printf("Environment: K_SERVICE=%s, PORT=%s", os.Getenv("K_SERVICE"), os.Getenv("PORT"))

	// Get port first (this should always work)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("✅ Port configured: %s", port)

	// Step 1: Try to load config
	log.Printf("📋 Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Printf("❌ Config loading failed: %v", err)
		log.Printf("🔄 Starting with minimal server...")
		startMinimalServer(port)
		return
	}
	log.Printf("✅ Configuration loaded successfully")

	// Step 2: Try database connection
	log.Printf("🗄️ Connecting to database...")
	log.Printf("Database URL exists: %t", cfg.DatabaseURL() != "")

	db, err := store.NewPostgresDB(cfg.DatabaseURL())
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		log.Printf("🔄 Starting with minimal server...")
		startMinimalServer(port)
		return
	}
	defer db.Close()
	log.Printf("✅ Database connected successfully")

	// Step 3: Initialize stores
	log.Printf("🏪 Initializing stores...")
	bidStore := store.NewBidStore(db)
	campaignStore := store.NewCampaignStore(db)
	log.Printf("✅ Stores initialized")

	// Step 4: Initialize ML predictor
	log.Printf("🤖 Initializing ML predictor...")
	
	var predictor *ml.Predictor
	
	// Get ML service URL from environment (set by Cloud Run)
	mlServiceURL := os.Getenv("ML_SERVICE_URL")
	
	if mlServiceURL != "" {
		// Use ML service
		log.Printf("Using ML service at: %s", mlServiceURL)
		mlPredictor, err := ml.NewMLPredictor(
			mlServiceURL,
			"", // encoders not needed for HTTP
			bidStore,
		)
		
		if err != nil {
			log.Printf("⚠️  ML service connection failed: %v", err)
			log.Printf("🔄 Falling back to OpenAI predictor...")
			predictor = ml.NewPredictor(cfg.OpenAI.APIKey, bidStore)
		} else {
			predictor = mlPredictor
			log.Printf("✅ ML service connected successfully!")
		}
	} else {
		// No ML service URL, use OpenAI
		log.Printf("No ML_SERVICE_URL set, using OpenAI predictor")
		if cfg.OpenAI.APIKey == "" {
			log.Printf("⚠️  No OpenAI API key configured")
		}
		predictor = ml.NewPredictor(cfg.OpenAI.APIKey, bidStore)
		log.Printf("✅ OpenAI predictor initialized")
	}

	// Step 5: Initialize tRPC handler
	log.Printf("🌐 Initializing tRPC handler...")
	trpcHandler := trpc.NewHandler(bidStore, campaignStore, predictor)
	log.Printf("✅ tRPC handler initialized")

	// Step 6: Setup routes and CORS
	log.Printf("🛣️ Setting up routes...")
	handler := corsMiddleware(trpcHandler.SetupRoutes())
	log.Printf("✅ Routes configured")

	// Step 7: Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("🎯 Starting tRPC server on port %s", port)
	log.Printf("✅ Server ready to accept connections")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// startMinimalServer starts a basic server when full initialization fails
func startMinimalServer(port string) {
	log.Printf("🔧 Starting minimal server on port %s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bidding Analysis API - Minimal Mode (Port: %s)", port)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Minimal Mode"))
	})

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"minimal_mode","port":"%s","message":"Full initialization failed"}`, port)
	})

	log.Printf("✅ Minimal server ready on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// corsMiddleware handles CORS for all requests
func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"https://bidding-analysis.web.app":         true,
		"https://bidding-analysis.firebaseapp.com": true,
		"http://localhost:3000":                    true,
		"http://localhost:3006":                    true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
