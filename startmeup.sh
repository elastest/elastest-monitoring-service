#!/bin/bash -e

mkdir /usr/share/logstash/pipes
mkfifo /usr/share/logstash/pipes/leftpipe
mkfifo /usr/share/logstash/pipes/staticrightpipe
mkfifo /usr/share/logstash/pipes/dynamicrightpipe
mkfifo /usr/share/logstash/pipes/swagpipe
mkfifo /usr/share/logstash/pipes/swageventspipe

mkdir /usr/share/logstash/in_data
mkdir /usr/share/logstash/out_data
mkdir /usr/share/logstash/outstatic_data

logstash -f /usr/share/logstash/pipeline/inlogstash.conf --path.data /usr/share/logstash/in_data &

logstash -f /usr/share/logstash/pipeline/outlogstash.conf --config.reload.automatic --path.data /usr/share/logstash/out_data &
logstash -f /usr/share/logstash/pipeline/staticoutlogstash.conf --path.data /usr/share/logstash/outstatic_data &

swagger --port=8888 --host=0.0.0.0 &

go_EMS "/usr/share/logstash/pipes/staticrightpipe" "/usr/share/logstash/pipes/dynamicrightpipe"
