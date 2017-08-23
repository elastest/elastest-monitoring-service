FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/elastest/elastest-monitoring-service
COPY go_EMS ./
RUN CGO_ENABLED=0 GOOS=linux go build -o go_EMS go_EMS.go Event.go ChannelInference.go Metric.go Session.go SignalsManager.go MetricDefinitions.go SessionDefinitions.go SessionsManager.go
# old workdir would lead to weird error of package not found
WORKDIR /go
COPY swagger-go ./
RUN go get github.com/go-openapi/errors
RUN go get github.com/go-openapi/loads
RUN go get github.com/go-openapi/runtime
RUN go get github.com/go-openapi/runtime/flagext
RUN go get github.com/go-openapi/runtime/middleware
RUN go get github.com/go-openapi/runtime/security
RUN go get github.com/go-openapi/spec
RUN go get github.com/go-openapi/strfmt
RUN go get github.com/go-openapi/swag
RUN go get github.com/jessevdk/go-flags
RUN go get github.com/tylerb/graceful
RUN CGO_ENABLED=0 GOOS=linux go build -o swagger cmd/monitoring-as-a-service-server/main.go

FROM docker.elastic.co/logstash/logstash:5.4.0
WORKDIR /root/
COPY --from=builder /go/src/github.com/elastest/elastest-monitoring-service/go_EMS /usr/local/bin/go_EMS
COPY --from=builder /go/swagger /usr/local/bin/swagger
COPY startmeup.sh /startmeup.sh
COPY logstashcfgs/* /usr/share/logstash/pipeline/
USER root
RUN chmod 666 /usr/share/logstash/pipeline/outlogstash.conf
ENTRYPOINT ["/startmeup.sh"]
