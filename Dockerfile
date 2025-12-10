
FROM golang:1.24.0 AS builder


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o /app/Telegram_bot ./cmd/main.go


FROM ubuntu:latest


WORKDIR /app


RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        tzdata \                
        ca-certificates \       
        curl \                
    && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Ekaterinburg
RUN ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    chmod 644 /etc/timezone


COPY --from=builder /app/Telegram_bot /app/Telegram_bot



CMD ["/app/Telegram_bot"]