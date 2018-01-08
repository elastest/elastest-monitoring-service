#!/bin/bash

# Generate .jar
cd ..
mvn -DskipTests=true package

# Generate and publish docker image

# Enable extended patterns
shopt -s extglob

cd target
mv $(echo !(*-sources|*-javadoc).jar) app.jar
mv app.jar ../docker
cd ../docker
docker build -t elastest/demo-rest-java-test-sut .

# Delete unwanted files
rm app.jar