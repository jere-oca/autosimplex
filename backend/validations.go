package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func validateReqConstraints(c *gin.Context, rows int, cols int, vars []float64) bool {
	if rows <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad de filas debe ser mayor a 0"})
		return true
	}
	if cols <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad de columnas debe ser mayor a 0"})
		return true
	}
	if len(vars) != rows*cols {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad de variables no coincide con filas x columnas"})
		return true
	}
	for i, v := range vars {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Valor inválido en la matriz de restricciones en la posición %d", i)})
			return true
		}
	}
	return false
}

func validateReqObjective(c *gin.Context, n int, coefs []float64) bool {
	if n <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad de variables de decisión debe ser mayor a 1"})
		return true
	}
	if len(coefs) != n {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad de coeficientes no coincide con n"})
		return true
	}
	for i, v := range coefs {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Coeficiente inválido en la posición %d", i)})
			return true
		}
	}
	return false
}
