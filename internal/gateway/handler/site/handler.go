package site

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		code, msg := mapGRPCError(err)
		writeError(w, code, msg)
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	resp, err := h.monitor.CreateSite(r.Context(), &monitorv1.CreateSiteRequest{
		Url:  body.URL,
		Name: body.Name,
	})
	if err != nil {
		code, msg := mapGRPCError(err)
		writeError(w, code, msg)
		return
	}

	site := resp.GetSite()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":   site.GetId(),
		"url":  site.GetUrl(),
		"name": site.GetName(),
	})
}

func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := h.monitor.DeleteSite(r.Context(), &monitorv1.DeleteSiteRequest{Id: id})
	if err != nil {
		code, msg := mapGRPCError(err)
		writeError(w, code, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetStatusByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	resp, err := h.monitor.GetSiteStatus(r.Context(), &monitorv1.GetSiteStatusRequest{Id: id})
	if err != nil {
		code, msg := mapGRPCError(err)
		writeError(w, code, msg)
		return
	}

	out := map[string]any{
		"url":          resp.GetUrl(),
		"availability": resp.GetAvailability(),
	}
	if resp.GetCode() != 0 {
		out["code"] = resp.GetCode()
	}
	if resp.GetCheckedAt() != nil {
		out["checked_at"] = resp.GetCheckedAt().AsTime()
	}
	if resp.GetDurationMs() != 0 {
		out["duration_ms"] = resp.GetDurationMs()
	}
	if resp.GetError() != "" {
		out["error"] = resp.GetError()
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *Handler) GetHistoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var limit, offset int32
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			writeError(w, http.StatusBadRequest, "некорректный параметр limit")
			return
		}
		limit = int32(n)
	}
	if raw := r.URL.Query().Get("offset"); raw != "" {
		n, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			writeError(w, http.StatusBadRequest, "некорректный параметр offset")
			return
		}
		offset = int32(n)
	}

	resp, err := h.monitor.GetSiteHistory(r.Context(), &monitorv1.GetSiteHistoryRequest{
		Id:     id,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		code, msg := mapGRPCError(err)
		writeError(w, code, msg)
		return
	}

	results := make([]map[string]any, 0, len(resp.GetResults()))
	for _, item := range resp.GetResults() {
		row := map[string]any{
			"id":           item.GetId(),
			"availability": item.GetAvailability(),
		}
		if item.GetCode() != 0 {
			row["code"] = item.GetCode()
		}
		if item.GetCheckedAt() != nil {
			row["checked_at"] = item.GetCheckedAt().AsTime()
		}
		if item.GetDurationMs() != 0 {
			row["duration_ms"] = item.GetDurationMs()
		}
		if item.GetError() != "" {
			row["error"] = item.GetError()
		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"result": results,
		"pagination": map[string]int32{
			"total":  resp.GetTotal(),
			"limit":  resp.GetLimit(),
			"offset": resp.GetOffset(),
		},
	})
}
