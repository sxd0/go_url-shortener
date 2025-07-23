FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/gateway
RUN go build -o /gateway

FROM alpine

WORKDIR /app

COPY --from=builder /gateway .

ENTRYPOINT ["/app/gateway"]
