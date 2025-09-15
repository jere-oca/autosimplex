# AutoSimplex

> Algoritmo simplex para resolución de problemas de programación lineal.

## Metodología

Utilizamos una metodología SCRUM, mediante iteraciones de 2 semanas. Cada iteración cuenta con una planificación, al menos una reunión semanal de seguimiento y una retrospectiva al finalizar.

## Tecnologías

- Lenguaje + framework propuesto: [Go](https://go.dev/) + [Gin](https://gin-gonic.com/).
- Libería: [Gonum](https://www.gonum.org/)

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
   cd autosimplex
   ```
   
2. Instalar las dependencias:
   ```bash
   go mod tidy
   ```
   
3. Ejecutar el servidor:
   ```bash
   go run .
   ```
   - No correr `run main.go`, ya que se requieren más archivos que el principal.

3.1. Ejecutar el Container:
```bash 

docker build -t  autosimplex:1.25 .

docker run -it --rm -p 8080:8080 autosimplex:1.25

```
   

4. Podrá acceder a la API en `http://localhost:8080`.



