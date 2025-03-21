# Use official Go image as the base
FROM golang:1.24-alpine

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum, download dependencies
COPY go.mod go.sum .secrets ./
RUN go mod download

# Copy the rest of the source code
COPY src/ ./src/

# Build the Go app
RUN go build -o task-manager ./src/cmd/server/main.go

# Expose port 8080 (our API port)
EXPOSE 8080

# Run the app when the container starts
CMD ["./task-manager"]