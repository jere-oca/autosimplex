package pdf

import (
	"fmt"
	"io"

	"autosimplex/internal/simplex"

	"github.com/johnfercher/maroto/pkg/color"
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
			// Iteracción
			mPdf.Row(10, func() {
				mPdf.Col(12, func() {
					mPdf.Text(fmt.Sprintf("Iteración %d", st.Iteration), props.Text{Top: 2, Align: "left", Size: 12})
				})
			})

			// Renderizar una representación compacta tipo tableau con líneas formateadas monoespacio
			// Encabezado: cj
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

			// Encabezado de tabla
			// Intentar renderizar una tabla usando Maroto TableList para un estilo visual más agradable
			headers := []string{"c_b", "Base"}
			if len(st.Table) > 0 {
				for j := 0; j < len(st.Table[0])-1; j++ {
					headers = append(headers, fmt.Sprintf("v%d", j+1))
				}
			}
			headers = append(headers, "R")

			contents := [][]string{}
			for rIdx, row := range st.Table {
				cols := []string{"", ""}
				if rIdx < len(st.Cb) {
					cols[0] = fmt.Sprintf("%.2f", st.Cb[rIdx])
				}
				if rIdx < len(st.BaseVariables) {
					bv := st.BaseVariables[rIdx]
					if bv <= len(solution) {
						cols[1] = fmt.Sprintf("X%d", bv)
					} else {
						cols[1] = fmt.Sprintf("S%d", bv-len(solution))
					}
				}
				for c := 0; c < len(row)-1; c++ {
					val := row[c]
					cell := fmt.Sprintf("%.2f", val)
					if rIdx == st.PivotRow && c == st.PivotCol {
						cell = "▶" + cell + "◀"
					}
					cols = append(cols, cell)
				}
				// Lado derecho (RHS)
				if len(row) > 0 {
					cols = append(cols, fmt.Sprintf("%.2f", row[len(row)-1]))
				} else {
					cols = append(cols, "")
				}
				contents = append(contents, cols)
			}

			// Renderizar tabla manualmente para permitir fondo de celda pivote
			// Construir tamaños de cuadrícula (suma a 12)
			nCols := len(headers)
			gridTotal := 12
			gridSizes := make([]uint, nCols)
			if nCols == 1 {
				gridSizes[0] = uint(gridTotal)
			} else {
				// dar a las primeras dos columnas ancho más pequeño, distribuir el resto
				first := 2
				second := 2
				remaining := gridTotal - first - second
				per := 1
				if nCols-2 > 0 {
					per = remaining / (nCols - 2)
				}
				for i := 0; i < nCols; i++ {
					if i == 0 {
						gridSizes[i] = uint(first)
					} else if i == 1 {
						gridSizes[i] = uint(second)
					} else if i == nCols-1 {
						// última columna toma el resto
						sum := 0
						for j := 0; j < nCols-1; j++ {
							sum += int(gridSizes[j])
						}
						gridSizes[i] = uint(gridTotal - sum)
					} else {
						gridSizes[i] = uint(per)
					}
				}
			}

			// Fila de encabezado
			headerHeight := 8.0
			headerBg := props.Text{Top: 2, Align: "CENTER", Size: 10}
			// habilitar bordes (dibujará líneas negras)
			mPdf.SetBorder(true)
			// dibujar encabezado con fondo
			mPdf.Row(headerHeight, func() {
				for i, h := range headers {
					gs := gridSizes[i]
					mPdf.Col(gs, func() {
						// fondo del encabezado
						mPdf.SetBackgroundColor(color.Color{Red: 14, Green: 165, Blue: 233})
						mPdf.Text(h, headerBg)
						// restablecer fondo
						mPdf.SetBackgroundColor(color.NewWhite())
					})
				}
			})

			// Filas de contenido
			contentHeight := 7.0
			for rIdx, row := range contents {
				mPdf.Row(contentHeight, func() {
					for cIdx, cell := range row {
						gs := gridSizes[cIdx]
						mPdf.Col(gs, func() {
							// si esta es la celda pivote, dibujar fondo coloreado
							if rIdx == st.PivotRow && cIdx == st.PivotCol+2 { // +2 porque contents incluye columnas cb y Base
								// fondo de burbuja pivote (azul)
								mPdf.SetBackgroundColor(color.Color{Red: 59, Green: 130, Blue: 246})
								// escribir celda en blanco
								mPdf.Text(cell, props.Text{Top: 1, Align: "CENTER", Size: 9, Color: color.Color{Red: 255, Green: 255, Blue: 255}})
								mPdf.SetBackgroundColor(color.NewWhite())
							} else {
								mPdf.Text(cell, props.Text{Top: 1, Align: "CENTER", Size: 9})
							}
						})
					}
				})
			}
			// deshabilitar bordes después de dibujar tabla
			mPdf.SetBorder(false)

			// Línea de resumen: entrante / saliente / t
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

// funciones auxiliares eliminadas (no usadas)
