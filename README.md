# AutoSimplex

> Algoritmo simplex para resolución de problemas de programación lineal.

## Metodología

Utilizamos la metodología SCRUM, mediante iteraciones de 2 semanas. Cada iteración cuenta con una planificación, al menos una reunión semanal de seguimiento y una retrospectiva al finalizar.

## Tecnologías

- Lenguaje + framework propuesto: [Go](https://go.dev/) + [Gin](https://gin-gonic.com/).
- Librería: [Gonum](https://www.gonum.org/)

### Integrantes

- Bayinay, Federico
- Ocaña, Jeremias
- Sabio, Santiago
- Sanz, Lautaro

# Instalación

> Se requiere tener instalado [Go](https://go.dev/doc/install).

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/jere-oca/autosimplex
   cd autosimplex/backend/
   ```

## Local

2. Instalar las dependencias:
 
   ```bash
   go mod tidy
   ```
   
3. Ejecutar el servidor:

   ```bash
   go run .
   ```
   - No correr `run main.go`, ya que se requieren más archivos que el principal.

## Docker

> Requiere tener instalado [Docker](https://www.docker.com/get-started/)

1. Construir imagen:

    ```bash
    docker build -t  autosimplex ./backend/
    ```

2. Correr:

    ```bash
    docker run -it --rm -p 8080:8080 autosimplex
    ```

# Testing

`curl` de prueba:

```bash
curl -X POST http://localhost:8080/process \
-H "Content-Type: application/json" \
-d '{
    "objective": {
    "n": 4,
    "coefficients": [3.2, 0.75, 5, 7.8]
    },
    "constraints": {
    "rows": 4,
    "cols": 5,
    "vars": [1, 1.5, 2, 3, 4,
        0, 1, 2.5, 6.3, 8,
        0, 1, 1, 0.8, 7,
        1, 5, 2.1, 3, 13]
    }
}'
```
