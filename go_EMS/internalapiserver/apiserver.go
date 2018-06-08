package internalapiserver

import (
    "log"
    "net"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
	pb "github.com/elastest/elastest-monitoring-service/protobuf"
    "google.golang.org/grpc/reflection"
	pe "github.com/elastest/elastest-monitoring-service/go_EMS/eventscounter"
	ep "github.com/elastest/elastest-monitoring-service/go_EMS/eventproc"
	et "github.com/elastest/elastest-monitoring-service/go_EMS/eventtag"
)

const (
    port = ":50051"
)

// server is used to implement protobuf.Health
type server struct{}

// GetHealth implements protobuf.Health
func (s *server) GetHealth(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
    return &pb.HealthReply{Healthstatus: "Ok", Processedevents: int32(pe.GetProcessedEvents())}, nil
}

// PostMom implements protobuf.PostMoM
func (s *server) PostMoM(ctx context.Context, in *pb.MomPostRequest) (*pb.MomPostReply, error) {
    switch in.Momtype {
    case "tag0.1":
        return et.DeployTaggerv01(in.Momdefinition), nil
    case "signals0.1":
        return ep.DeploySignals01(in.Momdefinition), nil
    }
    return &pb.MomPostReply{Deploymenterror:"Unrecognized tag "+in.Momtype, Momid:""}, nil
}

func Serve() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterEngineServer(s, &server{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
