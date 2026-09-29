package monitor

import (
	"math"

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoSite(site domain.Site) *monitorv1.Site {
	return &monitorv1.Site{
		Id:   site.ID.String(),
		Url:  site.URL,
		Name: site.Name,
	}
}

func toProtoSites(sites []domain.Site) []*monitorv1.Site {
	out := make([]*monitorv1.Site, 0, len(sites))
	for _, s := range sites {
		out = append(out, toProtoSite(s))
	}
	return out
}

func toProtoSiteStatus(site domain.Site) *monitorv1.GetSiteStatusResponse {
	resp := &monitorv1.GetSiteStatusResponse{
		Url:          site.URL,
		Availability: string(domain.Unknown),
	}
	if site.LastCheck == nil {
		return resp
	}

	lc := site.LastCheck
	resp.Availability = string(lc.Availability)
	resp.CheckedAt = timestamppb.New(lc.CheckedAt)
	resp.DurationMs = lc.Duration.Milliseconds()
	if lc.Code != 0 {
		resp.Code = intToInt32(lc.Code)
	}
	resp.Error = lc.Error
	return resp
}

func toProtoCheckResult(result domain.CheckResult) *monitorv1.CheckResult {
	out := &monitorv1.CheckResult{
		Id:           result.ID,
		Availability: string(result.Availability),
		CheckedAt:    timestamppb.New(result.CheckedAt),
		DurationMs:   result.Duration.Milliseconds(),
		Error:        result.Error,
	}
	if result.Code != 0 {
		out.Code = intToInt32(result.Code)
	}
	return out
}

func toProtoCheckResults(items []domain.CheckResult) []*monitorv1.CheckResult {
	out := make([]*monitorv1.CheckResult, 0, len(items))
	for _, item := range items {
		out = append(out, toProtoCheckResult(item))
	}
	return out
}

func intToInt32(v int) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v)
}
