package site

import (
	"errors"
)

var (
	ErrInvalidURL   = errors.New("некорректные URL сайта")
	ErrURLRequired  = errors.New("URL сайта обязателен")
	ErrNameRequired = errors.New("имя сайта обязательно")
)
