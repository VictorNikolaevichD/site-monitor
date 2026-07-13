package site

import (
	"github.com/google/uuid"
)

type Site struct {
	ID   uuid.UUID
	URL  string
	Name string
}
