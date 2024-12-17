# Step 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /

# Copy the go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Download dependencies (use cache if they haven't changed)
RUN go mod download

# Copy the entire project into the working directory
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app /cmd/main.go

# Step 2: Create a lightweight image with the binary
FROM alpine

# Set the working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app .

# Expose the port that the application listens on (replace 8080 if different)
EXPOSE 50051 50052 50053

# Set the entry point for the container
CMD ["./app", "-config", "/config.yaml"]
