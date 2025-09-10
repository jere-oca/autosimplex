package models

type Objective struct {
	Type         string    `json:"type"` // "max" o "min"
	Coefficients []float64 `json:"coefficients"`
}

type Constraint struct {
	Coefficients []float64 `json:"coefficients"`
	Operator     string    `json:"operator"` // <=, >=, =
	RHS          float64   `json:"rhs"`
}

type SimplexRequest struct {
	Objective   Objective    `json:"objective"`
	Constraints []Constraint `json:"constraints"`
}
