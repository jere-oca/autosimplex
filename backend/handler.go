package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

func process() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req MatrixRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Procesar 'req.Matrix'
		c.JSON(http.StatusOK, gin.H{"received_matrix": req.Matrix})
	}
}
