package main

import (
	"flag"

	"prom2zbx.com/internal/prom"
)

func main() {
	promURL := flag.String("promURL", "", "Prometheus URL")
	mode := flag.String("mode", "targets", "Mode: targets or rules")
	prefix := flag.String("prefix", "TEST", "Prefix for avoid duplicates")
	flag.Parse()
	switch *mode {
	case "targets":
		prom.GetTargetsProm2LLD(*promURL, *prefix)
	case "rules":
		prom.GetRules(*promURL)
	}
	// alerts.ListenAlerts()
	// _, err := http.NewRequest("POST", "http://127.0.0.1:10055", strings.NewReader(alertOK))
	// if err != nil {
	// 	fmt.Errorf("%v", err)
	// }
}
