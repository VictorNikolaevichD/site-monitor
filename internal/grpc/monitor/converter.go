package monitor

import (
	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
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
