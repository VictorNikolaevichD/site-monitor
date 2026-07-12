package main

import (
	"flag"
	"fmt"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/scheduler"
)

func main() {
	fmt.Println("Site Monitor started")

	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Site Monitor Error: %s\n", err)
		return
	}

	sh := scheduler.New(cfg.Sites, cfg.Interval)
	sh.Start()
}
