#!/bin/bash
sleep 100
curl -d "{ \"channel\": \"any\", \"ip\": \"elastest.software.imdea.org\", \"port\": 9202, \"user\": \"elastic\", \"password\": \"changeme\" }" -H "Content-Type: application/json" "http://${ET_EMS_LSBEATS_HOST}:8888/subscriber/elasticsearch"
curl -d "{ \"channel\": \"any\", \"endpoints\":[\"persistence\"] }" -H "Content-Type: application/json" "http://${ET_EMS_LSBEATS_HOST}:8888/subscriber/elastest"
while true; do
    echo STATUS_ON
    timeout 10s yes > /dev/null
    echo STATUS_OFF
    sleep 10
done
