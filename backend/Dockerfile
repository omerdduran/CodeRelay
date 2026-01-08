# syntax=docker/dockerfile:1.7

FROM golang:1.24 AS builder
WORKDIR /src

# Install build dependencies for CGO + SQLite
RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
RUN go mod download

COPY . .
# Build with CGO enabled so github.com/mattn/go-sqlite3 works
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/api ./cmd/api

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /srv/api
COPY --from=builder /out/api /usr/local/bin/api
COPY sqlite ./sqlite

ENV API_ADDR=:8080
EXPOSE 8080

CMD ["api"]
