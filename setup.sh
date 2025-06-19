#!/bin/bash

# Create directory structure
mkdir -p cmd/server cmd/migrate
mkdir -p api/proto/generated api/grpc api/rest
mkdir -p internal/{config,store,ml,processor,models}
mkdir -p migrations scripts docker deployments

# Initialize Go module
go mod init github.com/yourusername/bidding-analysis

# Install core dependencies
echo "Installing Go dependencies..."
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get github.com/lib/pq@latest
go get github.com/gin-gonic/gin@latest
go get github.com/golang-migrate/migrate/v4@latest
go get github.com/joho/godotenv@latest
go get github.com/sashabaranov/go-openai@latest

# Create basic Makefile
cat > Makefile << 'EOF'
.PHONY: build run test proto migrate-up migrate-down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto

include .env
export

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down


# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf api/proto/generated/*.pb.go
EOF

# Create initial README
cat > README.md << 'EOF'
# Bidding Analysis System

AI-powered digital advertising bid optimization platform.

## Features
- Real-time bid processing
- ML-based price prediction
- Fraud detection
- Campaign analytics

## Setup
1. Copy `.env.example` to `.env` and configure
2. Run `make proto` to generate protobuf files
3. Run `make migrate-up` to setup database
4. Run `make run` to start server

## Development
- `make build` - Build application
- `make test` - Run tests
- `make proto` - Generate protobuf files
EOF

echo "Project structure created successfully!"
echo "Next steps:"
echo "1. Update go.mod with your GitHub username"
echo "2. Run the commands to install dependencies"
echo "3. Start implementing core files"