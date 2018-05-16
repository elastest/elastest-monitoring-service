#!/bin/bash

ELASTESTURL=$1

# Posting EMS instance
INSTANCEID=$(curl -s -XPOST $ELASTESTURL/api/esm/services/bab3ae67-8c1d-46ec-a940-94183a443825/prov)

EMSURL=null

while [ "X$EMSURL" = "Xnull" ];
do 
    echo Waiting and retrying..
    sleep 1
    EMSURL=$(curl -s $ELASTESTURL/api/esm/services/instances/$INSTANCEID | docker run -i --rm nimmis/jq .urls.api)
done

EMSURL=$(echo "$EMSURL" | tr -d '"')
echo EMS URL: $EMSURL
EMSIPPORT=$(echo "$EMSURL" | sed 's/http:..//;s/:.*//')":5044"
echo EMS IP AND PORT: $EMSIPPORT

# project creation
echo Creating project
PROJ=$(curl -s -H "Content-Type: application/json" -d '{ "id": 0, "name": "proyecto" }' "$ELASTESTURL/api/project")
echo $PROJ

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

# Waiting for the T-Job to start producing events
echo Waiting for the T-Job to start producing events
PROCESSEDEVENTS=0
COUNTER=0
while [ $PROCESSEDEVENTS -eq 0 ];
do
    COUNTER=$((COUNTER+1))
    echo Processed events is still 0 ..
    if [ $COUNTER -eq 50 ]
    then
        echo Counter reached
        exit 1
    fi
    sleep 5
    PROCESSEDEVENTS=$(curl -s "${EMSURL}health" | docker run -i --rm nimmis/jq fromjson.ProcessedEvents)
done
echo received $PROCESSEDEVENTS events
exit 0
