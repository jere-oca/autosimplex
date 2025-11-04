package simplex

import (
	"math"
	"slices"

	"gonum.org/v1/gonum/mat"
)

// Solve is a convenience wrapper that assumes all constraints are "<=".
func Solve(maximize mat.Vector, constraints *mat.Dense) (float64, []float64) {
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
func SolveWithSigns(maximize mat.Vector, constraints *mat.Dense, signs []string) (float64, []float64) {
	const M = 1e7

	m, cols := constraints.Dims()
	n := maximize.Len()
	if cols != n+1 {
		// malformed matrix; return failure
		return 0, nil
	}

	// Count extra variables and build extended A matrix
	// We'll add slack (for <=), surplus+artificial (for >=), and artificial (for =)
	extraCols := 0
	artCols := 0
	for i := 0; i < m; i++ {
		s := "<="
		if i < len(signs) {
			s = signs[i]
		}
		switch s {
		case "<=":
			extraCols += 1 // slack
		case ">=":
			extraCols += 2 // surplus + artificial
			artCols += 1
		case "=":
			extraCols += 1 // artificial
			artCols += 1
		default:
			extraCols += 1
		}
	}

	totalVars := n + extraCols

	// Build A_extended
	A := mat.NewDense(m, totalVars, nil)
	// Fill original variables
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			A.Set(i, j, constraints.At(i, j))
		}
	}

	// Track indices and base variable per row
	col := n
	baseVars := make([]int, m) // 1-based indices of basic vars per row
	artIndices := []int{}
	for i := 0; i < m; i++ {
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
	for j := 0; j < n; j++ {
		c.Set(0, j, maximize.At(j, 0))
	}
	for _, ai := range artIndices {
		c.Set(0, ai, -M)
	}

	// Build b vector
	bData := make([]float64, m)
	for i := 0; i < m; i++ {
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
		for i := 0; i < m; i++ {
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
			return 0, nil
		}
		// y as row
		y := mat.NewDense(1, m, nil)
		for i := 0; i < m; i++ {
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
			for i := 0; i < m; i++ {
				bv := baseVars[i] - 1
				if bv < n {
					// original variable
					val := b.At(i, 0)
					solution[bv] = val
					optimal += c.At(0, bv) * val
				}
			}
			return optimal, solution
		}

		enteringVar := nonBase[entering]

		// Get column a for enteringVar
		raw := ATrans.RawRowView(enteringVar - 1)
		aVec := mat.NewVecDense(m, nil)
		for i := 0; i < m; i++ {
			aVec.SetVec(i, raw[i])
		}

		// Solve d = B^{-1} * aVec using LU
		dVec := mat.NewVecDense(m, nil)
		if err := lu.SolveVecTo(dVec, false, aVec); err != nil {
			return 0, nil
		}

		// Ratio test b_i / d_i for d_i > 0
		minRatio := math.Inf(1)
		leavingIndex := -1
		for i := 0; i < m; i++ {
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
			return 0, nil
		}

		// Update base: replace baseVars[leavingIndex] with enteringVar
		// Update b vector
		// Compute theta = minRatio
		theta := minRatio
		for i := 0; i < m; i++ {
			if i == leavingIndex {
				b.SetVec(i, theta)
				baseVars[i] = enteringVar
			} else {
				b.SetVec(i, b.At(i, 0)-theta*dVec.AtVec(i))
			}
		}

		iter++
	}

	return 0, nil
}

func contains(s []int, e int) bool {
	return slices.Contains(s, e)
}
