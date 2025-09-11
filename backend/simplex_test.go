package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestSimplexPositives(t *testing.T) {
	// create new example matrix
	maximize := mat.NewVecDense(3, []float64{5, 4, 3})

	constraints := mat.NewDense(3, 4, []float64{
		2, 3, 1, 5,
		4, 1, 2, 11,
		3, 4, 2, 8})

	actualResult := Solve(maximize, constraints)
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

	actualResult := Solve(maximize, constraints)
	var expectedResult float64 = 7

	if actualResult != expectedResult {
		t.Fatalf("Expected %f but got %f", expectedResult, actualResult)
	}
}

func TestSimplexFractions(t *testing.T) {
	// TODO: usar fracciones
	// create new example matrix
	maximize := mat.NewVecDense(3, []float64{2, 3, 2})

	constraints := mat.NewDense(3, 4, []float64{
		1, 1, 0, 8,
		0, 1, 2, 12,
		0, 1, 1, 7})

	actualResult := Solve(maximize, constraints)
	var expectedResult float64 = 28

	if actualResult != expectedResult {
		t.Fatalf("Expected %f but got %f", expectedResult, actualResult)
	}
}
