package gw2models

type GW2MaterialCategory struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Items []int  `json:"items"`
}
