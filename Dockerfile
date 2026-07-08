FROM golang:1.24.13-alpine

RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    ca-certificates \
    tzdata \
    && go install github.com/go-delve/delve/cmd/dlv@v1.24.0

ENV PATH="/go/bin:${PATH}"

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum* ./
RUN go mod download && go mod verify

# Копируем исходники
COPY . .

EXPOSE 8080
EXPOSE 2345

CMD ["dlv", "debug", "--headless", "--api-version=2", "--listen=:2345", "--accept-multiclient", "--continue", "--build-flags=-buildvcs=false", "./cmd/server"]