package errors

import (
	"errors"
)

var (
	ErrInvalidURL        = errors.New("некорректные URL сайта")
	ErrURLRequired       = errors.New("URL сайта обязателен")
	ErrNameRequired      = errors.New("имя сайта обязательно")
	ErrSiteAlreadyExists = errors.New("сайт с таким URL уже существует")
)
