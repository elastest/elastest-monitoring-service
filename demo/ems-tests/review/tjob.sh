#!/bin/sh

git clone https://github.com/pchico83/full-teaching-experiment
cd full-teaching-experiment
mvn -Dtest=FullTeachingTestE2EDualVideoSession -B test &

cd
git clone https://github.com/fgorostiaga/elastest-monitoring-service
cd elastest-monitoring-service
git checkout review_demo

cd demo/ems-tests/review
echo Posting stamper
curl -H "Content-Type:text/plain"  --data-binary @stampers.txt http://${ET_EMS_LSBEATS_HOST}:8888/stamper/tag0.1
echo Posting sessions
curl -H "Content-Type:text/plain" --data-binary @sessiondef.txt http://${ET_EMS_LSBEATS_HOST}:8888/MonitoringMachine/signals0.1
echo "Building the GO binaries..."
go build -o /usr/local/bin/tjob
echo "Done!"

echo "Executing Go agent..."
/usr/local/bin/tjob
exit $?
