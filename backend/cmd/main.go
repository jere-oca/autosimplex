package main

import (
	"autosimplex/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/process", handler.Process())

	if err := r.Run(":8080"); err != nil {
		panic("Error al iniciar el servidor: " + err.Error())
	}
}
