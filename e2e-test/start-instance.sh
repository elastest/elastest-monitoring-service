#!/bin/bash

ELASTESTURL=localhost:37000

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

curl -H "Content-Type: application/json" -d '{ "channel": "chan", "ip": "localhost", "port": 9201, "user": "elastic", "password": "changeme" }' "${EMSURL}subscriber/elasticsearch"
