package internalapiserver

import (
    "log"
    "net"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
	pb "github.com/elastest/elastest-monitoring-service/protobuf"
    "google.golang.org/grpc/reflection"
	pe "github.com/elastest/elastest-monitoring-service/go_EMS/eventscounter"
)

const (
    port = ":50051"
)

// server is used to implement protobuf.Health
type server struct{}

// SayHello implements protobuf.Health
func (s *server) GetHealth(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
    return &pb.HealthReply{Healthstatus: "Ok", Processedevents: int32(pe.GetProcessedEvents())}, nil
}

func Serve() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterHealthServer(s, &server{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
