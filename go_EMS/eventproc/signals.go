package eventproc

import (
  dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
  pb "github.com/elastest/elastest-monitoring-service/protobuf"
  "math/rand"
  "strconv"
  "fmt"
  "github.com/elastest/elastest-monitoring-service/go_EMS/moms"
  "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/session"
  "strings"
)

var momengines map[int]dt.MoMEngine01 = make(map[int]dt.MoMEngine01)

func DeploySignals01(signalsDef string) *pb.MomPostReply {
  fmt.Println("Deploying signal defs: ", signalsDef)
  reader := strings.NewReader(signalsDef)
  monitorif, err := session.ParseReader("Monitoring Machine", reader)
  if err != nil {
    fmt.Println("deployment error: ", err.Error())
    return &pb.MomPostReply{Deploymenterror:err.Error(), Momid:""}
  }
  momifs := monitorif.([]interface{})
  moms := make([]session.MoM, len(momifs))
  for i,momif := range momifs {
    moms[i] = momif.(session.MoM)
  }
  momid := rand.Int()
  deployRealMoM01(moms, momid)
  fmt.Println("with momid: ", momid)
  return &pb.MomPostReply{Deploymenterror:"", Momid:strconv.Itoa(momid)}
}

func deployRealMoM01(signaldefs []session.MoM, momid int) {
  momengines[momid] = moms.StartEngine(signaldefs)
}

func DeleteSignal(momidstr string) *pb.MomDeleteReply {
  momid, err := strconv.Atoi(momidstr);
  if err != nil {
    return &pb.MomDeleteReply{Deletionerror:err.Error()}
  }
  if engine,ok := momengines[momid]; ok {
    delete(momengines, momid)
    engine.Kill()
    return &pb.MomDeleteReply{Deletionerror:""}
  }
  return &pb.MomDeleteReply{Deletionerror:"No such id"}
}

func ProcessEvent(evt dt.Event) {
  for _,engine := range momengines {
    engine.Sampler.ProcessEvent(evt)
  }
}

