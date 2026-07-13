package site

import (
	"errors"
)

var (
	ErrSiteAlreadyExists = errors.New("сайт с таким URL уже существует")
	ErrSiteNotFound      = errors.New("Сайт с таким ID не найден")
)
