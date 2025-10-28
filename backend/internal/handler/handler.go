package handler

import (
	"autosimplex/internal/models"
	"autosimplex/internal/simplex"
	"net/http"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"
)

func Process() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req models.SimplexRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()})
			return
		}

		// Función objetivo
		n := req.Objective.N
		coefs := req.Objective.Coefficients
		if validateReqObjective(c, n, coefs) {
			return
		}
		objective := mat.NewVecDense(n, coefs)

		// Matriz de restricciones (incluye lado derecho)
		rows := req.Constraints.Rows
		cols := req.Constraints.Cols
		vars := req.Constraints.Vars
		if validateReqConstraints(c, rows, cols, vars) {
			return
		}
		constraintMatrix := mat.NewDense(rows, cols, vars)

		// Ahora Solve retorna también los pasos
		result, solution, steps := simplex.Solve(objective, constraintMatrix)

		c.JSON(http.StatusOK, gin.H{
			"optimal_value": result,
			"solution":      solution,
			"steps":         steps,
		})
	}
}
