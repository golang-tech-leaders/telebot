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
COPY config.yml ./
COPY --from=builder /bin/server ./

CMD ["./server", "-config", "config.yml"]