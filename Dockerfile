FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/elastest/elastest-monitoring-service
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ems .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/elastest/elastest-monitoring-service/ems /usr/local/bin/ems
ENTRYPOINT ["ems"]
