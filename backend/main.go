package main

import (
	"autosimplex/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configurar CORS manualmente
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/process", handler.Process())

	if err := r.Run(":8080"); err != nil {
		panic("Error al iniciar el servidor: " + err.Error())
	}
}
