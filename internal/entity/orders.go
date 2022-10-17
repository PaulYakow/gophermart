package entity

import "time"

//todo: Использовать разные объекты для JSON и БД (DTO?)

type UploadOrder struct {
	ID         int       `json:"-" db:"id"`
	UserID     int       `json:"-" db:"user_id"`
	Number     string    `json:"number" db:"number"`
	Status     string    `json:"status,omitempty" db:"status,omitempty"`
	Accrual    float32   `json:"accrual,omitempty" db:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty" db:"created_at,omitempty"`
}

type WithdrawOrder struct {
	ID          int       `json:"-" db:"id"`
	UserID      int       `json:"-" db:"user_id"`
	Order       string    `json:"order" db:"number"`
	Sum         float32   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"created_at"`
}
