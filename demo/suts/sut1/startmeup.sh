#!/bin/bash
/usr/bin/printenv
/stress_periodically.sh &
/usr/local/bin/docker-entrypoint -e
