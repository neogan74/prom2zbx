package prom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//RulesStruct ...
type RulesStruct struct {
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

//Rule ...
type Rule struct {
	Name string `json:"{#RNAME}"`
}

//Rules ...
type Rules map[string][]Rule

//LLD ...
type LLD struct {
	Res []Rules `json:"data"`
}

//GetRules ...
func GetRules() {

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
	var rulesJS RulesStruct

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
	// fmt.Println(len(names))
	// for k := range names {
	// 	// fmt.Printf("Key: %v = %v len=%v\n\n", k, len(names[k]), names[k])
	// }
	// fmt.Println(RRules)
	final, _ := json.Marshal(RRules)
	// fmt.Println("keys %v %v", len(RRules))
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
