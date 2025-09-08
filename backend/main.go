package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/process", process())

	if err := r.Run(":8080"); err != nil {
		panic("Error al iniciar el servidor: " + err.Error())
	}
}
