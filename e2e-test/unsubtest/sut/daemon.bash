#!/bin/bash

while true
do
  sleep 1
  STATS=`curl http://localhost:9200/logstash-*/_stats`
  #echo $STATS
  curl -H "content-type: application/json" -XPUT "http://localhost:8181" -d $STATS
done
