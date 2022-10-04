package entity

import "time"

type UploadOrder struct {
	ID         int       `json:"-" db:"id"`
	UserID     int       `json:"-" db:"user_id"`
	Number     int       `json:"number" db:"number"`
	Status     string    `json:"status" db:"status"`
	Accrual    float32   `json:"accrual" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"created_at"`
}

type WithdrawOrder struct {
	ID          int       `json:"-" db:"id"`
	UserID      int       `json:"-" db:"user_id"`
	Order       int       `json:"order" db:"number"`
	Sum         float32   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"created_at"`
}
