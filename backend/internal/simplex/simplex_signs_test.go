package simplex

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestSimplexEqualityConstraint(t *testing.T) {
	// Maximize 3 x1 + 2 x2
	maximize := mat.NewVecDense(2, []float64{3, 2})

	// Constraint: x1 + x2 = 4
	constraints := mat.NewDense(1, 3, []float64{
		1, 1, 4,
	})

	signs := []string{"="}

	result, sol, _ := SolveWithSigns(maximize, constraints, signs)

	expected := 12.0
	if result != expected {
		t.Fatalf("Expected optimal %v but got %v, solution: %v", expected, result, sol)
	}
	if len(sol) != 2 {
		t.Fatalf("Unexpected solution length: %v", sol)
	}
	if sol[0] != 4 || sol[1] != 0 {
		t.Fatalf("Expected solution [4,0] but got %v", sol)
	}
}

func TestSimplexGreaterEqualConstraint(t *testing.T) {
	// Maximize 2 x1 + 1 x2
	maximize := mat.NewVecDense(2, []float64{2, 1})

	// Constraints:
	// x1 + x2 >= 3
	// x1 <= 2
	// x2 <= 3
	constraints := mat.NewDense(3, 3, []float64{
		1, 1, 3,
		1, 0, 2,
		0, 1, 3,
	})

	signs := []string{">=", "<=", "<="}

	result, sol, _ := SolveWithSigns(maximize, constraints, signs)

	expected := 7.0
	if result != expected {
		t.Fatalf("Expected optimal %v but got %v, solution: %v", expected, result, sol)
	}
	if len(sol) != 2 {
		t.Fatalf("Unexpected solution length: %v", sol)
	}
	if sol[0] != 2 || sol[1] != 3 {
		t.Fatalf("Expected solution [2,3] but got %v", sol)
	}
}
