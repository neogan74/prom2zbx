package alerts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"prom2zbx.com/internal/zbxsender"
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

func ListenAlerts() {
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

	r.URL.Hostname()
	m1 := req.CommonLabels["severity"] + "[" + req.Alerts[0].Labels["alertname"] + "]"
	m2 := req.CommonLabels["severity"] + ".summary[" + req.Alerts[0].Labels["alertname"] + "]"
	fmt.Println()
	var data []*zbxsender.Metric
	data = append(data, zbxsender.NewMetric("Promth", m1, "1"))
	data = append(data, zbxsender.NewMetric("Promth", m2, req.Alerts[0].Annotations["summary"]))
	fmt.Println(data)
	pkg := zbxsender.NewPacket(data, time.Now().Unix())
	fmt.Println(pkg)

	zsnd := zbxsender.NewSender("51.15.213.9", 10144)
	res, err := zsnd.Send(pkg)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	fmt.Printf("Request %v\n\n", req)
	fmt.Printf("Request %v\n\n", req.Receiver)
	fmt.Printf("rem addr %v\n\n", r.RemoteAddr)
}
