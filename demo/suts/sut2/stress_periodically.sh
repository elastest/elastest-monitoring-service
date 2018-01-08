#/bin/bash
while true; do
    curl -H "content-type: application/json" -XPUT "http://${ET_EMS_LSBEATS_HOST}:8181" -d '{ "message" : "Stress", "status" : "on" }' || true
    timeout 5s yes > /dev/null
    curl -H "content-type: application/json" -XPUT "http://${ET_EMS_LSBEATS_HOST}:8181" -d '{ "message" : "Stress", "status" : "off" }' || true
    #echo sleeping..
    sleep 5
    #echo wakeup
done
