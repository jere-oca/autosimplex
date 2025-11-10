package simplex

import (
	"math"
	"slices"

	"gonum.org/v1/gonum/mat"
)

// Solve es una función auxiliar que asume que todas las restricciones son "<=".
func Solve(maximize mat.Vector, constraints *mat.Dense) (float64, []float64, []SimplexStep, string) {
	rows, _ := constraints.Dims()
	signs := make([]string, rows)
	for i := range signs {
		signs[i] = "<="
	}
	return SolveWithSigns(maximize, constraints, signs)
}

// SolveWithSigns resuelve un problema de maximización de PL dado un vector objetivo y una
// matriz de restricciones (las filas son [a1 ... an b]). 'signs' contiene uno de
// "<=", ">=", o "=" por restricción. Usa una estrategia Big-M para variables artificiales.
func SolveWithSigns(maximize mat.Vector, constraints *mat.Dense, signs []string) (float64, []float64, []SimplexStep, string) {
	const M = 1e7
	var steps []SimplexStep

	m, cols := constraints.Dims()
	n := maximize.Len()
	if cols != n+1 {
		warning := "Cantidad de columnas no coinciden con variables"
		return 0, nil, steps, warning
	}

	// Contar variables extra y construir matriz A extendida
	// Agregaremos holgura (para <=), exceso+artificial (para >=), y artificial (para =)
	extraCols := 0
	for i := range m {
		s := "<="
		if i < len(signs) {
			s = signs[i]
		}
		switch s {
		case "<=":
			extraCols += 1 // holgura
		case ">=":
			extraCols += 2 // exceso + artificial
		case "=":
			extraCols += 1 // artificial
		default:
			extraCols += 1
		}
	}

	totalVars := n + extraCols

	// Construir A_extendida
	A := mat.NewDense(m, totalVars, nil)
	// Llenar variables originales
	for i := range m {
		for j := range n {
			A.Set(i, j, constraints.At(i, j))
		}
	}

	// Rastrear índices y variable base por fila
	col := n
	baseVars := make([]int, m) // índices base 1 de variables básicas por fila
	artIndices := []int{}
	for i := range m {
		s := "<="
		if i < len(signs) {
			s = signs[i]
		}
		switch s {
		case "<=":
			A.Set(i, col, 1)
			baseVars[i] = col + 1
			col++
		case ">=":
			// exceso
			A.Set(i, col, -1)
			col++
			// artificial
			A.Set(i, col, 1)
			baseVars[i] = col + 1
			artIndices = append(artIndices, col)
			col++
		case "=":
			// artificial
			A.Set(i, col, 1)
			baseVars[i] = col + 1
			artIndices = append(artIndices, col)
			col++
		default:
			A.Set(i, col, 1)
			baseVars[i] = col + 1
			col++
		}
	}

	// Construir objetivo c (1 x totalVars). Variables artificiales obtienen penalidad -M
	c := mat.NewDense(1, totalVars, make([]float64, totalVars))
	for j := range n {
		c.Set(0, j, maximize.At(j, 0))
	}
	for _, ai := range artIndices {
		c.Set(0, ai, -M)
	}

	// Construir vector b
	bData := make([]float64, m)
	for i := range m {
		bData[i] = constraints.At(i, n)
	}
	b := mat.NewVecDense(m, bData)

	// Ahora proceder con iteraciones simplex similar a implementación anterior,
	// pero usando nuestras A, c, b y baseVars construidas.
	ATrans := mat.DenseCopyOf(A.T())

	const maxIter = 200
	iter := 0
	for {
		if iter > maxIter {
			break
		}

		// Construir B desde baseVars
		B := mat.NewDense(m, m, nil)
		cB := make([]float64, m)
		for i := range m {
			// baseVars almacena índice base 1
			B.SetCol(i, ATrans.RawRowView(baseVars[i]-1))
			cB[i] = c.At(0, baseVars[i]-1)
		}

		// Resolver B^T * y = cB  (y es columna)
		var lu mat.LU
		lu.Factorize(B)
		yCol := mat.NewVecDense(m, nil)
		cBVec := mat.NewVecDense(m, cB)
		if err := lu.SolveVecTo(yCol, true, cBVec); err != nil {
			// base singular -> infactible
			warning := "Matriz singular, problema infactible o mal planteado"
			return 0, nil, steps, warning
		}
		// y como fila
		y := mat.NewDense(1, m, nil)
		for i := range m {
			y.Set(0, i, yCol.AtVec(i))
		}

		// Construir AN y cN para variables no básicas
		var nonBase []int
		for j := 1; j <= totalVars; j++ {
			if !contains(baseVars, j) {
				nonBase = append(nonBase, j)
			}
		}
		AN := mat.NewDense(m, len(nonBase), nil)
		cN := mat.NewDense(1, len(nonBase), nil)
		for i := range nonBase {
			AN.SetCol(i, ATrans.RawRowView(nonBase[i]-1))
			cN.SetCol(i, []float64{c.At(0, nonBase[i]-1)})
		}

		// Costos reducidos: cN - y * AN
		yAN := mat.NewDense(1, len(nonBase), nil)
		yAN.Mul(y, AN)

		// Elegir variable entrante: cualquier índice donde cN > yAN
		entering := -1
		enteringVal := 0.0
		for i := range nonBase {
			if cN.At(0, i) > yAN.At(0, i)+1e-9 {
				val := cN.At(0, i) - yAN.At(0, i)
				if entering == -1 || val > enteringVal || nonBase[i] < nonBase[entering] {
					entering = i
					enteringVal = val
				}
			}
		}

		if entering == -1 {
			// Verificar si hay variables artificiales en la base (problema infactible)
			for i := range m {
				bv := baseVars[i] - 1
				// Si la variable básica es artificial y tiene valor positivo
				if contains(artIndices, bv) && b.At(i, 0) > 1e-9 {
					warning := "Problema infactible: no existe solución"
					// Devolver solución parcial alcanzada hasta el momento
					solution := make([]float64, n)
					for j := range m {
						bvj := baseVars[j] - 1
						if bvj < n {
							val := b.At(j, 0)
							solution[bvj] = val
						}
					}
					return 0, solution, steps, warning
				}
			}

			// Verificar si hay infinitas soluciones (costo reducido = 0 para variables no básicas)
			hasInfiniteSolutions := false
			for i := range nonBase {
				// Si el costo reducido es cero para una variable no básica
				if math.Abs(cN.At(0, i)-yAN.At(0, i)) < 1e-9 {
					hasInfiniteSolutions = true
					break
				}
			}

			// Si llegamos aquí, la solución es óptima
			// Construir solución para variables originales
			solution := make([]float64, n)
			var optimal float64
			for i := range m {
				bv := baseVars[i] - 1
				if bv < n {
					// variable original
					val := b.At(i, 0)
					solution[bv] = val
					optimal += c.At(0, bv) * val
				}
			}

			warning := ""
			if hasInfiniteSolutions {
				warning = "Solución óptima no única: existen infinitas soluciones"
			}
			return optimal, solution, steps, warning
		}

		enteringVar := nonBase[entering]

		// Obtener columna a para enteringVar
		raw := ATrans.RawRowView(enteringVar - 1)
		aVec := mat.NewVecDense(m, nil)
		for i := range m {
			aVec.SetVec(i, raw[i])
		}

		// Resolver d = B^{-1} * aVec usando LU
		dVec := mat.NewVecDense(m, nil)
		if err := lu.SolveVecTo(dVec, false, aVec); err != nil {
			warning := "Solución no única, problema infactible o degenerado"
			return 0, nil, steps, warning
		}

		// Prueba de razón b_i / d_i para d_i > 0
		minRatio := math.Inf(1)
		leavingIndex := -1
		for i := range m {
			dv := dVec.AtVec(i)
			if dv > 1e-12 {
				ratio := b.At(i, 0) / dv
				if ratio < minRatio {
					minRatio = ratio
					leavingIndex = i
				}
			}
		}
		if leavingIndex == -1 {
			warning := "Problema no acotado"
			return 0, nil, steps, warning
		}

		// Preparar paso
		currentBaseVars := append([]int{}, baseVars...)
		currentNonBaseVars := append([]int{}, nonBase...)
		reducedCosts := matDenseToSlice(yAN)
		bVector := matVecToSlice(b)
		newBaseVar := enteringVar
		leavingVar := baseVars[leavingIndex]
		lowestValueOfT := minRatio

		// Construir tableau completo: para cada columna de variable j, calcular B^{-1} * A[:, j]
		// las filas del tableau contendrán coeficientes para cada variable y el RHS final
		tableRows := make([][]float64, m)
		for i := 0; i < m; i++ {
			tableRows[i] = make([]float64, totalVars+1) // +1 para R
		}
		// encabezado cj
		cj := make([]float64, totalVars)
		for j := 0; j < totalVars; j++ {
			cj[j] = c.At(0, j)
			// resolver B^{-1} * A[:, j]
			rawCol := ATrans.RawRowView(j)
			colVec := mat.NewVecDense(m, nil)
			aVec := mat.NewVecDense(m, nil)
			for i := 0; i < m; i++ {
				aVec.SetVec(i, rawCol[i])
			}
			if err := lu.SolveVecTo(colVec, false, aVec); err != nil {
				// si la resolución falla, llenar con ceros y continuar
				for i := 0; i < m; i++ {
					tableRows[i][j] = 0
				}
			} else {
				for i := 0; i < m; i++ {
					tableRows[i][j] = colVec.AtVec(i)
				}
			}
		}
		// llenar columna RHS
		for i := 0; i < m; i++ {
			tableRows[i][totalVars] = b.At(i, 0)
		}

		// cb: coeficientes objetivo de variables básicas
		cb := make([]float64, m)
		for i := 0; i < m; i++ {
			cb[i] = c.At(0, currentBaseVars[i]-1)
		}

		step := SimplexStep{
			Iteration:        iter,
			BaseVariables:    currentBaseVars,
			NonBaseVariables: currentNonBaseVars,
			ReducedCosts:     reducedCosts,
			BVector:          bVector,
			EnteringVar:      newBaseVar,
			LeavingVar:       leavingVar,
			TValue:           lowestValueOfT,
			Table:            tableRows,
			Cj:               cj,
			Cb:               cb,
			PivotRow:         leavingIndex,
			PivotCol:         enteringVar - 1,
		}
		steps = append(steps, step)

		// Actualizar base: reemplazar baseVars[leavingIndex] con enteringVar
		// Actualizar vector b
		// Calcular theta = minRatio
		theta := minRatio
		for i := range m {
			if i == leavingIndex {
				b.SetVec(i, theta)
				baseVars[i] = enteringVar
			} else {
				b.SetVec(i, b.At(i, 0)-theta*dVec.AtVec(i))
			}
		}

		iter++
	}
	warning := ""
	return 0, nil, steps, warning
}

func contains(s []int, e int) bool {
	return slices.Contains(s, e)
}

func matDenseToSlice(m *mat.Dense) []float64 {
	r, c := m.Dims()
	out := make([]float64, r*c)
	for i := range r {
		for j := range c {
			out[i*c+j] = m.At(i, j)
		}
	}
	return out
}

func matVecToSlice(v *mat.VecDense) []float64 {
	n := v.Len()
	out := make([]float64, n)
	for i := range n {
		out[i] = v.AtVec(i)
	}
	return out
}
