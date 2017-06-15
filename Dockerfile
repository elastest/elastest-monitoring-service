FROM golang:1.7.3
WORKDIR /go/src/github.com/elastest/elastest-monitoring-service
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/ems .
ENTRYPOINT ["ems"]
