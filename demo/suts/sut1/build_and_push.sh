#!/bin/bash

# Make sure that you are logged in to docker. To do so, use 
# $ sudo docker login
# and provide username and password

sudo docker build -t imdeasoftware/ems-metricbeat .

sudo docker push imdeasoftware/ems-metricbeat
