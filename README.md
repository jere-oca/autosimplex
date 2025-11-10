# AutoSimplex

> Algoritmo simplex para resolución de problemas de programación lineal.

## Metodología

Utilizamos la metodología SCRUM, mediante iteraciones de 2 semanas. Cada iteración cuenta con una planificación, al menos una reunión semanal de seguimiento y una retrospectiva al finalizar.

### Integrantes

- Bayinay, Federico
- Ocaña, Jeremias
- Sabio, Santiago
- Sanz, Lautaro

## Tecnologías

Backend:
  - [Go](https://go.dev/) + [Gin](https://gin-gonic.com/).
  - Librería: [Gonum](https://www.gonum.org/)

Frontend:
  - Vite + Preact

# Instalación

## Docker

> Se requiere [Docker](https://www.docker.com/get-started/).

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/jere-oca/autosimplex
   cd autosimplex/
   ```

2. Levantar los contenedores:
 
   ```bash
   docker compose up -d
   ```
   
El frontend se ejecutará en http://localhost:3000/.

## Local

> Se requiere [Go](https://go.dev/doc/install) y [Node.js](https://nodejs.org/es/download).

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/jere-oca/autosimplex
   cd autosimplex/backend/
   ```

2. Instalar las dependencias:
 
   ```bash
   go mod tidy
   ```
   
3. Ejecutar el servidor:

   ```bash
   go run .
   ```

La API se ejecutará en http://localhost:8080/.
   
4. Correr el frontend:

   ```bash
   cd ../frontend
   ```
   
   ```bash
   npm run dev
   ```

Se ejecutará en http://localhost:5173/.
