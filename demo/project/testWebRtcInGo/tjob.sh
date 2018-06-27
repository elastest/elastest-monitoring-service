#!/bin/sh

git clone https://github.com/elastest/full-teaching-experiment

git clone https://github.com/elastest/elastest-monitoring-service
mv elastest-monitoring-service/demo/project/FullTeachingTestE2EDualVideoSession.java full-teaching-experiment/src/test/java/com/fullteaching/backend/e2e/

cd full-teaching-experiment
mvn -Dtest=FullTeachingTestE2EDualVideoSession -B test &

cd
cd elastest-monitoring-service/demo/project/testWebRtcInGo
echo "Building the GO binaries..."
go build -o /usr/local/bin/tjob
echo "Done!"

echo "Executing Go agent..."
/usr/local/bin/tjob
exit $?
