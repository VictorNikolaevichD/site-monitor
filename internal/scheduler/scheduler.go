package scheduler

import (
	"log/slog"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
)

type Scheduler struct {
	sites    []config.Site
	interval time.Duration
	stop     chan struct{}
	done     chan struct{}
	logger   *slog.Logger
}

func New(sites []config.Site, interval time.Duration, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		sites:    sites,
		interval: interval,
		logger:   logger,
	}
}

func (s *Scheduler) Start() {
	if s.stop != nil || s.done != nil {
		return
	}

	s.stop = make(chan struct{})
	s.done = make(chan struct{})

	go s.schedule()
}

func (s *Scheduler) Stop() {
	s.logger.Info("shutting down...")

	if s.stop == nil || s.done == nil {
		return
	}

	close(s.stop)
	<-s.done
}

func (s *Scheduler) schedule() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	defer close(s.done)

	s.checkSites()
	for {
		select {
		case <-ticker.C:
			s.checkSites()
		case <-s.stop:
			return
		}
	}
}

func (s *Scheduler) checkSites() {
	var result checker.Result
	for _, v := range s.sites {
		result = checker.CheckSite(v.URL)

		if result.Error != nil {
			s.logger.Error("site check failed", "status", "NOT ok", "url", v.URL, "error", result.Error)
			continue
		}

		if !result.AvailabilityStatus {
			s.logger.Warn("site unavailable", "status", "NOT ok", "code", result.Code, "url", v.URL)
			continue
		}
		
		s.logger.Info("site available", "status", "ok", "code", result.Code, "url", v.URL)
	}
}
