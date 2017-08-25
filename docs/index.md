# ElasTest Monitoring Service (EMS)

The ElasTest Monitoring Service (EMS)provides a monitoring infrastructure suitable for inspecting executions of a System Under Test (hereinafter "SuT") and the ElasTest platform itself online.
This service allows the user and the platform to deploy machines able to process events in real time and generate complex, higher level events from them. This can help to better understand what's happening, detect anomalies, correlate issues, and even stress the tests automatically; all of which aims to maximize the chances of uncover bugs and their causes.
To achieve its goal, it provides an OpenAPI endpoint whose specification can be found at http://elastest.io/docs/api/ems, along with input endpoints over which events can be fed to the service.

## Features
The version 0.1 of the ElasTest Monitoring Service accepts the subscription of rabbitmq and elasticsearch endpoints, which will receive all the events sent through all the cannels.
The only current input endpoint for events is for [beats](https://www.elastic.co/products/beats) on port 5044.

## How to run

To run the EMS as a standalone component, you can download the docker-compose file available at https://github.com/elastest/elastest-monitoring-service/blob/master/docker-compose.yml and then run it with the following command line:
```
$ docker-compose up
```

## Basic usage

When the EMS is started, a server for managing the monitoring machines and the subscription endpoints, in compliance with the [EMS API](http://elastest.io/docs/api/ems) is started at port 8888.
As specified by the API, the user can subscribe a new RabbitMQ endpoint by executing the following command:
```
$ echo '{"channel": "in", "ip": "rabbitHost", "port": 5672, "user": "rabbituser", "password": "passw0rd", "key": "key", "exchange": "exc", "exchange_type": "fanout"}' | curl -i -H "Content-Type: application/json" --data @- http://127.0.0.1:8888/subscriber/rabbitmq
```

A client can send events to the EMS configuring a beats server to send its output as specified by the following lines in its configuration file:
```
output.logstash:
  hosts: ["logstash:5044"]
```

## Development documentation

### Architecture

The EMS is distributed as a single Docker image, running the following processes:
* A logstash instance acting as an input endpoint
* A logstash instance acting as an output endpoint
* A webserver handling the OpenAPI requests
* The EMS engine itself

#### Input Logstash instance

The configuration file of this Logstash instance is static and specifies the input endpoints that are currently supported by the component. Adding a new input endpoint in future versions of this component means editing this configuration file.

#### Output Logstash instance

The configuration file of this Logstash instance is dynamic and it's manipulated by the webserver upon requests.

#### OpenAPI webserver

The webserver is generated on runtime using [Swagger-go](https://github.com/go-swagger/go-swagger). The implementation of its methods can be found in the directory elastest-monitoring-service/swagger-go/

#### The EMS Engine

The EMS Engine is the core of this component, and implements the logic of the monitoring machines registered via the webserver.

### Prepare development environment

Clone the project from GitHub:
```
$ git clone https://github.com/elastest/elastest-monitoring-service.git
```

Every architecture subcomponent is generated and run inside a docker image, so the development and is carried out in them, making Docker the only requisite for it. Anyway, feel free to set up your own local deveolpment environment.

### Development procedure

#### Input logstash instance

The input logstash instance will read its configuration file from /usr/share/logstash/pipeline/inlogstash.conf, and its output is expected to be written to a FIFO at /usr/share/logstash/pipes/leftpipe.

#### Output logstash instance

The out logstash instance will read its configuration file from /usr/share/logstash/pipeline/outlogstash.conf, and its input is expected to be read from a FIFO at /usr/share/logstash/pipes/rightpipe.

#### OpenAPI webserver

To generate the webserver, you can run the following command in a shell inside the directory elastest-monitoring-service/swagger-go/ :
```
$ docker run --rm -it -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger generate server -f ../api.yaml
```

Unfortunately, the import paths of the generated files are incorrectly generated. To fix them, you may find useful tuning and running the script file convertpaths.sh located on the same folder.

#### The EMS Engine

To build the engine, run the following command in a shell inside the directory elastest-monitoring-service/swagger-go/ :
```
$ go build -o bin/go_EMS .
```
