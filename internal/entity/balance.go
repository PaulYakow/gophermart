package entity

type Balance struct {
	User      string  `json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
