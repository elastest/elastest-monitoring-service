package main

import (
    "os"
	"testing"
    "fmt"
    "net/http"
)

func TestAll(t *testing.T) {

elastesturl := os.args[1]

client := &http.Client{}

// project creation
ftm.Println("Creating project")
req, err := http.NewRequest("POST", elastesturl + "/api/project", `{ "id": 0, "name": "proyecto" }`)
req.Header.Set("Content-Type", "application/json")

resp, err := http.Client.Do(req)
/*
# SuT creation
echo Creating SuT
SUT=$(curl -s -H "Content-Type: application/json" -d '{"project":'${PROJ}',"eimConfig":null,"parameters":[],"exProject":null,"id":0,"name":"nombredelsut","specification":"williamyeh/dummy","description":"descripciondelsut","sutType":"MANAGED","instrumentalize":false,"currentSutExec":1,"instrumentedBy":"WITHOUT","port":null,"managedDockerType":"IMAGE","mainService":"","commands":"","exTJobs":null,"commandsOption":"DEFAULT"}' "$ELASTESTURL/api/sut")
echo $SUT

# T-Job creation
echo Creating T-Job
DATA='{ "id": 0, "name": "nombredeltjob", "imageName": "docker.elastic.co/beats/metricbeat:5.4.0", "sut": '"${SUT}"', "project": '"${PROJ}"' , "tjobExecs": [], "parameters": [], "commands": "metricbeat -e -E output.logstash.hosts=[\"'"${EMSIPPORT}"'\"] -E output.elasticsearch.hosts=[\"edm-elasticsearch:9200\"]", '"$(cat suffix.json)"
TJOB=$(curl -s -H "Content-Type: application/json" -d "${DATA}" "$ELASTESTURL/api/tjob")
TJOBID=$(echo $TJOB | docker run -i --rm nimmis/jq .id)
echo TJOBID: $TJOBID

# T-Job execution
echo Executing T-Job
TJOBEXEC=$(curl -s -H "Content-Type: application/json" -d '{"tJobParams": []}' "$ELASTESTURL/api/tjob/$TJOBID/exec")
echo $TJOBEXEC
    resp, err := http.Get("http://example.com/")
    if err != nil {
    }
    fmt.Println(resp)
    */
}
