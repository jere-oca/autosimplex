package pdf

import (
	"fmt"
	"io"

	m "github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// GenerateSimplexPDF escribe en w un PDF sencillo con el valor óptimo y la solución
func GenerateSimplexPDF(optimalValue float64, solution []float64, w io.Writer) error {
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

	// maroto.Output() devuelve (bytes.Buffer, error) en versiones recientes.
	// Obtenemos el buffer y lo copiamos al writer proporcionado.
	buf, err := mPdf.Output()
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}
