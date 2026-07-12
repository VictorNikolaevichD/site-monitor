package scheduler

import (
	"fmt"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
)

type Scheduler struct {
	sites    []config.Site
	interval time.Duration
	isActive bool
	stop     chan struct{}
}

func New(sites []config.Site, interval time.Duration) *Scheduler {
	return &Scheduler{
		sites:    sites,
		interval: interval,
	}
}

func (s *Scheduler) Start() {
	if s.isActive {
		return
	}

	s.isActive = true
	s.schedule()
}

func (s *Scheduler) Stop() {
	if !s.isActive {
		return
	}

	s.stop <- struct{}{}
	close(s.stop)
	s.isActive = false
}

func (s *Scheduler) schedule() {
	ticker := time.NewTicker(s.interval)

	s.checkSites()
	for {
		select {
		case <-ticker.C:
			s.checkSites()
		case <-s.stop:
			ticker.Stop()
			return
		}
	}
}

func (s *Scheduler) checkSites() {
	var result checker.Result
	for _, v := range s.sites {
		t := time.Now()

		result = checker.CheckSite(v.URL)

		if result.Error != nil || !result.AvailabilityStatus {
			fmt.Printf("[%s] Site %s NOT ok\n", t.Format("2006-01-02 15:04:05"), v.URL)
		} else {
			fmt.Printf("[%s] Site %s ok\n", t.Format("2006-01-02 15:04:05"), v.URL)
		}
	}
}
