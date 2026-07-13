package dto

import (
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

func ToSiteResponse(site domain.Site) SiteResponse {
	return SiteResponse{
		ID:   site.ID,
		URL:  site.URL,
		Name: site.Name,
	}
}

func ToSiteResponses(sites []domain.Site) []SiteResponse {
	responses := make([]SiteResponse, 0, len(sites))
	for _, s := range sites {
		responses = append(responses, ToSiteResponse(s))
	}
	return responses
}
