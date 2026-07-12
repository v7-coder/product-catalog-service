FROM golang:1.24.13-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /tmp/server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /tmp/server .
COPY migrations ./migrations

EXPOSE 8080
CMD ["./server"]