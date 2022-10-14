package token

import "time"

// For future experiments with PaseTo

// IMaker is an interface to manage tokens
type IMaker interface {
	// CreateToken creates a new token for specific user_id and duration
	CreateToken(userID int, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
