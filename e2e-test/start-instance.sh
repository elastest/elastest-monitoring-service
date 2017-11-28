#!/bin/bash

ELASTESTURL=localhost:37000

# Posting EMS instance
INSTANCEID=$(curl -s -XPOST http://$ELASTESTURL/api/esm/services/bab3ae67-8c1d-46ec-a940-94183a443825/prov)

EMSURL=null

while [ "X$EMSURL" = "Xnull" ];
do 
    echo Waiting and retrying..
    sleep 1
    EMSURL=$(curl -s http://$ELASTESTURL/api/esm/services/instances/$INSTANCEID | sudo docker run -i --rm nimmis/jq .urls.api)
done

EMSURL=$(echo "$EMSURL" | tr -d '"')
echo EMS URL: $EMSURL
EMSIPPORT=$(echo "$EMSURL" | sed 's/http:..//;s/:.*//')":5044"
echo EMS IP AND PORT: $EMSIPPORT

# project creation
echo Creating project
PROJ=$(curl -s -H "Content-Type: application/json" -d '{ "id": 0, "name": "proyecto" }' "http://$ELASTESTURL/api/project")
echo $PROJ

# SuT creation
echo Creating SuT
SUT=$(curl -s -H "Content-Type: application/json" -d '{ "id": 1, "name": "dename", "specification": "williamyeh/dummy", "sutType": "MANAGED", "description": "dedescript", "project": '${PROJ}', "instrumentalize": false, "currentSutExec": 1, "instrumentedBy": "WITHOUT", "port": null, "managedDockerType": "IMAGE", "mainService": "mbeat", "parameters": [] }' "http://$ELASTESTURL/api/sut")
echo $SUT

# T-Job creation
echo Creating T-Job
DATA='{ "id": 0, "name": "nombredeltjob", "imageName": "docker.elastic.co/beats/metricbeat:5.4.0", "sut": '"${SUT}"', "project": '"${PROJ}"' , "tjobExecs": [], "parameters": [], "commands": "metricbeat -e -E output.logstash.hosts=[\"'"${EMSIPPORT}"'\"] -E output.elasticsearch.hosts=[\"edm-elasticsearch:9200\"]", '"$(cat suffix.json)"
TJOB=$(curl -s -H "Content-Type: application/json" -d "${DATA}" "http://$ELASTESTURL/api/tjob")
TJOBID=$(echo $TJOB | sudo docker run -i --rm nimmis/jq .id)
echo TJOBID: $TJOBID

# T-Job execution
echo Executing T-Job
TJOBEXEC=$(curl -s -H "Content-Type: application/json" -d '{"tJobParams": []}' "http://$ELASTESTURL/api/tjob/$TJOBID/exec")
echo $TJOBEXEC


# tjob commands: metricbeat -e -E output.logstash.hosts=['172.19.0.14:5044'] -E output.elasticsearch.hosts=['edm-elasticsearch:9200']
# tjob environment docker image: docker.elastic.co/beats/metricbeat:5.4.0

# subscriber. Doesn't seem necessary anymore
# curl -H "Content-Type: application/json" -d '{ "channel": "chan", "ip": "edm-elasticsearch", "port": 9201, "user": "elastic", "password": "changeme" }' "${EMSURL}subscriber/elasticsearch" # might be wrong..
