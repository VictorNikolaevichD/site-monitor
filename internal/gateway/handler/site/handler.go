package site

import (
	"encoding/json"
	"net/http"

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
)

type Handler struct {
	monitor monitorv1.MonitorServiceClient
}

func NewHandler(monitor monitorv1.MonitorServiceClient) *Handler {
	return &Handler{
		monitor: monitor,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp, err := h.monitor.GetSites(r.Context(), &monitorv1.GetSitesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sites := make([]map[string]string, 0, len(resp.GetSites()))
	for _, s := range resp.GetSites() {
		sites = append(sites, map[string]string{
			"id":   s.GetId(),
			"url":  s.GetUrl(),
			"name": s.GetName(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sites)
}
