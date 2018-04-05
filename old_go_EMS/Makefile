GO = go
GOBUILD = $(GO) build

all: goems 

goems:
	$(GOBUILD) -o bin/go_EMS go_EMS.go Event.go ChannelInference.go Metric.go Session.go SignalsManager.go MetricDefinitions.go SessionDefinitions.go SessionsManager.go MoMRegister.go
