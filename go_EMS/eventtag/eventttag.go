package eventtag

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
    pb "github.com/elastest/elastest-monitoring-service/protobuf"
    "github.com/elastest/elastest-monitoring-service/go_EMS/parsers/stamp"
    "math/rand"
    "strconv"
    "fmt"
    "strings"
)

var tagMonitors map[int]stamp.Filters = make(map[int]stamp.Filters)

func DeployTaggerv01(taggerDef string) *pb.MomPostReply {
    fmt.Println("Deploying def: ", taggerDef)
    reader := strings.NewReader(taggerDef)
    monitorif, err := stamp.ParseReader("Tagger", reader)
    if err != nil {
        fmt.Println("deployment error: ", err.Error())
        return &pb.MomPostReply{Deploymenterror:err.Error(), Momid:""}
    }
    monitor := monitorif.(stamp.Filters)
    momid := rand.Int()
    DeployRealSamplerv01(monitor, momid)
    fmt.Println("with momid: ", momid)
    return &pb.MomPostReply{Deploymenterror:"", Momid:strconv.Itoa(momid)}
}

func DeployRealSamplerv01(monitor stamp.Filters, momid int) {
    // TODO make this method private in the future
    tagMonitors[momid] = monitor
}

func TagEvent(ev *dt.Event) {
    checkDefs := []stamp.Filter{}
    for _,monitor := range tagMonitors {
        for _,def := range monitor.Defs {
            checkDefs = append(checkDefs, def)
        }
    }
    dirty:=true

    for (dirty) {
        dirty=false

		tmp := checkDefs[:0]
        for _,def := range checkDefs {
            if def.Pred.Eval(*ev) {
                dirty = true
                (*ev).Channels = sets.SetAdd(ev.Channels, def.Tag.Tag)
            } else if !sets.SetIn(def.Tag.Tag, ev.Channels) {
                tmp = append(tmp, def)
            }
        }
		checkDefs = tmp
    }
}
