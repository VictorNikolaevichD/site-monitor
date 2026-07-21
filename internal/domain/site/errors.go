package site

import (
	"errors"
)

var (
	ErrInvalidURL   = errors.New("invalid site URL")
	ErrURLRequired  = errors.New("site URL is required")
	ErrNameRequired = errors.New("site name is required")

	ErrSiteNotFound = errors.New("site not found")
	ErrStorage      = errors.New("storage error")
)
