# Step 1: Build the Go binary
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Download dependencies (use cache if they haven't changed)
RUN go mod download

# Copy the entire project into the working directory
COPY ./ ./

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Step 2: Create a lightweight image with the binary
FROM scratch

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy config yaml from builder stage
COPY --from=builder /app/config.yaml .

# Expose the port that the application listens on (replace 8080 if different)
EXPOSE 3000

# Set the entry point for the container
CMD ["./main"]
