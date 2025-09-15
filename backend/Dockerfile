# Etapa 1: Compilación
FROM golang:1.25 AS builder

# Configuramos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos go.mod y go.sum primero (mejora el cache)
COPY go.mod go.sum ./

# Descargamos dependencias
RUN go mod download

# Copiamos el resto del código fuente
COPY . .

# Compilamos la aplicación en un binario estático
RUN go build -o main .

# Etapa 2: Imagen final
FROM debian:bookworm-slim

# Creamos un directorio de trabajo
WORKDIR /app

# Copiamos el binario desde la etapa de compilación
COPY --from=builder /app/main .

# Exponemos el puerto (ejemplo: 8080)
EXPOSE 8080

# Ejecutamos el binario
CMD ["./main"]