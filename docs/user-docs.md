# User documentation of the ElasTest Monitoring Service (EMS)

The ElasTest Monitoring Service (EMS)provides a monitoring infrastructure suitable for inspecting executions of a System Under Test (hereinafter "SuT") and the ElasTest platform itself online.
This service allows the user and the platform to deploy machines able to process events in real time and generate complex, higher level events from them. This can help to better understand what's happening, detect anomalies, correlate issues, and even stress the tests automatically; all of which aims to maximize the chances of uncover bugs and their causes.
To achieve its goal, it provides an OpenAPI endpoint whose specification can be found at http://elastest.io/docs/api/ems, along with input endpoints over which events can be fed to the service.

## Features
### Subscription and event feeding
The version 0.1 of the ElasTest Monitoring Service accepts the subscription of rabbitmq and elasticsearch endpoints, which will receive all the events sent through all the cannels.
The only current input endpoint for events is for [beats](https://www.elastic.co/products/beats) on port 5044.
### Monitoring machines deployment
The EMS comes in two flavors: a full-fledged version, which supports the deployment of Monitoring machines, and a lightweight one, which does not. If a MoM method is called on a lightweight EMS, an error is returned.

## How to run

To run the EMS as a standalone component, you can download the docker-compose file available at https://github.com/elastest/elastest-monitoring-service/blob/master/docker-compose.yml and then run it with the following command line:
```
$ docker-compose up
```

## Basic usage

### Subscription

When the EMS is started, a server for managing the monitoring machines and the subscription endpoints, in compliance with the [EMS API](http://elastest.io/docs/api/ems) is started at port 8888.
As specified by the API, the user can subscribe a new RabbitMQ endpoint by executing the following command:
```
$ echo '{"channel": "in", "ip": "rabbitHost", "port": 5672, "user": "rabbituser", "password": "passw0rd", "key": "key", "exchange": "exc", "exchange_type": "fanout"}' | curl -i -H "Content-Type: application/json" --data @- http://127.0.0.1:8888/subscriber/rabbitmq
```

### Event feeding

A client can send events to the EMS configuring a beats server to send its output as specified by the following lines in its configuration file:
```
output.logstash:
  hosts: ["logstash:5044"]
```

### Monitoring machines deployment
[Full-fledged version only]
The tester can deploy monitoring machines by performing a POST request to the MonitoringMachine path, providing the parameters specified at the [EMS API](http://elastest.io/docs/api/ems). The definition depends on the type of monitoring machine being deployed (field "momType"), and consists of a JSON with the corresponding format.
Examples for deploying SampledSignals and WriteDefinitions can be found in the directory elastest-monitoring-service/go\_EMS/testinputs. For example, to declare a metric called "cpuload x", listening over channel "in", whose value is scrapped from field "event[system][load][1]" and whose parameter "x" must be retrieved from field "event[beat][hostname]", we can execute the following command:
```
$ echo '{ "momType":"sampledSignal", "definition":"{ \"name\": \"cpuload\", \"paramsPaths\": {\"x\": \"beat.hostname\"}, \"inChannel\": \"in\", \"valuePath\": \"system.load.1\" }"}' | curl -i -H "Content-Type: application/json" --data @- http://127.0.0.1:8888/MonitoringMachine
```

