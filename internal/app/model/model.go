package model

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no suitable entry was found")
	//If a user tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	//If a user tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
