package site

import (
	"sync"

	"github.com/google/uuid"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
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

func (r *Repository) AddIfAbsent(site domain.Site) (domain.Site, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.storage {
		if s.URL == site.URL {
			return domain.Site{}, siteusecase.ErrSiteAlreadyExists
		}
	}

	r.storage[site.ID] = site
	return site, nil
}

func (r *Repository) DeleteByID(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[id]; !exists {
		return domain.ErrSiteNotFound
	}

	delete(r.storage, id)
	return nil
}

func (r *Repository) UpdateLastCheckByID(id uuid.UUID, checkStatus domain.CheckStatus) (domain.Site, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	site, exists := r.storage[id]
	if !exists {
		return domain.Site{}, domain.ErrSiteNotFound
	}

	site.LastCheck = &checkStatus
	r.storage[id] = site

	return site, nil
}

func (r *Repository) GetByID(id uuid.UUID) (domain.Site, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	site, exists := r.storage[id]
	if !exists {
		return domain.Site{}, domain.ErrSiteNotFound
	}

	return site, nil
}
