package handler

import (
	"autosimplex/internal/models"
	"autosimplex/internal/pdf"
	"autosimplex/internal/simplex"
	"fmt"
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

		// Debug: print what we are sending to the solver
		fmt.Printf("isMinimize=%v, objectiveVec=%v\n", isMinimize, mat.Formatted(maximizeVec, mat.Prefix(" ")))
		fmt.Printf("constraintMatrix:\n%v\n", mat.Formatted(constraintMatrix, mat.Prefix(" ")))

		// Allow a debug-only response that returns the transformed inputs
		// without executing the solver. Use query param ?debug=true.
		if c.Query("debug") == "true" {
			// Convert maximizeVec to a plain slice
			coeffs := make([]float64, maximizeVec.Len())
			for i := 0; i < maximizeVec.Len(); i++ {
				coeffs[i] = maximizeVec.At(i, 0)
			}
			// Convert constraintMatrix to nested slices
			rowsOut, colsOut := constraintMatrix.Dims()
			matOut := make([][]float64, rowsOut)
			for i := 0; i < rowsOut; i++ {
				row := make([]float64, colsOut)
				for j := 0; j < colsOut; j++ {
					row[j] = constraintMatrix.At(i, j)
				}
				matOut[i] = row
			}
			c.JSON(http.StatusOK, gin.H{
				"debug":       true,
				"isMinimize":  isMinimize,
				"objective":   coeffs,
				"constraints": matOut,
			})
			return
		}

		// Build signs slice: use provided signs if any, otherwise default to "<=" for all rows
		signs := req.Constraints.Signs
		if len(signs) == 0 {
			signs = make([]string, rows)
			for i := 0; i < rows; i++ {
				signs[i] = "<="
			}
		}

		result, solution := simplex.SolveWithSigns(maximizeVec, constraintMatrix, signs)

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
		})
	}
}
