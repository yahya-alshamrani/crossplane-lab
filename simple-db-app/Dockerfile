# --- Stage 1: Build Stage ---
# Using 1.25-alpine to satisfy the 'go >= 1.25.2' requirement in your go.mod
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy dependency files first. The * ensures it doesn't fail if go.sum is missing
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the application source code and the templates directory
COPY . .

# Build a statically linked binary
# CGO_ENABLED=0 is required for the 'static' distroless image
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app main.go

# --- Stage 2: Final Runtime Stage ---
# Distroless static is the smallest and most secure image for Go
FROM gcr.io/distroless/static-debian12

# Copy the binary from the builder
COPY --from=builder /go-app /go-app

# IMPORTANT: Copy the templates folder so the app can find your index.html
COPY --from=builder /app/templates /templates

# Expose the application port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/go-app"]
