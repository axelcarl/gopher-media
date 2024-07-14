FROM golang:1.22.5-alpine

# Install Air for hot-reloading.
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Command to run Air to serve the application.
CMD ["air", "-c", ".air.toml"]

EXPOSE 8080
