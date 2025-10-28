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
}
