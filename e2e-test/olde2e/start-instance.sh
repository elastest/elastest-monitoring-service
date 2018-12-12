#!/bin/bash

ELASTESTURL=$1

# SuT creation
echo Creating SuT
SUT=$(curl -s -H "Content-Type: application/json" -d @sutdesc.txt "$ELASTESTURL/api/sut")
echo $SUT

# T-Job creation
echo Creating T-Job
TJOB=$(curl -s -H "Content-Type: application/json" -d @tjobdesc.txt "$ELASTESTURL/api/tjob")
echo $TJOB

# T-Job execution
echo Executing T-Job
TJOBEXEC=$(curl -s -H "Content-Type: application/json" -d '{"tJobParams": []}' "$ELASTESTURL/api/tjob/2/exec")
echo $TJOBEXEC

# Getting result
n=0
while [ $n -le 3000 ]
do
	n=$(( n+1 ))	 # increments $n
	sleep 1
	TJOBEXEC=$(curl -s "$ELASTESTURL/api/tjob/2/exec/1/result")
    if [[ $TJOBEXEC = *"SUCCESS"* ]]; then
        echo Test successful
        exit 0
    fi
    if [[ $TJOBEXEC = *"FAIL"* ]]; then
        echo Test failed
        exit -1
    fi
    if [[ $TJOBEXEC = *"ERROR"* ]]; then
        echo Test erroneous
        exit -1
    fi
done

echo Test took too long
exit -1
