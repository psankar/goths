# Build stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the generated files and source code
COPY . /app

# Build the application
RUN CGO_ENABLED=0 go build -o goths

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/goths .

# Copy migration files
COPY sqlc/schema/ /app/schema/

# Run the application
CMD ["./goths"]
