package service

import "errors"

var (
	ErrLoginNotExist    = errors.New("login not exist")
	ErrMismatchPassword = errors.New("password mismatch")
)
