package dto

import (
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SiteResponse struct {
	ID   uuid.UUID `json:"id"`
	URL  string    `json:"url"`
	Name string    `json:"name"`
}
