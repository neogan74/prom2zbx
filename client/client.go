package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {

	url := "http://localhost:9094"
	method := "POST"

	payload := strings.NewReader("{\n    \"version\": \"4\",\n    \"groupKey\": \"{")
	/* }:{alertname=\\\"InstanceDown\\\"}\",\n
	\"status\": \"resolved\",\n
	\"receiver\": \"testing\",\n
	\"groupLabels\": {\n
	\"alertname\": \"InstanceDown\"\n
	},\n
	\"commonLabels\": {\n
		\"alertname\": \"InstanceDown\",\n
		\"instance\": \"abc\",\n
		\"job\": \"node_exporter\",\n
		\"severity\": \"critical\"\n
	},\n
	\"commonAnnotations\": {\n
			\"description\": \"localhost:9100 of job node_exporter has been down for more than 1 minute.\",\n
			\"summary\": \"Instance localhost:9100 down\"\n
	},\n
	\"externalURL\": \"http://edas-GE72-6QC:9093\",\n
	\"alerts\": [\n
		{\n
			\"labels\": {\n
			\"alertname\": \"InstanceDown\",\n
			\"instance\": \"ABCDEFABDSDWDP02.web.abc02.domain.com\",\n
			\"job\": \"node_exporter\",\n
			\"severity\": \"critical\"\n
		},\n
		\"annotations\": {\n
			\"description\": \"localhost:9100 of job node_exporter has been down for more than 1 minute.\",\n
			\"summary\": \"Instance localhost:9100 down\"\n
		},\n
		\"startsAt\": \"2018-08-30T16:59:09.653872838+03:00\",\n
		\"EndsAt\": \"2018-08-30T17:01:09.656110177+03:00\"\n
		}\n
	]\n}") */

	time.Sleep(time.Second * 10)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
