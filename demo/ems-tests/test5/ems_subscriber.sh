#!/bin/sh
set -e

echo "Waiting for EMS to be ready..."
sleep 15
echo "Done!"
echo "Subscribing elasticsearch..."
curl -d "{ \"channel\": \"any\", \"ip\": \"elastest.software.imdea.org\", \"port\": 9202, \"user\": \"elastic\", \"password\": \"changeme\" }" -H "Content-Type: application/json" "http://${ET_EMS_LSBEATS_HOST}:8888/subscriber/elasticsearch"
echo "Done!"
echo "Waiting for EMS to be ready..."
sleep 15
echo "Done!"
