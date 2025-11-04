package models

type Objective struct {
	N            int       `json:"n"`
	Coefficients []float64 `json:"coefficients"`
	// Type indicates whether to "maximize" or "minimize" the objective.
	// Optional: defaults to "maximize" when omitted.
	Type string `json:"type,omitempty"`
}

type Constraints struct {
	Rows int       `json:"rows"`
	Cols int       `json:"cols"`
	Vars []float64 `json:"vars"`
	// Signs optionally holds the sense of each constraint: "<=", ">=", or "=".
	// If omitted, all constraints are assumed to be "<=".
	Signs []string `json:"signs,omitempty"`
}

type SimplexRequest struct {
	Objective   Objective   `json:"objective"`
	Constraints Constraints `json:"constraints"`
}
