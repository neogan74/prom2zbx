package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const alertOK = `{  
	"version":"4",
	"groupKey":"{}:{alertname=\"InstanceDown\"}",
	"status":"resolved",
	"receiver":"testing",
	"groupLabels":{  
	   "alertname":"InstanceDown"
	},
	"commonLabels":{  
	   "alertname":"InstanceDown",
	   "instance":"abc",
	   "job":"node_exporter",
	   "severity":"critical"
	},
	"commonAnnotations":{  
	   "description":"localhost:9100 of job node_exporter has been down for more than 1 minute.",
	   "summary":"Instance localhost:9100 down"
	},
	"externalURL":"http://edas-GE72-6QC:9093",
	"alerts":[  
	   {  
		  "labels":{  
			 "alertname":"InstanceDown",
			 "instance":"localhost:9100",
			 "job":"node_exporter",
			 "severity":"critical"
		  },
		  "annotations":{  
			 "description":"localhost:9100 of job node_exporter has been down for more than 1 minute.",
			 "summary":"Instance localhost:9100 down"
		  },
		  "startsAt":"2018-08-30T16:59:09.653872838+03:00",
		  "EndsAt":"2018-08-30T17:01:09.656110177+03:00"
	   }
	]
 }`

func main() {
	// getTargetsProm2LLD()
	// getRules()
	listenAlerts()
	_, err := http.NewRequest("POST", "http://127.0.0.1:10055", strings.NewReader(alertOK))
	if err != nil {
		fmt.Errorf("%v", err)
	}

}

//getTargetsProm2LLD ...
func getTargetsProm2LLD() {

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

	for i, v := range rulesJS.Data.Groups {
		fmt.Println(i, v.Name)

		for j, vv := range v.Rules {
			fmt.Printf("\t\t= %v - %v + %v\n", j, vv.Name, vv.Labels.Severity)
			if contains(names[vv.Labels.Severity], vv.Name) {
				continue
			}
			RRules[vv.Labels.Severity] = append(RRules[vv.Labels.Severity], Rule{Name: vv.Name})
			names[vv.Labels.Severity] = append(names[vv.Labels.Severity], vv.Name)
		}
	}
	// fmt.Println(len(names))
	for k := range names {
		fmt.Printf("Key: %v = %v len=%v\n\n", k, len(names[k]), names[k])
	}
	// fmt.Println(RRules)
	final, _ := json.Marshal(RRules)
	fmt.Println("keys %v %v", len(RRules))
	fmt.Println(string(final))

}

func listenAlerts() {
	// type JSONHandler struct {
	// 	// Sender      *zabbixsnd.Sender
	// 	KeyPrefix   string
	// 	DefaultHost string
	// 	Hosts       map[string]string
	// }
	http.HandleFunc("/", HandlePost)

	log.Println("Listening on localhost:3000")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))

}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	// AlertManager
	type Alert struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
		StartAt     map[string]string `json:"startAt,omitempty"`
		EndAt       map[string]string `json:"endAt,omitempty"`
	}
	type AlertManagerRequest struct {
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		Status            string            `json:"status"`
		Receiver          string            `json:"receiver"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Alerts            []Alert           `json:"alerts"`
	}

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var req AlertManagerRequest
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "request body is not valid json", http.StatusBadRequest)
		return
	}

	if req.Status == "" || req.CommonLabels["alertname"] == "" {
		http.Error(w, "missing fields in request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Request %v\n\n", req)
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
