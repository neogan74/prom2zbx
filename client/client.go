package main

import (
	"log"
	"net/http"
	"strings"
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
	_, err := http.NewRequest("POST", "http://127.0.0.1:3000", strings.NewReader(alertOK))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
