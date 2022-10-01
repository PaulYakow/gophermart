package entity

import "time"

type Order struct {
	Id         int       `json:"-"`
	User       string    `json:"-"`
	Number     int       `json:"number" db:"number"`
	Status     string    `json:"status" db:"status"`
	Accrual    float64   `json:"accrual" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
