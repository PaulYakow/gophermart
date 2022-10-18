package entity

import "time"

// UploadOrderDTO - object for API
type UploadOrderDTO struct {
	Number     string    `json:"number"`
	Status     string    `json:"status,omitempty"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
}

// UploadOrderDAO - object for database
type UploadOrderDAO struct {
	Number    string    `db:"number"`
	Status    string    `db:"status,omitempty"`
	Accrual   float32   `db:"accrual,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	UserID    int       `db:"user_id"`
}

// WithdrawOrderDTO - object for API
type WithdrawOrderDTO struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

// WithdrawOrderDAO - object for database
type WithdrawOrderDAO struct {
	Number    string    `db:"number"`
	Sum       float32   `db:"sum"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	UserID    int       `db:"user_id"`
}
