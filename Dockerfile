# --- Stage 1: Build ---
FROM docker.io/library/golang:alpine AS builder
# Install git and build tools (go needs that to get some dependencies)
RUN apk add --no-cache git ca-certificates
# Set working directory
WORKDIR /app
# Copy Go modules manifests first for caching
COPY go.mod go.sum ./
RUN go mod download
# Get Templ
RUN go install github.com/a-h/templ/cmd/templ@latest
# Copy the rest of the source code
COPY . .
# Build the Go binary for Linux ARM
RUN templ generate
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o glyphtones .

# --- Stage 2: Runtime ---
FROM docker.io/library/alpine:latest
# Install ffmpeg, certificates and curl for healthcheck
RUN apk add --no-cache ffmpeg ca-certificates curl
WORKDIR /app
# Copy binary and static from builder
COPY --from=builder /app/glyphtones ./glyphtones
COPY --from=builder /app/static ./static
# Expose backend port
EXPOSE 8080
# Define healthcheck
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1
# Run the backend
ENTRYPOINT ["./glyphtones"]
