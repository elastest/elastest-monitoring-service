#!/bin/bash

function containerIp () {
	ip=$(docker inspect --format=\"{{.NetworkSettings.Networks."$COMPOSE_PROJECT_NAME"_elastest.IPAddress}}\" "$COMPOSE_PROJECT_NAME"_$1_1 2> /dev/null)
	error=$?
	if [ -z "$2" ]; then
		echo $( echo $ip | cut -f2 -d'"' )
	elif [ "$2" = 'check' ]; then
		echo $error
	fi
}

projectName="elastest"

export COMPOSE_PROJECT_NAME=$projectName

# Start

echo 'Starting ElasTest Platform...'
docker pull elastest/platform
docker run -d -v /var/run/docker.sock:/var/run/docker.sock --rm elastest/platform start --lite --forcepull --nocheck

# Check if ETM container is created
ERROR=$(containerIp "etm" "check")

initial=70
counter=$initial

while [ $ERROR -gt 0 ] ; do
	echo "Waiting to ElasTest ETM container"
	ERROR=$(containerIp "etm" "check")
	sleep 2
	# prevent infinite loop
	counter=$(($counter-1))
		if [ $counter = 0 ]; then
		    echo "Timeout while checking if ETM container is created"
		    exit 1
		fi
done

ET_ETM_HOST=$(docker inspect --format=\"{{.NetworkSettings.Networks.elastest_elastest.IPAddress}}\" elastest_etm_1 2> /dev/null)
export ET_ETM_HOST=$ET_ETM_HOST

echo "ElasTest ETM container is started with IP $ET_ETM_HOST"

docker logs -f elastest_etm_1 &

# Connect test container to docker-compose network

containerId=$(cat /proc/self/cgroup | grep "docker" | sed s/\\//\\n/g | tail -1)

if [ ! -z $containerId ];
then
   echo "Script executing inside the container = ${containerId}"
   docker network connect ${projectName}_elastest ${containerId}
else			
   echo "Script executing in host (not in container)"
fi

echo "Waiting to ETM service ready inside the container (port 8091 available in IP $ET_ETM_HOST)"

# wait ETM started
initial=90
counter=$initial
while ! nc -z -v $ET_ETM_HOST 8091 2> /dev/null; do
	echo "Waiting to ETM ready inside the container"
    sleep 2
    # prevent infinite loop
    counter=$(($counter-1))
    if [ $counter = 0 ]; then
	    echo "Timeout while checking if ETM service is started"
	    exit 1
    fi
done

echo ''
echo "ETM is ready in address $ET_ETM_HOST and port 8091"

echo 'Check if ETM is working... (return 200 OK)'
responseCheck=$(curl --write-out %{http_code} --silent --output /dev/null http://${ET_ETM_HOST}:8091)

if [ $responseCheck = '200' ]; then
	echo "ElasTest ETM is working (returned 200 OK in http://${ET_ETM_HOST}:8091)"
else
	echo "ElasTest ETM is not working (returned ${responseCheck} in http://${ET_ETM_HOST}:8091)"
	exit 1
fi

echo ''
echo ''
echo ''
echo ''
echo ''
echo ''
echo ''
echo "http://$ET_ETM_HOST:8091"


