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

func ToStatusResponse(site domain.Site) StatusResponse {
	response := StatusResponse{
		URL:          site.URL,
		Availability: domain.Unknown,
	}

	if site.LastCheck == nil {
		return response
	}

	lastCheck := site.LastCheck

	response.Availability = lastCheck.Availability

	checkedAt := lastCheck.CheckedAt
	response.CheckedAt = &checkedAt

	durationMs := lastCheck.Duration.Milliseconds()
	response.Duration = &durationMs

	if lastCheck.Code != 0 {
		code := lastCheck.Code
		response.Code = &code
	}

	if lastCheck.Error != "" {
		errorMessage := lastCheck.Error
		response.Error = &errorMessage
	}

	return response
}
