package domain

import "errors"

// Domain errors.
var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrSyncFailed    = errors.New("sync failed")
	ErrDatabaseError = errors.New("database error")
	ErrExchangeError = errors.New("exchange server error")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrConflict      = errors.New("conflict")
	ErrInternalError = errors.New("internal error")
)
