package auth

import "errors"

var (
	// ErrInvalidToken 無效的 token
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidID 無效的 ID
	ErrInvalidID = errors.New("invalid id")
	// ErrExpiredToken 過期的 token
	ErrExpiredToken = errors.New("token has expired")
	// ErrUnauthorized 沒有 Authorization Header
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden Forbidden
	ErrForbidden = errors.New("forbidden")
)
