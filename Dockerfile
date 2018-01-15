
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

FROM golang:latest as builder
WORKDIR /go/src/github.com/elastest/elastest-monitoring-service
COPY . /go/src/github.com/elastest/elastest-monitoring-service
RUN CGO_ENABLED=0 GOOS=linux go build -o ems ./go_EMS

WORKDIR /go
COPY --from=swaggerbuilder /go/src/swagger-go ./
RUN go get github.com/go-openapi/runtime/flagext
RUN go get github.com/jessevdk/go-flags
COPY vendor /go/src
RUN CGO_ENABLED=0 GOOS=linux go build -o swagger cmd/monitoring-as-a-service-server/main.go

FROM docker.elastic.co/logstash/logstash:5.4.0
WORKDIR /root/
COPY --from=builder /go/src/github.com/elastest/elastest-monitoring-service/ems /usr/local/bin/go_EMS
COPY --from=builder /go/swagger /usr/local/bin/swagger
COPY startmeup.sh /startmeup.sh
COPY logstashcfgs/* /usr/share/logstash/pipeline/
COPY logstash.yml /usr/share/logstash/config/logstash.yml
USER root
RUN chmod 666 /usr/share/logstash/pipeline/outlogstash.conf
ENTRYPOINT ["/startmeup.sh"]
