# Development documentation

## Architecture

The full-fledged version of the EMS is distributed as a single Docker image, running the following processes:
* A logstash instance acting as an input endpoint
* A logstash instance acting as an output endpoint
* A webserver handling the OpenAPI requests
* The EMS engine itself

While the lightweight version runs only these two processes:
* A logstash instance acting as an i/o endpoint
* A webserver handling the OpenAPI requests

### Input Logstash instance

[Full-fledged version]
The configuration file of this Logstash instance is static and specifies the input endpoints that are currently supported by the component. Adding a new input endpoint in future versions of this component means editing this configuration file.

### Output Logstash instance

[Full-fledged version]
The configuration file of this Logstash instance is dynamic and it's manipulated by the webserver upon requests.

### IO Logstash instance

[Lightweight version]
The configuration file of this Logstash instance is dynamic and it's manipulated by the webserver upon requests.

### OpenAPI webserver

The webserver is generated on runtime using [Swagger-go](https://github.com/go-swagger/go-swagger). The implementation of its methods can be found in the directory elastest-monitoring-service/swagger-go/

### The EMS Engine

[Full-fledged version]
The EMS Engine is the core of this component, and implements the logic of the monitoring machines registered via the webserver.

## Prepare development environment

Clone the project from GitHub:
```
$ git clone https://github.com/elastest/elastest-monitoring-service.git
```

Every architecture subcomponent is generated and run inside a docker image, so the development and is carried out in them, making Docker the only requisite for it. Anyway, feel free to set up your own local deveolpment environment.

## Development procedure

### Input logstash instance

[Full-fledged version]
The input logstash instance will read its configuration file from /usr/share/logstash/pipeline/inlogstash.conf, and its output is expected to be written to a FIFO at /usr/share/logstash/pipes/leftpipe.

### Output logstash instance

[Full-fledged version]
The out logstash instance will read its configuration file from /usr/share/logstash/pipeline/outlogstash.conf, and its input is expected to be read from a FIFO at /usr/share/logstash/pipes/rightpipe.

### IO Logstash instance

[Lightweight version]
The io logstash instance will read its configuration file from /usr/share/logstash/pipeline/minlogstash.conf.

### OpenAPI webserver

To generate the webserver, you can run the following command in a shell inside the directory elastest-monitoring-service/swagger-go/ :
```
$ docker run --rm -it -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger generate server -f ../api.yaml
```

Unfortunately, the import paths of the generated files are incorrectly generated. To fix them, you may find useful tuning and running the script file convertpaths.sh located on the same folder.

### The EMS Engine

[Full-fledged version]
To build the engine, run the following command in a shell inside the directory elastest-monitoring-service/swagger-go/ :
```
$ go build -o bin/go_EMS .
```
