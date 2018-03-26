FROM golang:1.9.0-alpine3.6 as builder
RUN apk update && \
    apk add --no-cache ca-certificates git
RUN go get github.com/gorilla/websocket

ENV SRC_DIR=/go/src/github.com/elastest/demo/ems-tests/test5
WORKDIR $SRC_DIR
ADD . $SRC_DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/tjob

FROM alpine:3.6
RUN apk update && \
    apk add --no-cache ca-certificates curl

COPY --from=builder /usr/local/bin/tjob /usr/local/bin/tjob

WORKDIR /root/elastest-demo
CMD ["/usr/local/bin/tjob"]
