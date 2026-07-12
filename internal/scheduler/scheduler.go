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
	stop     chan struct{}
	done     chan struct{}
}

func New(sites []config.Site, interval time.Duration) *Scheduler {
	return &Scheduler{
		sites:    sites,
		interval: interval,
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
	fmt.Println("Shutting down...")

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
		t := time.Now()

		result = checker.CheckSite(v.URL)

		if result.Error != nil || !result.AvailabilityStatus {
			fmt.Printf("[%s] Site %s NOT ok\n", t.Format("2006-01-02 15:04:05"), v.URL)
		} else {
			fmt.Printf("[%s] Site %s ok\n", t.Format("2006-01-02 15:04:05"), v.URL)
		}
	}
}
