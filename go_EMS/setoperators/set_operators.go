package setoperators

import (
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
    //"fmt"
)

func SetFromList(channels []string) dt.ChannelSet {
    chans := dt.ChannelSet(map[dt.Channel]interface{}{})
    for _,channel := range channels {
        chans[dt.Channel(channel)] = nil
    }
    return chans
}

func SetIn(ch dt.Channel, chans dt.ChannelSet) bool {
    _,ok := chans[ch]
    return ok
}

func SetMinus(chans dt.ChannelSet, out dt.ChannelSet) dt.ChannelSet {
    for ch,_ := range out {
        delete(chans, ch)
    }
    return chans
}

func SetIsEmpty(chans dt.ChannelSet) bool {
    return len(chans)==0
}

func SetAdd(chans dt.ChannelSet, ch dt.Channel) dt.ChannelSet {
    chans[ch] = nil
    return chans
}

func SetUnion(chans1 dt.ChannelSet, chans2 dt.ChannelSet) dt.ChannelSet {
    for ch,_ := range chans2 {
        chans1[ch] = nil
    }
    return chans1
}

func SetIsIncluded(chans1 dt.ChannelSet, chans2 dt.ChannelSet) bool {
    for ch,_ := range chans1 {
        _,ok := chans2[ch]
        if (!ok) {
            return false
        }
    }
    return true
}

func SetsAreEqual(chans1 dt.ChannelSet, chans2 dt.ChannelSet) bool {
    return SetIsIncluded(chans1, chans2) && SetIsIncluded(chans2, chans1)
}
