# ------------ Stage 1: Build ------------
# Use an official Go image with Alpine Linux. Choose a specific Go version.
FROM golang:1.22-alpine AS builder

# Install build tools and UPX for compression
# git might be needed if you have private go modules or specific dependencies
RUN apk add --no-cache git upx

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
RUN go mod verify

# Copy the rest of the source code
COPY bot.go ./

# Build the Go application statically
# CGO_ENABLED=0 is crucial for static linking, especially for minimal base images.
# -ldflags="-w -s" strips debug symbols and symbol table, reducing binary size.
# Replace 'mybot' if your desired output binary name is different.
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/telegram-userinfo-bot .

# Compress the binary using UPX for maximum size reduction
# --best uses more time for potentially better compression
# --lzma often provides good results
RUN upx --best --lzma /app/telegram-userinfo-bot

# ------------ Stage 2: Runtime ------------
# Use a minimal base image like Alpine.
# Alpine includes CA certificates needed for HTTPS requests by the bot library.
# Using 'scratch' is smaller but requires manually copying CA certs if HTTPS is needed.
FROM alpine:latest

# Create a non-privileged user and group for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy only the compressed, static binary from the builder stage
COPY --from=builder /app/telegram-userinfo-bot /app/telegram-userinfo-bot

# Alpine includes CA certificates at /etc/ssl/certs/ca-certificates.crt
# If you were using 'scratch' as a base, you would need to copy them:
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Ensure the binary is executable by the user
RUN chown appuser:appgroup /app/telegram-userinfo-bot && chmod +x /app/telegram-userinfo-bot

# Switch to the non-root user
USER appuser

# Define the entry point for the container. This command runs when the container starts.
ENTRYPOINT ["/app/telegram-userinfo-bot"]