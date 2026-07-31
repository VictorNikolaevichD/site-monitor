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
	response.DurationMs = &durationMs

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

func ToCheckResultResponse(result domain.CheckResult) CheckResultResponse {
	checkedAt := result.CheckedAt
	durationMs := result.Duration.Milliseconds()

	response := CheckResultResponse{
		ID:           result.ID,
		Availability: result.Availability,
		CheckedAt:    &checkedAt,
		DurationMs:   &durationMs,
	}

	if result.Code != 0 {
		code := result.Code
		response.Code = &code
	}

	if result.Error != "" {
		errorMessage := result.Error
		response.Error = &errorMessage
	}

	return response
}

func ToCheckHistoryResponse(
	items []domain.CheckResult,
	total, limit, offset int,
) CheckHistoryResponse {
	result := make([]CheckResultResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ToCheckResultResponse(item))
	}

	return CheckHistoryResponse{
		Result: result,
		Pagination: PaginationMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}
}
