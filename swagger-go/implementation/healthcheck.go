package implementation

import (
	"net/http"
	openapiruntime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"runtime"
	"encoding/json"
    "log"
    "time"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "github.com/elastest/elastest-monitoring-service/protobuf"
)

const (
    address     = "localhost:50051"
)

type HealthStatus pb.HealthReply

// WriteResponse to the client
func (healthStatus HealthStatus) WriteResponse(rw http.ResponseWriter, producer openapiruntime.Producer) {
    rw.WriteHeader(200)
	healthstatusjson, _ := json.Marshal(healthStatus)
	if err := producer.Produce(rw, string(healthstatusjson)); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

type EnvironmentInfo struct { generalinfo map[string]interface{} }

// WriteResponse to the client
func (info EnvironmentInfo) WriteResponse(rw http.ResponseWriter, producer openapiruntime.Producer) {
    rw.WriteHeader(200)
	envjson, _ := json.Marshal(info.generalinfo)
	if err := producer.Produce(rw, string(envjson)); err != nil {
		panic(err) // let the recovery middleware deal with this
    }
}

func GetHealth() middleware.Responder {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
        // return error instead
    }
    defer conn.Close()
    c := pb.NewEngineClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := c.GetHealth(ctx, &pb.HealthRequest{})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
        // return error instead
    }
	return HealthStatus(*r)
}

func GetEnvironment() middleware.Responder {
	// gather info here!
	var info map[string]interface{} = map[string]interface{} {
		"OS" : runtime.GOOS,
		"Architecture" : runtime.GOARCH,
	}
	return EnvironmentInfo {info}
}
