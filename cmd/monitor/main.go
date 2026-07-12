package main

import (
	"flag"
	"fmt"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
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

	var result checker.Result
	for _, v := range cfg.Sites {
		result = checker.CheckSite(v.URL)

		if result.Error != nil || !result.AvailabilityStatus {
			fmt.Printf("Site %s NOT ok\n", v.URL)
		} else {
			fmt.Printf("Site %s ok\n", v.URL)
		}
	}
}
