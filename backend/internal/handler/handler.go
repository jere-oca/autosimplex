package handler

import (
	"autosimplex/internal/models"
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
			for i := 0; i < n; i++ {
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

		result, solution := simplex.Solve(maximizeVec, constraintMatrix)

		// If it was a minimization request, invert the returned optimal value
		// because we solved the equivalent maximization of -c.
		if isMinimize {
			result = -result
		}

		c.JSON(http.StatusOK, gin.H{
			"optimal_value": result,
			"solution":      solution,
		})
	}
}
