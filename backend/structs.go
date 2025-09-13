package main

type Objective struct {
	N            int       `json:"n"`
	Coefficients []float64 `json:"coefficients"`
}

type Constraints struct {
	Rows int       `json:"rows"`
	Cols int       `json:"cols"`
	Vars []float64 `json:"vars"`
}

type SimplexRequest struct {
	Objective   Objective   `json:"objective"`
	Constraints Constraints `json:"constraints"`
}
