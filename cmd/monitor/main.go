package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/scheduler"
)

func main() {
	fmt.Println("Site Monitor started. Press Ctrl+C to stop.")

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Site Monitor Error: %s\n", err)
		return
	}

	sh := scheduler.New(cfg.Sites, cfg.Interval)
	sh.Start()

	<-signals
	sh.Stop()

	fmt.Println("Site Monitor stopped.")
}
