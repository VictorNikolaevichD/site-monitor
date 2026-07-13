package dto

import (
	"github.com/google/uuid"
)

type SiteResponse struct {
	ID   uuid.UUID `json:"id"`
	URL  string    `json:"url"`
	Name string    `json:"name"`
}
