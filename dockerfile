# Use Go 1.24 Alpine as the base image
FROM golang:1.24-alpine

ENV APP_ENV=cont
# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire src/ directory
COPY src/ ./src/

# Copy .secrets/.env.local
COPY .secrets/.env.local .secrets/.env.local

COPY .secrets/.env.cont .secrets/.env.cont

# Install protoc for proto generation
RUN apk add --no-cache protobuf

# Build the Go application
RUN go build -o task-manager ./src/cmd/server/main.go

# Expose both REST and gRPC ports
EXPOSE 8080 50051

# Command to run the application
CMD ["./task-manager"]