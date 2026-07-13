package site

import (
	"sync"

	"github.com/google/uuid"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type Repository struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]domain.Site
}

func NewRepository(sites []domain.Site) *Repository {
	storage := make(map[uuid.UUID]domain.Site, len(sites))

	for _, site := range sites {
		storage[site.ID] = site
	}

	return &Repository{
		storage: storage,
	}
}

func (r *Repository) GetAll() []domain.Site {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := make([]domain.Site, 0, len(r.storage))
	for _, rec := range r.storage {
		records = append(records, rec)
	}
	return records
}
