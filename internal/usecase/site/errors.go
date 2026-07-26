package site

import (
	"errors"
)

var (
	ErrSiteAlreadyExists = errors.New("site already exists")
	ErrInvalidLimit      = errors.New("invalid limit")
	ErrInvalidOffset     = errors.New("invalid offset")
)
