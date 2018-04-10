#!/bin/bash
while true; do
    echo STATUS_ON
    timeout 10s yes > /dev/null &
    timeout 10s yes > /dev/null &
    timeout 10s yes > /dev/null &
    timeout 10s yes > /dev/null
    echo STATUS_OFF
    sleep 10
done
