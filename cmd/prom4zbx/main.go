package main

import (
	"prom2zbx.com/internal/prom"
	"prom2zbx.com/internal/alerts"
)

func main() {
	// prom.GetTargetsProm2LLD()
	prom.GetRules()


	alerts.ListenAlerts()
	// _, err := http.NewRequest("POST", "http://127.0.0.1:10055", strings.NewReader(alertOK))
	// if err != nil {
	// 	fmt.Errorf("%v", err)
	// }
}
