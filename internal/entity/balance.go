package entity

type Balance struct {
	ID        int     `json:"-" db:"id"`
	UserID    int     `json:"-" db:"user_id"`
	Current   float32 `json:"current" db:"current"`
	Withdrawn float32 `json:"withdrawn" db:"withdrawn"`
}
