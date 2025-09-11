package simplex

import (
	"testing"
	"autosimplex/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewSolver(t *testing.T) {
	solver := NewSolver()
	assert.NotNil(t, solver)
}

func TestSolver_Solve(t *testing.T) {
	solver := NewSolver()
	
	request := models.SimplexRequest{
		Objective: models.Objective{
			Type:         "max",
			Coefficients: []float64{3, 2},
		},
		Constraints: []models.Constraint{
			{
				Coefficients: []float64{1, 1},
				Operator:     "<=",
				RHS:          4,
			},
		},
	}
	
	result, err := solver.Solve(request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "success", result["status"])
}