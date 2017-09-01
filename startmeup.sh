#!/bin/bash -e

mkdir /usr/share/logstash/pipes
mkfifo /usr/share/logstash/pipes/leftpipe
mkfifo /usr/share/logstash/pipes/rightpipe
mkfifo /usr/share/logstash/pipes/swagpipe

mkdir /usr/share/logstash/in_data
mkdir /usr/share/logstash/out_data

logstash -f /usr/share/logstash/pipeline/inlogstash.conf --path.data /usr/share/logstash/in_data &

logstash -f /usr/share/logstash/pipeline/outlogstash.conf --config.reload.automatic --path.data /usr/share/logstash/out_data &

swagger --port=8888 --host=0.0.0.0 &

go_EMS > /usr/share/logstash/pipes/rightpipe
