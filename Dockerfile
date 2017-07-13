
FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/elastest/elastest-monitoring-service
COPY go_EMS ./
RUN CGO_ENABLED=0 GOOS=linux go build -o go_EMS go_EMS.go Event.go ChannelInference.go Metric.go SignalsManager.go Definitions.go

FROM docker.elastic.co/logstash/logstash:5.4.0
WORKDIR /root/
COPY --from=builder /go/src/github.com/elastest/elastest-monitoring-service/go_EMS /usr/local/bin/go_EMS
COPY startmeup.sh /startmeup.sh
COPY logstashcfgs/* /usr/share/logstash/pipeline/
ENTRYPOINT ["/startmeup.sh"]
