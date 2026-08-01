# --- Build stage ---
FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Needed for the cgo-based sqlite3 driver
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
RUN go build -o /minibooking .

# --- Runtime stage ---
FROM debian:bookworm-slim

WORKDIR /app

# sqlite3 driver is statically linked via cgo but still needs libc; slim
# debian already has it. curl is included only to make the HEALTHCHECK work.
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /minibooking /app/minibooking

# The sqlite file is written here; mount a volume at /app/data in
# docker-compose so booking data survives container restarts.
ENV DB_PATH=/app/data/minibooking.db
RUN mkdir -p /app/data

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD curl -f http://localhost:8080/api/schedules || exit 1

CMD ["/app/minibooking"]
