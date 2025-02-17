# Build stage
FROM golang:1.21-alpine AS builder

# Install git and SSL certificates for private repos (if needed)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cert-exporter

# Final stage
FROM alpine:3.18

# Add CA certificates for SSL verification
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -H -h /app certexporter
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cert-exporter .

# Copy example config
COPY example.yml /etc/cert-exporter/config.yml

# Set ownership
RUN chown -R certexporter:certexporter /app /etc/cert-exporter

# Use non-root user
USER certexporter

# Expose prometheus metrics port
EXPOSE 9117

# Set default config location
ENV CONFIG_FILE=/etc/cert-exporter/config.yml

# Run the exporter
ENTRYPOINT ["./cert-exporter"]
CMD ["--config.file", "/etc/cert-exporter/config.yml"] 