FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/auth
RUN go build -o /auth

FROM alpine

WORKDIR /app

COPY --from=builder /auth .

ENTRYPOINT ["/app/auth"]
