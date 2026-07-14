package scheduler

import (
	"log/slog"
	"time"
)

type checkSitesUseCase interface {
	Execute()
}

type Scheduler struct {
	checkSitesUseCase checkSitesUseCase
	interval          time.Duration
	logger            *slog.Logger

	stop chan struct{}
	done chan struct{}
}

func New(checkSitesUseCase checkSitesUseCase, interval time.Duration, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		checkSitesUseCase: checkSitesUseCase,
		interval:          interval,
		logger:            logger,
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

	s.checkSitesUseCase.Execute()
	for {
		select {
		case <-ticker.C:
			s.checkSitesUseCase.Execute()
		case <-s.stop:
			return
		}
	}
}
