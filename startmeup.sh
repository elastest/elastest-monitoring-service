#!/bin/bash -e

mkdir /usr/share/logstash/pipes
mkfifo /usr/share/logstash/pipes/leftpipe
mkfifo /usr/share/logstash/pipes/rightpipe
mkfifo /usr/share/logstash/pipes/swagpipe
mkfifo /usr/share/logstash/pipes/swageventspipe

mkdir /usr/share/logstash/in_data
mkdir /usr/share/logstash/out_data

EDM_ES=$(echo $ET_EDM_ELASTICSEARCH_API | sed 's/http:..//;s/:.*//')
ES_IP=$(ping -c 1 $EDM_ES | head -n1 | sed 's/).*//;s/.*(//')
echo -e "$ES_IP\telasticsearch" >> /etc/hosts

logstash -f /usr/share/logstash/pipeline/inlogstash.conf --path.data /usr/share/logstash/in_data &

logstash -f /usr/share/logstash/pipeline/outlogstash.conf --config.reload.automatic --path.data /usr/share/logstash/out_data &

swagger --port=8888 --host=0.0.0.0 &

go_EMS > /usr/share/logstash/pipes/rightpipe
