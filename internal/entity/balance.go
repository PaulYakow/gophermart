package entity

// BalanceDTO - object for API
type BalanceDTO struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

// BalanceDAO - object for database
type BalanceDAO struct {
	UserID    int     `db:"user_id"`
	Current   float32 `db:"current"`
	Withdrawn float32 `db:"withdrawn"`
}
