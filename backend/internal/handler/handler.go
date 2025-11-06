package handler

import (
	"autosimplex/internal/models"
	"autosimplex/internal/pdf"
	"autosimplex/internal/simplex"
	"net/http"
	"strings"

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
		if validateReqObjective(c, n, coefs, req.Objective.Type) {
			return
		}
		// Build objective vector. If the request asks to minimize, convert
		// the problem into a maximization by negating the coefficients.
		objective := mat.NewVecDense(n, coefs)
		isMinimize := strings.ToLower(strings.TrimSpace(req.Objective.Type)) == "minimize"
		var maximizeVec *mat.VecDense
		if isMinimize {
			neg := make([]float64, n)
			for i := range n {
				neg[i] = -coefs[i]
			}
			maximizeVec = mat.NewVecDense(n, neg)
		} else {
			maximizeVec = objective
		}

		// Matriz de restricciones (incluye lado derecho)
		rows := req.Constraints.Rows
		cols := req.Constraints.Cols
		vars := req.Constraints.Vars
		if validateReqConstraints(c, rows, cols, vars) {
			return
		}
		constraintMatrix := mat.NewDense(rows, cols, vars)

		// Build signs slice: use provided signs if any, otherwise default to "<=" for all rows
		signs := req.Constraints.Signs
		if len(signs) == 0 {
			signs = make([]string, rows)
			for i := range rows {
				signs[i] = "<="
			}
		}

		result, solution, steps, warning := simplex.SolveWithSigns(maximizeVec, constraintMatrix, signs)

		// If it was a minimization request, invert the returned optimal value
		// because we solved the equivalent maximization of -c.
		if isMinimize {
			result = -result
		}

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
			"steps":         steps,
			"warning":       warning,
		})
	}
}
