package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

func main() {
	r := gin.Default()

	r.POST("/process", func(c *gin.Context) {
		var req MatrixRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Procesar 'req.Matrix'
		c.JSON(http.StatusOK, gin.H{"matriz_recibida": req.Matrix})
	})

	if err := r.Run(":8080"); err != nil {
		panic("Error al iniciar el servidor: " + err.Error())
	}
}
