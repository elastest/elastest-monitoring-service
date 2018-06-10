
FROM quay.io/goswagger/swagger:0.13.0 as swaggerbuilder

# Set Image Labels
ARG GIT_COMMIT=unspecified
LABEL git_commit=$GIT_COMMIT

ARG COMMIT_DATE=unspecified
LABEL commit_date=$COMMIT_DATE

ARG VERSION=unspecified
LABEL version=$VERSION

WORKDIR /go/src/swagger-go
ENV GOPATH /go
COPY swagger-go ./
COPY api.yaml ./swagger.yaml
RUN swagger generate server
RUN sh ./convertpaths.sh

FROM golang:latest as builder2
WORKDIR /go

RUN go get github.com/go-openapi/analysis
RUN go get github.com/go-openapi/errors
RUN go get github.com/go-openapi/loads
RUN go get github.com/go-openapi/spec
RUN go get github.com/go-openapi/strfmt
RUN go get github.com/go-openapi/swag
RUN go get github.com/go-openapi/validate
RUN go get github.com/tylerb/graceful
RUN go get github.com/go-openapi/runtime/flagext
RUN go get github.com/jessevdk/go-flags
RUN go get github.com/golang/protobuf/proto
RUN go get google.golang.org/grpc
RUN go get google.golang.org/grpc/reflection

COPY . /go/src/github.com/elastest/elastest-monitoring-service
COPY --from=swaggerbuilder /go/src/swagger-go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o swagger cmd/monitoring-as-a-service-server/main.go

FROM golang:latest as builder
WORKDIR /go/src/github.com/elastest/elastest-monitoring-service
RUN go get github.com/golang/protobuf/proto
RUN go get google.golang.org/grpc
RUN go get google.golang.org/grpc/reflection
RUN go get github.com/mna/pigeon
COPY . /go/src/github.com/elastest/elastest-monitoring-service
RUN mkdir -p /go/src/gitlab.software.imdea.org/felipe.gorostiaga
RUN mv /go/src/github.com/elastest/elastest-monitoring-service/striver-go /go/src/gitlab.software.imdea.org/felipe.gorostiaga
RUN cd go_EMS/parsers/session; make; cd -
RUN cd go_EMS/parsers/stamp; make; cd -
RUN CGO_ENABLED=0 GOOS=linux go build -o ems ./go_EMS

FROM docker.elastic.co/logstash/logstash:5.4.0
WORKDIR /root/
USER root
RUN /usr/share/logstash/bin/logstash-plugin install logstash-output-websocket
COPY new_startmeup.sh /startmeup.sh
COPY logstashcfgs/* /usr/share/logstash/pipeline/
COPY keystore.jks /
COPY logstash.yml /usr/share/logstash/config/logstash.yml
RUN chmod 666 /usr/share/logstash/pipeline/outlogstash.conf
RUN chmod 666 /usr/share/logstash/pipeline/staticoutlogstash.conf
RUN mkdir /usr/share/logstash/pipes
RUN mkfifo /usr/share/logstash/pipes/leftpipe
RUN mkfifo /usr/share/logstash/pipes/staticrightpipe
RUN mkfifo /usr/share/logstash/pipes/dynamicrightpipe
RUN mkdir /usr/share/logstash/in_data
RUN mkdir /usr/share/logstash/out_data
RUN mkdir /usr/share/logstash/outstatic_data

COPY --from=builder /go/src/github.com/elastest/elastest-monitoring-service/ems /usr/local/bin/go_EMS
COPY --from=builder2 /go/swagger /usr/local/bin/swagger

ENTRYPOINT ["/startmeup.sh"]
