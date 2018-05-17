#!/bin/sh
set -e

apk update
apk add --no-cache ca-certificates git curl

git clone https://github.com/elastest/elastest-monitoring-service;
cd elastest-monitoring-service

cd demo/ems-tests/test6
echo "Installing dependencies..."
go get github.com/gorilla/websocket
echo "Done!"
echo "Building the binaries..."
go build -o /usr/local/bin/tjob
echo "Done!"

exec /usr/local/bin/tjob
