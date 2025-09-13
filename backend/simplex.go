package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func Solve(maximize mat.Vector, constraints *mat.Dense) (float64, []float64) {

	// Dimensiones de la matriz de restricciones sin los valores del lado derecho
	constraintCount, variablesCount := constraints.Dims()
	variablesCount--

	totalVariables := constraintCount + variablesCount

	// Matriz A: coeficientes de las restricciones
	A := mat.DenseCopyOf(constraints.Grow(0, constraintCount-1))

	// Se asignan las variables de holgura
	tempVector := make([]float64, constraintCount, constraintCount)
	tempVector[0] = 1
	A.SetCol(variablesCount, tempVector)

	// Resto de la matriz identidad
	for i := 1; i < constraintCount; i++ {
		A.Set(i, i+variablesCount, 1)
	}

	// Vector c: coeficientes de la función objetivo y variables de holgura cero
	c := mat.NewDense(1, totalVariables, make([]float64, totalVariables, totalVariables))
	for i := 0; i < maximize.Len(); i++ {
		c.Set(0, i, maximize.At(i, 0))
	}

	// Vector b: lado derecho de las restricciones
	bTemp := make([]float64, constraintCount, constraintCount)
	for i := 0; i < constraintCount; i++ {
		bTemp[i] = constraints.At(i, variablesCount)
	}
	b := mat.NewVecDense(constraintCount, bTemp)

	// Variables básicas iniciales
	currentBaseVars := make([]int, constraintCount, constraintCount)
	for i := range currentBaseVars {
		currentBaseVars[i] = variablesCount + i + 1
	}

	// Se utilizarán para mostrar el progreso al usuario
	//fmt.Printf("Current base vars:\n %v\n\n", currentBaseVars)
	//
	//fmt.Printf("A matrix:\n %v\n\n", mat.Formatted(A, mat.Prefix(" "), mat.Excerpt(8)))
	//
	//fmt.Printf("c vector:\n %v\n\n", mat.Formatted(c, mat.Prefix(" "), mat.Excerpt(8)))
	//
	//fmt.Printf("b vector:\n %v\n\n", mat.Formatted(b, mat.Prefix(" "), mat.Excerpt(8)))

	// Iteraciones
	const maxSimplexIterations = 10
	iterations := 0
	for {
		if iterations > maxSimplexIterations {
			break
		}

		// Paso 1: resolver el sistema dual (y^T)(B) = (c_B)^T
		// Se construye la matriz B de columnas básicas y se obtienen los multiplicadores simplex 'y'
		B := mat.NewDense(constraintCount, constraintCount, nil)
		AT := mat.DenseCopyOf(A.T()) // Transpuesta
		cBData := make([]float64, constraintCount, constraintCount)
		for i := range currentBaseVars {
			B.SetCol(i, AT.RawRowView(currentBaseVars[i]-1))
			cBData[i] = c.At(0, currentBaseVars[i]-1)
		}
		//fmt.Printf("B matrix:\n %v\n\n", mat.Formatted(B, mat.Prefix(" "), mat.Excerpt(8)))
		y := mat.NewDense(1, constraintCount, cBData)
		//fmt.Printf("cBT vector:\n %v\n\n", mat.Formatted(y, mat.Prefix(" "), mat.Excerpt(8)))

		Bi := mat.DenseCopyOf(B)
		err := Bi.Inverse(B)
		if err != nil {
			panic("Error en matriz inversa!")
		}
		//fmt.Printf("Bi matrix:\n %v\n\n", mat.Formatted(Bi, mat.Prefix(" "), mat.Excerpt(8)))
		y.Mul(y, Bi)
		//fmt.Printf("y^T vector:\n %v\n\n", mat.Formatted(y, mat.Prefix(" "), mat.Excerpt(8)))

		// Paso 2: calcular y^T A_N y comparar con c_{N}^T component-wise
		// Paso 2: calcular costos reducidos
		// Matriz de variables no básicas
		AN := mat.NewDense(constraintCount, variablesCount, nil)
		cNT := mat.NewDense(1, variablesCount, nil)
		var currentNonBaseVars []int
		for i := 1; i < totalVariables+1; i++ {
			if !contains(currentBaseVars, i) {
				currentNonBaseVars = append(currentNonBaseVars, i)
			}
		}
		//fmt.Printf("Non-Base vars:\n %v\n\n", currentNonBaseVars)

		for i := range currentNonBaseVars {
			AN.SetCol(i, AT.RawRowView(currentNonBaseVars[i]-1))
			cNT.SetCol(i, []float64{c.At(0, currentNonBaseVars[i]-1)})
		}

		// Costos reducidos (si todos <= 0, la solución es óptima)
		yTAN := mat.NewDense(1, variablesCount, nil)
		yTAN.Mul(y, AN)

		//fmt.Printf("y^T A_N vector:\n %v\n\n", mat.Formatted(yTAN, mat.Prefix(" "), mat.Excerpt(8)))
		//fmt.Printf("AN matrix:\n %v\n\n", mat.Formatted(AN, mat.Prefix(" "), mat.Excerpt(8)))
		//fmt.Printf("cNT vector:\n %v\n\n", mat.Formatted(cNT, mat.Prefix(" "), mat.Excerpt(8)))

		newBaseVar := variablesCount + constraintCount + 1
		var largestVal float64
		hasLargestVal := false
		a := mat.NewDense(constraintCount, 1, nil)

		// Paso 3: elegir variable entrante
		for i := range currentNonBaseVars {
			if cNT.At(0, i) > yTAN.At(0, i) {
				// Mayor que el máximo valor actual y de índice menor que el del máximo actual
				if !hasLargestVal || cNT.At(0, i) >= largestVal {
					if currentNonBaseVars[i] < newBaseVar {
						newBaseVar = currentNonBaseVars[i]
						largestVal = cNT.At(0, i)
						hasLargestVal = true
					}
				}
			}
		}
		//fmt.Printf("new base var:\n %v\n\n", newBaseVar)

		// Si no hay mejora posible, se devuelve la función objetivo óptima
		if !hasLargestVal {
			var result float64
			solution := make([]float64, maximize.Len())
			for i := range currentBaseVars {
				baseVarIndex := currentBaseVars[i] - 1
				if baseVarIndex < variablesCount { // evitar variables de holgura
					//fmt.Printf("b vector:\n %v\n\n", mat.Formatted(b, mat.Prefix(" "), mat.Excerpt(8)))
					val := b.At(i, 0)
					result += maximize.At(baseVarIndex, 0) * val
					solution[baseVarIndex] = val
				}
			}
			return result, solution
		}

		// Paso 4: calcular dirección y paso permitido
		a.SetCol(0, AT.RawRowView(newBaseVar-1))

		//fmt.Printf("a vector:\n %v\n\n", mat.Formatted(a, mat.Prefix(" "), mat.Excerpt(8)))

		a.Mul(Bi, a)

		//fmt.Printf("d vector:\n %v\n\n", mat.Formatted(a, mat.Prefix(" "), mat.Excerpt(8)))

		// Paso 5: máximo 't' posible tal que b - t * d <= 0 (determina qué variable sale de la base)
		lowest := -1.0
		lowestIndex := -1
		lowestValueOfT := 0.0
		for i := range currentBaseVars {
			baseValue := b.At(i, 0)
			dValue := a.At(i, 0)
			if dValue > 0 {
				tValue := baseValue / dValue
				//fmt.Println(tValue)
				if lowest < 0 || tValue < lowest {
					lowest = tValue
					lowestIndex = i
					lowestValueOfT = tValue
				}
			}
		}
		if lowest <= 0 {
			fmt.Println("couldn't find appropriate t value")
			return 0, nil
		}

		// Paso 6: actualizar la base y el vector 'b'
		for i := range currentBaseVars {
			if i == lowestIndex {
				b.SetVec(i, lowest)
				currentBaseVars[lowestIndex] = newBaseVar
			} else {
				b.SetVec(i, b.At(i, 0)-lowestValueOfT*a.At(i, 0))
			}
		}
		//fmt.Printf("new b vector:\n %v\n\n", mat.Formatted(b, mat.Prefix(" "), mat.Excerpt(8)))
		//fmt.Printf("new base vars:\n %v\n\n", currentBaseVars)
		//
		//fmt.Println("--------------------------------------------------------")
		//fmt.Printf("iteration: %v\n", iterations)

		iterations++
	}

	return 0, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
