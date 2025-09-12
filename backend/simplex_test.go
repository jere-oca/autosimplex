package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func CompareFloat64Slice(a []float64, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSimplexPositives(t *testing.T) {
	// create new example matrix
	maximize := mat.NewVecDense(3, []float64{5, 4, 3})

	constraints := mat.NewDense(3, 4, []float64{
		2, 3, 1, 5,
		4, 1, 2, 11,
		3, 4, 2, 8})

	actualResult, _ := Solve(maximize, constraints)
	var expectedResult float64 = 13

	if actualResult != expectedResult {
		t.Fatalf("Expected %f but got %f", expectedResult, actualResult)
	}
}

func TestSimplexNegatives(t *testing.T) {
	// create new example matrix
	maximize := mat.NewVecDense(4, []float64{6, -9, 1, -11})

	constraints := mat.NewDense(2, 5, []float64{
		2, -3, -1, -7, 1,
		2, 1, 1, 3, 3})

	actualResult, _ := Solve(maximize, constraints)
	var expectedResult float64 = 7

	if actualResult != expectedResult {
		t.Fatalf("Expected %f but got %f", expectedResult, actualResult)
	}
}
func TestSimplexFractions(t *testing.T) {
	// create new example matrix
	maximize := mat.NewVecDense(4, []float64{3.2, .75, 5, 7.8})

	constraints := mat.NewDense(4, 5, []float64{
		1, 1.5, 2, 3, 4,
		0, 1, 2.5, 6.3, 8,
		0, 1, 1, .8, 7,
		1, 5, 2.1, 3, 13})

	actualResult, actualSolution := Solve(maximize, constraints)
	var expectedResult float64 = 12.8
	var expectedSolution = []float64{4, 0, 0, 0}

	if actualResult != expectedResult {
		t.Fatalf("Expected %f but got %f", expectedResult, actualResult)
	}

	if !CompareFloat64Slice(actualSolution, expectedSolution) {
		t.Fatalf("Expected %v but got %v", expectedSolution, actualSolution)
	}
}
