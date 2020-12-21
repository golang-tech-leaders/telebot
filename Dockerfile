FROM golang:1.15-buster AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN go build -v -o /bin/server cmd/botapp/main.go

FROM debian:buster-slim
RUN set -x && apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY cmd/botapp/config.yml ./config.yml
COPY --from=builder /bin/server ./
COPY db/migrations /app/migrations

ENV DATABASE_URL=""

EXPOSE $PORT

CMD ["./server", "-config", "config.yml"]