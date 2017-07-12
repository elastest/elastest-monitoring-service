GO = go
GOBUILD = $(GO) build

all: timestamper goems 

timestamper:
	$(GOBUILD) -o bin/timestamper timestamper.go

goems:
	$(GOBUILD) -o bin/go_EMS go_EMS.go Event.go ChannelInference.go Metric.go SignalsManager.go Definitions.go
