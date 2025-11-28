# --- Stage 1: Build ---
FROM docker.io/library/golang:1.25-alpine AS builder
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
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-s -w" -o glyphtones .

# --- Stage 2: Runtime ---
FROM docker.io/library/alpine:3.22
# Install ffmpeg and certificates
RUN apk add --no-cache ffmpeg ca-certificates
WORKDIR /app
# Copy binary and static from builder
COPY --from=builder /app/glyphtones ./glyphtones
COPY --from=builder /app/static ./static
# Expose backend port
ENV LISTEN_PORT=8080
EXPOSE $LISTEN_PORT
# Run the backend
ENTRYPOINT ["./glyphtones"]
