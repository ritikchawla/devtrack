# Dockerfile

# Build Stage
# Update Go version
# Use Go 1.24
FROM golang:1.24-alpine AS builder

# Install build tools and migrate CLI
RUN apk add --no-cache wget tar
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz -O migrate.tar.gz && \
    tar -xvf migrate.tar.gz && \
    # Corrected filename
    mv migrate /usr/local/bin/migrate && \
    rm migrate.tar.gz

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Build for linux amd64 which is the typical architecture for containers
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/server

# Final Stage
# Use specific version
FROM alpine:3.19

# Install runtime dependencies (wget for healthcheck, migrate needs nothing extra here)
RUN apk add --no-cache wget

WORKDIR /app

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built binary, migrate tool, and migration files from the builder stage with correct ownership
COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder --chown=appuser:appgroup /usr/local/bin/migrate /usr/local/bin/migrate
# Copy migration files
COPY --chown=appuser:appgroup db/migration ./db/migration/

# Expose port (use the ARG/ENV pattern for flexibility if needed later)
EXPOSE 8080

# Add a basic healthcheck (adjust command as needed once app has health endpoint)
# This assumes the app will listen on 8080 and respond to a simple TCP connection attempt
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Entrypoint script to run migrations then the app
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Switch to the non-root user BEFORE setting the entrypoint
USER appuser

ENTRYPOINT ["/app/entrypoint.sh"] # Use entrypoint script

# Default command for the entrypoint (the app itself)
CMD ["/app/main"]