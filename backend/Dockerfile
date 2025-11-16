FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache git build-base
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /trading-backend ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /trading-backend /trading-backend

RUN mkdir -p /data
ENV DB_PATH=/data/trading.db
ENV JWT_SECRET=dev-secret

EXPOSE 8080
CMD ["/trading-backend"]
