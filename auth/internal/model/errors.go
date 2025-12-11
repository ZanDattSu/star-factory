package model

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrLoginAlreadyTaken = errors.New("login already taken")
	ErrEmailAlreadyTaken = errors.New("email already taken")
)

var (
	ErrInvalidUserFilter  = errors.New("invalid user filter")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSessionUUID = errors.New("invalid session uuid")
	ErrInvalidUserUUID    = errors.New("invalid user uuid")
	ErrInvalidTimestamp   = errors.New("invalid timestamp")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenGeneration    = errors.New("failed to generate token")
)
