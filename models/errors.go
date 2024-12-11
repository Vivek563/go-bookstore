package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")
var NoEnvFile = errors.New("models: no matching record found")
var ErrIncorrectPassword = errors.New("models: incorrect password")
var ErrUserNotFound = errors.New("models: no such user exist")
