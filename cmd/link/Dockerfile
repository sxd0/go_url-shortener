FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/link
RUN go build -o /link

FROM alpine

WORKDIR /app

COPY --from=builder /link .

ENTRYPOINT ["/app/link"]
