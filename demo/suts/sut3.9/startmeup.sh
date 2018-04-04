#!/bin/bash
/usr/bin/printenv
# wait for ems
sleep 100
# subscribe es
curl -d "{ \"channel\": \"any\", \"ip\": \"elastest.software.imdea.org\", \"port\": 9202, \"user\": \"elastic\", \"password\": \"changeme\" }" -H "Content-Type: application/json" "http://${ET_EMS_LSBEATS_HOST}:8888/subscriber/elasticsearch"
sleep 15
/stress.sh &
/usr/local/bin/docker-entrypoint -e
