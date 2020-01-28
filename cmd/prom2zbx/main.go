package main

import (
	"flag"
	"fmt"
	"regexp"

	"prom2zbx.com/internal/alerts"
	"prom2zbx.com/internal/prom"
	// "prom2zbx.com/internal/zbxsender"
)

func main() {
	if flag.NArgs()

	promURL := flag.String("promURL", "", "Prometheus URL")
	mode := flag.String("mode", "targets", "Mode: targets or rules")
	prefix := flag.String("prefix", "TEST", "Prefix for avoid duplicates")
	flag.Parse()
	switch *mode {
	case "targets":
		prom.GetTargetsProm2LLD(*promURL, *prefix)
	case "rules":
		prom.GetRules(*promURL)
	case "listen":
		alerts.ListenAlerts()
	case "test":
		test()
	}
}

func test() {
	testcase := []byte(`ABCDEFABDSDWDO5`)
	testcase2 := []byte(`ABCDEFABDSDWDO1.web.test.domain.com:9090`)
	testcase3 := []byte(`ABCDEFABDSDWDO2.web.test.domain.com:9090`)
	testcase4 := []byte(`ABCDEFABDSDWDP1.web.test.domain.com:9090`)
	testcase5 := []byte(`ABCDEFABDSDWDP2.web.test.domain.com:9090`)

	re := regexp.MustCompile("([A-Z0-9]+)")
	res := string(re.Find(testcase))
	res2 := string(re.Find(testcase2))
	res3 := string(re.Find(testcase3))
	res4 := string(re.Find(testcase4))
	res5 := string(re.Find(testcase5))

	fmt.Printf("test: %v\n", res[len(res)-7:])
	fmt.Printf("test: %v\n", res2[len(res2)-7:])
	fmt.Printf("test: %v\n", res3[len(res3)-7:])
	fmt.Printf("test: %v\n", res4[len(res4)-7:])
	fmt.Printf("test: %v\n", res5[len(res5)-7:])
}
