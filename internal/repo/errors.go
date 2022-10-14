package repo

import "errors"

var (
	ErrDuplicateKey = errors.New("duplicate key")
	ErrNoFunds      = errors.New("no funds available")
)
