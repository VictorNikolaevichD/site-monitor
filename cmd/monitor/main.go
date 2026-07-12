package main

import (
	"fmt"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
)

var sites = [10]string{
	"https://vk.com",
	"https://vk.ru",
	"https://viniapps.ru",
	"https://viniapps.com",
	"https://google.com",
	"https://yandex.ru",
	"https://habr.com",
	"https://youtube.com",
	"https://abrakadabra.ru",
	"https://kremlin.ru",
}

func main() {
	fmt.Println("Site Monitor started")

	var result checker.Result
	for _, v := range sites {
		result = checker.CheckSite(v)
		
		if result.Error != nil || !result.AvailabilityStatus {
			fmt.Printf("Site %s NOT ok\n", v)
		} else {
			fmt.Printf("Site %s ok\n", v)
		}
	}
}
