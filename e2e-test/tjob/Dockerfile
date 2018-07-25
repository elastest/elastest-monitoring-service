FROM golang:latest
WORKDIR /go

COPY . /go
RUN go get github.com/gorilla/websocket
RUN CGO_ENABLED=0 GOOS=linux go build -o tjob .
ENTRYPOINT ["/go/tjob"]
