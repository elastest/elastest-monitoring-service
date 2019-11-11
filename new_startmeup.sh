#!/bin/bash -e

if [ -n "$EMS_STANDALONE" ]
then
  mv /usr/share/logstash/pipeline/standaloneinlogstash.conf /usr/share/logstash/pipeline/inlogstash.conf
fi

logstash -f /usr/share/logstash/pipeline/inlogstash.conf --config.reload.automatic --path.data /usr/share/logstash/in_data &>/inlogs.txt &

logstash -f /usr/share/logstash/pipeline/outlogstash.conf --config.reload.automatic --path.data /usr/share/logstash/out_data &>/outlogs.txt &
logstash -f /usr/share/logstash/pipeline/staticoutlogstash.conf --path.data /usr/share/logstash/outstatic_data &>/outstaticlogs.txt &

swagger --port=8888 --host=0.0.0.0 &>/swaggerlogs.txt &

# (while true; do cat /usr/share/logstash/pipes/leftpipe; done; echo GBFST) | tee /gologs.txt /usr/share/logstash/pipes/staticrightpipe >(while true; do cat > /usr/share/logstash/pipes/dynamicrightpipe; done; echo GODDBYE)

go_EMS "/usr/share/logstash/pipes/staticrightpipe" "/usr/share/logstash/pipes/dynamicrightpipe" | tee /gologs.txt
