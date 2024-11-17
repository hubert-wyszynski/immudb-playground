# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy only go.mod (since go.sum doesn't exist yet)
COPY go.mod ./

# Download dependencies (this will create go.sum if needed)
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"] 