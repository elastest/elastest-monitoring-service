#!/bin/bash -e

mkdir /usr/share/logstash/pipes
mkfifo /usr/share/logstash/pipes/leftpipe
mkfifo /usr/share/logstash/pipes/rightpipe

mkdir /usr/share/logstash/in_data
mkdir /usr/share/logstash/out_data

logstash -f /usr/share/logstash/pipeline/inlogstash.conf --path.data /usr/share/logstash/in_data &

logstash -f /usr/share/logstash/pipeline/outlogstash.conf --path.data /usr/share/logstash/out_data &

go_EMS < /usr/share/logstash/pipes/leftpipe > /usr/share/logstash/pipes/rightpipe
