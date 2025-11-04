package handler

import (
	"autosimplex/internal/models"
	"autosimplex/internal/pdf"
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

			result, solution := simplex.Solve(objective, constraintMatrix)

			// Si se solicita formato PDF, generar y devolver PDF
			format := c.Query("format")
			if format == "pdf" {
				c.Writer.Header().Set("Content-Type", "application/pdf")
				c.Writer.Header().Set("Content-Disposition", "attachment; filename=resultado_simplex.pdf")
				if err := pdf.GenerateSimplexPDF(result, solution, c.Writer); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"optimal_value": result,
				"solution":      solution,
			})
	}
}
