# Stage 1: Builder
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ports ./cmd

# Stage 2: Final
FROM alpine:latest
COPY --from=builder /app/ports /usr/local/bin/ports
CMD ["ports"]