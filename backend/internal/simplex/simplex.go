package simplex

import (
	"math"
	"slices"

	"gonum.org/v1/gonum/mat"
)

// Solve is a convenience wrapper that assumes all constraints are "<=".
func Solve(maximize mat.Vector, constraints *mat.Dense) (float64, []float64, []SimplexStep) {
	rows, _ := constraints.Dims()
	signs := make([]string, rows)
	for i := range signs {
		signs[i] = "<="
	}
	return SolveWithSigns(maximize, constraints, signs)
}

// SolveWithSigns solves a maximization LP given an objective vector and a
// constraint matrix (rows are [a1 ... an b]). 'signs' contains one of
// "<=", ">=", or "=" per constraint. Uses a Big-M strategy for artificials.
func SolveWithSigns(maximize mat.Vector, constraints *mat.Dense, signs []string) (float64, []float64, []SimplexStep) {
	const M = 1e7
	var steps []SimplexStep

	m, cols := constraints.Dims()
	n := maximize.Len()
	if cols != n+1 {
		// malformed matrix; return failure
		return 0, nil, steps
	}

	// Count extra variables and build extended A matrix
	// We'll add slack (for <=), surplus+artificial (for >=), and artificial (for =)
	extraCols := 0
	for i := range m {
		s := "<="
		if i < len(signs) {
			s = signs[i]
		}
		switch s {
		case "<=":
			extraCols += 1 // slack
		case ">=":
			extraCols += 2 // surplus + artificial
		case "=":
			extraCols += 1 // artificial
		default:
			extraCols += 1
		}
	}

	totalVars := n + extraCols

	// Build A_extended
	A := mat.NewDense(m, totalVars, nil)
	// Fill original variables
	for i := range m {
		for j := range n {
			A.Set(i, j, constraints.At(i, j))
		}
	}

	// Track indices and base variable per row
	col := n
	baseVars := make([]int, m) // 1-based indices of basic vars per row
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
			// surplus
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

	// Build objective c (1 x totalVars). Artificial variables get -M penalty
	c := mat.NewDense(1, totalVars, make([]float64, totalVars))
	for j := range n {
		c.Set(0, j, maximize.At(j, 0))
	}
	for _, ai := range artIndices {
		c.Set(0, ai, -M)
	}

	// Build b vector
	bData := make([]float64, m)
	for i := range m {
		bData[i] = constraints.At(i, n)
	}
	b := mat.NewVecDense(m, bData)

	// Now proceed with simplex iterations similar to prior implementation,
	// but using our constructed A, c, b and the provided baseVars.
	ATrans := mat.DenseCopyOf(A.T())

	const maxIter = 200
	iter := 0
	for {
		if iter > maxIter {
			break
		}

		// Build B from baseVars
		B := mat.NewDense(m, m, nil)
		cB := make([]float64, m)
		for i := range m {
			// baseVars stores 1-based index
			B.SetCol(i, ATrans.RawRowView(baseVars[i]-1))
			cB[i] = c.At(0, baseVars[i]-1)
		}

		// Solve B^T * y = cB  (y is column)
		var lu mat.LU
		lu.Factorize(B)
		yCol := mat.NewVecDense(m, nil)
		cBVec := mat.NewVecDense(m, cB)
		if err := lu.SolveVecTo(yCol, true, cBVec); err != nil {
			// singular base -> infeasible
			return 0, nil, steps
		}
		// y as row
		y := mat.NewDense(1, m, nil)
		for i := range m {
			y.Set(0, i, yCol.AtVec(i))
		}

		// Build AN and cN for non-basic vars
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

		// Reduced costs: cN - y * AN
		yAN := mat.NewDense(1, len(nonBase), nil)
		yAN.Mul(y, AN)

		// Choose entering variable: any index where cN > yAN
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
			// optimal
			// Build solution for original variables
			solution := make([]float64, n)
			var optimal float64
			for i := range m {
				bv := baseVars[i] - 1
				if bv < n {
					// original variable
					val := b.At(i, 0)
					solution[bv] = val
					optimal += c.At(0, bv) * val
				}
			}
			return optimal, solution, steps
		}

		enteringVar := nonBase[entering]

		// Get column a for enteringVar
		raw := ATrans.RawRowView(enteringVar - 1)
		aVec := mat.NewVecDense(m, nil)
		for i := range m {
			aVec.SetVec(i, raw[i])
		}

		// Solve d = B^{-1} * aVec using LU
		dVec := mat.NewVecDense(m, nil)
		if err := lu.SolveVecTo(dVec, false, aVec); err != nil {
			return 0, nil, steps
		}

		// Ratio test b_i / d_i for d_i > 0
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
			// unbounded
			return 0, nil, steps
		}

		// Prepare step
		currentBaseVars := append([]int{}, baseVars...)
		currentNonBaseVars := append([]int{}, nonBase...)
		reducedCosts := matDenseToSlice(yAN)
		bVector := matVecToSlice(b)
		newBaseVar := enteringVar
		leavingVar := baseVars[leavingIndex]
		lowestValueOfT := minRatio

		step := SimplexStep{
			Iteration:        iter,
			BaseVariables:    currentBaseVars,
			NonBaseVariables: currentNonBaseVars,
			ReducedCosts:     reducedCosts,
			BVector:          bVector,
			EnteringVar:      newBaseVar,
			LeavingVar:       leavingVar,
			TValue:           lowestValueOfT,
		}
		steps = append(steps, step)

		// Update base: replace baseVars[leavingIndex] with enteringVar
		// Update b vector
		// Compute theta = minRatio
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

	return 0, nil, steps
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
