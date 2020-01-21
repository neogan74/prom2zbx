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
		targetid := target{v.Labels.Job + "." + targname}
		lldres.Res = append(lldres.Res, targetid)
	}
	final, _ := json.Marshal(lldres)
	fmt.Println(string(final))

}
