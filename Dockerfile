FROM golang:1.24.13-alpine

RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    ca-certificates \
    tzdata

RUN go install github.com/go-delve/delve/cmd/dlv@v1.24.0
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ENV PATH="/go/bin:${PATH}"

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download && go mod verify

COPY . .

EXPOSE 8080
EXPOSE 2345

RUN go build -gcflags="all=-N -l" -o /tmp/server ./cmd/server
CMD ["dlv", "exec", "--headless", "--api-version=2", "--listen=:2345", "--accept-multiclient", "--continue", "/tmp/server"]