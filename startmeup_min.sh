#!/bin/bash -e

logstash -f /usr/share/logstash/pipeline/outlogstash.conf --config.reload.automatic &

swagger --port=8888 --host=0.0.0.0
