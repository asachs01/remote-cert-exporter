# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache make git

# Copy only the files needed for go mod download first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN make build

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/remote-cert-exporter .

# Create a non-root user
RUN adduser -D -H -h /app exporter
USER exporter

# Expose metrics port
EXPOSE 9117

ENTRYPOINT ["/app/remote-cert-exporter"] 