# syntax=docker/dockerfile:1.7

FROM golang:1.22 AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /srv/api
COPY --from=builder /out/api /usr/local/bin/api
COPY sqlite ./sqlite

ENV API_ADDR=:8080
EXPOSE 8080

CMD ["api"]
