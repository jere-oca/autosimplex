FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod tidy

# Compilación estática y optimizada
RUN CGO_ENABLED=0 go build -o main -ldflags="-s -w"

# Imagen final
FROM gcr.io/distroless/static-debian12

USER nobody

ENV HOME=/home/nobody
ENV GIN_MODE=release

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080

CMD ["/app/main"]
