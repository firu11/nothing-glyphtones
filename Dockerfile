# --- Stage 1: Build ---
FROM docker.io/library/golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go tool templ generate
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o glyphtones .

# --- Stage 2: Runtime ---
FROM docker.io/library/alpine:latest

RUN apk add --no-cache ffmpeg ca-certificates curl
WORKDIR /app

COPY --from=builder /app/glyphtones ./glyphtones
COPY --from=builder /app/static ./static

EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1

ENTRYPOINT ["./glyphtones"]
