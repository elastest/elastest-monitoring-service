#!/bin/bash
/usr/bin/printenv
/stress.sh &
/usr/local/bin/docker-entrypoint -e
