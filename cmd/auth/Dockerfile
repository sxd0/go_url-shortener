FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/auth
RUN go build -o /auth

FROM gcr.io/distroless/base-debian11

COPY --from=builder /auth /auth

ENV AUTH_GRPC_PORT=50051

ENTRYPOINT ["/auth"]
