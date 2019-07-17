package implementation

import (
  "net/http"
  runtime "github.com/go-openapi/runtime"
  middleware "github.com/go-openapi/runtime/middleware"

  "log"
  "time"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "../restapi/operations/monitoring_machine"
  "../restapi/operations/stamper"
  pb "github.com/elastest/elastest-monitoring-service/protobuf"
)

type MomPostReply pb.MomPostReply
type MomDeleteReply pb.MomDeleteReply

// WriteResponse to the client
func (o MomPostReply) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
  var ret string
  if len(o.Deploymenterror) == 0 {
    rw.WriteHeader(200)
    ret = o.Momid
  } else {
    rw.WriteHeader(406)
    ret = o.Deploymenterror
  }
  if err := producer.Produce(rw, ret); err != nil {
    panic(err) // let the recovery middleware deal with this
  }
}

func (o MomDeleteReply) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
  var ret string
  if len(o.Deletionerror) == 0 {
    rw.WriteHeader(200)
    ret = ""
  } else {
    rw.WriteHeader(406)
    ret = o.Deletionerror
  }
  if err := producer.Produce(rw, ret); err != nil {
    panic(err) // let the recovery middleware deal with this
  }
}

type DeployNotAllowed struct { }

// WriteResponse to the client
func (o DeployNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

  rw.WriteHeader(405)
  if err := producer.Produce(rw, "Deployment is not allowed in this version of the ElasTest Monitoring Service"); err != nil {
    panic(err) // let the recovery middleware deal with this
  }
}

func PostMoM(params monitoring_machine.PostMoMParams) middleware.Responder {
  req := pb.MomPostRequest{Momtype:params.Version, Momdefinition:*params.Mom}
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
    // return error instead
  }
  defer conn.Close()
  c := pb.NewEngineClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  r, err := c.PostMoM(ctx, &req)
  if err != nil {
    log.Fatalf("could not greet: %v", err)
    // return error instead
  }
  return MomPostReply(*r)
}

func PostStamper(params stamper.PostStamperParams) middleware.Responder {
  req := pb.MomPostRequest{Momtype:params.Version, Momdefinition:*params.Stamper}
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
    // return error instead
  }
  defer conn.Close()
  c := pb.NewEngineClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  r, err := c.PostStamper(ctx, &req)
  if err != nil {
    log.Fatalf("could not greet: %v", err)
    // return error instead
  }
  return MomPostReply(*r)
}

func DeleteMoM(params monitoring_machine.DeleteMoMParams) middleware.Responder {
  req := pb.MomDeleteRequest{Momid:params.MoMID}
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
    // return error instead
  }
  defer conn.Close()
  c := pb.NewEngineClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  r, err := c.DeleteMoM(ctx, &req)
  if err != nil {
    log.Fatalf("could not greet: %v", err)
    // return error instead
  }
  return MomDeleteReply(*r)
}

func DeleteStamper(params stamper.DeleteStamperParams) middleware.Responder {
  req := pb.MomDeleteRequest{Momid:params.StamperID}
  conn, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
    // return error instead
  }
  defer conn.Close()
  c := pb.NewEngineClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  r, err := c.DeleteStamper(ctx, &req)
  if err != nil {
    log.Fatalf("could not greet: %v", err)
    // return error instead
  }
  return MomDeleteReply(*r)
}
