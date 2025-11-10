package pdf

import (
	"fmt"
	"io"
	"strings"

	"autosimplex/internal/simplex"

	m "github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// GenerateSimplexPDF escribe en w un PDF sencillo con el valor óptimo y la solución
// GenerateSimplexPDF escribe en w un PDF sencillo con el valor óptimo, la solución
// y las tablas intermedias (steps)
func GenerateSimplexPDF(optimalValue float64, solution []float64, steps []simplex.SimplexStep, w io.Writer) error {
	mPdf := pdf.NewMaroto(m.Portrait, m.A4)

	mPdf.Row(20, func() {
		mPdf.Col(12, func() {
			mPdf.Text("Resultado del Simplex", props.Text{Top: 2, Align: "center", Size: 16})
		})
	})

	mPdf.Row(10, func() {
		mPdf.Col(12, func() {
			mPdf.Text(fmt.Sprintf("Valor óptimo: %.6f", optimalValue), props.Text{Top: 2, Align: "left", Size: 12})
		})
	})

	// Lista de variables
	for i, v := range solution {
		idx := i + 1
		mPdf.Row(8, func() {
			mPdf.Col(12, func() {
				mPdf.Text(fmt.Sprintf("x%d = %.6f", idx, v), props.Text{Top: 2, Align: "left", Size: 11})
			})
		})
	}

	// Tablas intermedias (steps)
	if len(steps) > 0 {
		mPdf.Row(12, func() {
			mPdf.Col(12, func() {
				mPdf.Text("Tablas intermedias:", props.Text{Top: 2, Align: "left", Size: 14})
			})
		})

		for _, st := range steps {
			// Iteration header
			mPdf.Row(10, func() {
				mPdf.Col(12, func() {
					mPdf.Text(fmt.Sprintf("Iteración %d", st.Iteration), props.Text{Top: 2, Align: "left", Size: 12})
				})
			})

			// Render a compact tableau-like representation as monospaced formatted lines
			// Header: cj
			if len(st.Cj) > 0 {
				line := "cj:"
				for _, v := range st.Cj {
					line += fmt.Sprintf(" %8.2f", v)
				}
				mPdf.Row(8, func() {
					mPdf.Col(12, func() {
						mPdf.Text(line, props.Text{Top: 2, Align: "left", Size: 9})
					})
				})
			}

			// Table header
			totalCols := 0
			if len(st.Table) > 0 {
				totalCols = len(st.Table[0]) // includes RHS at last position
			}
			// cb | Base | vars... | R
			header := fmt.Sprintf("%6s | %8s |", "c_b", "Base")
			for j := 0; j < totalCols-1; j++ {
				header += fmt.Sprintf(" %8s", fmt.Sprintf("v%d", j+1))
			}
			header += fmt.Sprintf(" | %8s", "R")
			mPdf.Row(8, func() {
				mPdf.Col(12, func() {
					mPdf.Text(header, props.Text{Top: 2, Align: "left", Size: 9})
				})
			})

			// Rows
			for rIdx, row := range st.Table {
				cbVal := ""
				if rIdx < len(st.Cb) {
					cbVal = fmt.Sprintf("%6.2f", st.Cb[rIdx])
				}
				baseName := ""
				if rIdx < len(st.BaseVariables) {
					bv := st.BaseVariables[rIdx]
					if bv <= len(solution) {
						baseName = fmt.Sprintf("X%d", bv)
					} else {
						baseName = fmt.Sprintf("S%d", bv-len(solution))
					}
				}
				line := fmt.Sprintf("%6s | %8s |", cbVal, baseName)
				for c := 0; c < totalCols-1; c++ {
					val := 0.0
					if c < len(row) {
						val = row[c]
					}
					cell := fmt.Sprintf(" %8.2f", val)
					// mark pivot cell with surrounding brackets
					if rIdx == st.PivotRow && c == st.PivotCol {
						cell = fmt.Sprintf("[%6.2f] ", val)
					}
					line += cell
				}
				// RHS
				rhs := 0.0
				if totalCols-1 < len(row) {
					rhs = row[totalCols-1]
				}
				line += fmt.Sprintf(" | %8.2f", rhs)
				mPdf.Row(8, func() {
					mPdf.Col(12, func() {
						mPdf.Text(line, props.Text{Top: 2, Align: "left", Size: 9})
					})
				})
			}

			// Summary line: entering / leaving / t
			mPdf.Row(8, func() {
				mPdf.Col(12, func() {
					mPdf.Text(fmt.Sprintf("Entra: %d   Sale: %d   t: %.6f", st.EnteringVar, st.LeavingVar, st.TValue), props.Text{Top: 2, Align: "left", Size: 10})
				})
			})
		}
	}

	// maroto.Output() devuelve (bytes.Buffer, error) en versiones recientes.
	// Obtenemos el buffer y lo copiamos al writer proporcionado.
	buf, err := mPdf.Output()
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

func intsToString(a []int) string {
	if len(a) == 0 {
		return ""
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ", ")
}

func floatsToString(f []float64) string {
	if len(f) == 0 {
		return ""
	}
	parts := make([]string, len(f))
	for i, v := range f {
		parts[i] = fmt.Sprintf("%.6f", v)
	}
	return strings.Join(parts, ", ")
}
