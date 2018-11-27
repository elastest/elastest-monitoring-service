#!/bin/bash

ELASTESTURL=$1

function cleanexit() {
    # Project destruction
    echo Destroying project
    curl -X "DELETE" "$ELASTESTURL/api/project/$PROJID"
    exit -1
}

# Project creation
echo Creating Project
PROJ=$(curl -s -H "Content-Type: application/json" -d '{ "id": 0, "name": "EMSe2e" }' "$ELASTESTURL/api/project")
echo $PROJ
PROJID=`echo "$PROJ" | tr '\n' ' ' | docker run -i stedolan/jq '.id'`
echo Proj ID: $PROJID

# SuT creation
echo Creating SuT
DESC=`sed "s/PROJID/$PROJID/" sutdesc.txt`
SUT=$(curl -s -H "Content-Type: application/json" -d "$DESC" "$ELASTESTURL/api/sut")
echo $SUT

# T-Job creation
echo Creating T-Job
DESC=`sed "s/PROJID/$PROJID/" tjobdesc.txt`
TJOB=$(curl -s -H "Content-Type: application/json" -d "$DESC" "$ELASTESTURL/api/tjob")
echo $TJOB

TJOBID=`echo "$TJOB" | tr '\n' ' ' | docker run -i stedolan/jq '.id'`
echo TJob ID: $TJOBID

# T-Job execution
echo Executing T-Job
TJOBEXEC=$(curl -s -H "Content-Type: application/json" -d '{"tJobParams": []}' "$ELASTESTURL/api/tjob/$TJOBID/exec")
echo $TJOBEXEC
TJOBEXECID=`echo "$TJOBEXEC" | tr '\n' ' ' | docker run -i stedolan/jq '.monitoringIndex'`
echo TJobEXEC ID: $TJOBEXECID

# Getting result
n=0
while [ $n -le 3000 ]
do
	n=$(( n+1 ))	 # increments $n
	sleep 1
	TJOBEXEC=$(curl -s "$ELASTESTURL/api/tjob/$TJOBID/exec/$TJOBEXECID/result")
    if [[ $TJOBEXEC = *"SUCCESS"* ]]; then
        echo Test successful
        cleanexit
    fi
    if [[ $TJOBEXEC = *"FAIL"* ]]; then
        echo Test failed
        cleanexit
    fi
    if [[ $TJOBEXEC = *"ERROR"* ]]; then
        echo Test erroneous
        cleanexit
    fi
done

echo Test took too long
cleanexit
