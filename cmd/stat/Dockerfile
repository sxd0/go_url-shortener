FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/stat
RUN go build -o /stat

FROM alpine

WORKDIR /app

COPY --from=builder /stat .

ENTRYPOINT ["/app/stat"]
