package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// GetTargetsProm2LLD()
	getRules()
}

//GetTargetsProm2LLD ...
func GetTargetsProm2LLD() {

	type TargetsList struct {
		Status string `json:"status"`
		Data   struct {
			ActiveTargets []struct {
				DiscoveredLabels struct {
					Address     string `json:"__address__"`
					MetricsPath string `json:"__metrics_path__"`
					Scheme      string `json:"__scheme__"`
					Group       string `json:"group"`
					Job         string `json:"job"`
				} `json:"discoveredLabels"`
				Labels struct {
					Group    string `json:"group"`
					Instance string `json:"instance"`
					Job      string `json:"job"`
				} `json:"labels"`
				ScrapeURL  string    `json:"scrapeUrl"`
				LastError  string    `json:"lastError"`
				LastScrape time.Time `json:"lastScrape"`
				Health     string    `json:"health"`
			} `json:"activeTargets"`
			DroppedTargets []interface{} `json:"droppedTargets"`
		} `json:"data"`
	}
	resp, err := http.Get("http://192.168.33.11/targets.html")
	if err != nil {
		fmt.Errorf("Error while get targets: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error with body: %v\n", err)
	}
	var targets []string
	var targetsjs TargetsList
	err = json.Unmarshal(data, &targetsjs)
	if err != nil {
		fmt.Errorf("Error while JSON unmarshal %s", err)
	}

	type target struct {
		Name string `json:"{#TRGNAME}"`
		ZTag string `json:"{#JOB}"`
	}

	type LLD struct {
		Res []target `json:"data"`
	}

	var lldres LLD
	for _, v := range targetsjs.Data.ActiveTargets {
		tmp := v.Labels.Instance[:strings.LastIndex(v.Labels.Instance, ":")]
		if strings.Contains(tmp, "https") {
			tmp = strings.TrimLeft(tmp, "https://")
		}
		if len(tmp) == 0 {
			continue
		}
		// fmt.Println(tmp)
		targname := tmp
		targets = append(targets, v.Labels.Job+"."+targname)
		targetid := target{Name: v.Labels.Job + "." + targname,
			ZTag: v.Labels.Job}
		lldres.Res = append(lldres.Res, targetid)
	}
	final, _ := json.Marshal(lldres)
	fmt.Println(string(final))

}

func getRules() {

	type PromRules struct {
		Status string `json:"status"`
		Data   struct {
			Groups []struct {
				Name  string `json:"name"`
				File  string `json:"file"`
				Rules []struct {
					Name     string `json:"name"`
					Query    string `json:"query"`
					Duration string `json:"duration"`
					Labels   struct {
						Receiver string `json:"receiver"`
						Severity string `json:"severity"`
					} `json:"labels"`
					Annotations struct {
						Description string `json:"description"`
						Summary     string `json:"summary"`
					} `json:"annotations"`
					Alerts []struct {
						Alert string `json:"alert"`
					} `json:"alerts"`
					Health string `json:"health"`
					Type   string `json:"type"`
				} `json:"rules"`
				Interval int `json:"interval"`
			} `json:"groups"`
		} `json:"data"`
	}
	type Rule struct {
		Name string `json:"{#RNAME}"`
	}

	type Rules map[string][]Rule

	type LLD struct {
		Res []Rules `json:"data"`
	}

	resp, err := http.Get("http://192.168.33.11/rules.html")
	if err != nil {
		fmt.Errorf("Error while get targets: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error with body: %v\n", err)
	}

	// var rules []string{}
	var rulesJS PromRules

	err = json.Unmarshal(data, &rulesJS)
	if err != nil {
		fmt.Errorf("Error while JSON unmarshal %s", err)
	}

	names := make(map[string][]string)
	RRules := map[string][]Rule{}

	for _, v := range rulesJS.Data.Groups {
		// fmt.Println(i, v.Name)

		for _, vv := range v.Rules {
			// fmt.Printf("\t\t= %v - %v + %v\n", j, vv.Name, vv.Labels.Severity)
			if contains(names[vv.Labels.Severity], vv.Name) {
				continue
			}
			RRules[vv.Labels.Severity] = append(RRules[vv.Labels.Severity], Rule{Name: vv.Name})
			names[vv.Labels.Severity] = append(names[vv.Labels.Severity], vv.Name)
		}
	}
	// for _ := range names {
	// 	// fmt.Printf("Key: %v = %v\n\n", k, names[k])
	// }
	// fmt.Println(RRules)
	final, _ := json.Marshal(RRules)
	// fmt.Println()
	// fmt.Println()
	fmt.Println(string(final))

}

//ADDITIONAL FUNCTIONS

func contains(sl []string, str string) bool {
	for _, a := range sl {
		if a == str {
			return true
		}
	}
	return false
}
