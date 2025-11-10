package simplex

type SimplexStep struct {
	Iteration        int       `json:"iteration"`
	BaseVariables    []int     `json:"base_variables"`
	NonBaseVariables []int     `json:"non_base_variables"`
	ReducedCosts     []float64 `json:"reduced_costs"`
	BVector          []float64 `json:"b_vector"`
	EnteringVar      int       `json:"entering_var"`
	LeavingVar       int       `json:"leaving_var"`
	TValue           float64   `json:"t_value"`

	// Tabla completa por iteración: filas = restricciones, columnas = variables extendidas, última columna = R (b)
	Table [][]float64 `json:"table,omitempty"`
	// Cj: cabecera de coeficientes objetivo (uno por variable extendida)
	Cj []float64 `json:"cj,omitempty"`
	// Cb: coeficientes objetivo de las variables básicas (uno por fila)
	Cb []float64 `json:"cb,omitempty"`
	// Pivot position: fila (0-based) y columna (0-based dentro de variables extendidas)
	PivotRow int `json:"pivot_row,omitempty"`
	PivotCol int `json:"pivot_col,omitempty"`
}
