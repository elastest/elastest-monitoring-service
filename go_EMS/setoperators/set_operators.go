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
    retchans := dt.ChannelSet(map[dt.Channel]interface{}{})
    for ch,_ := range chans {
        if _,ok := out[ch]; !ok {
            retchans[ch] = nil
        }
    }
    return retchans
}

func SetIsEmpty(chans dt.ChannelSet) bool {
    return len(chans)==0
}

func SetAdd(chans dt.ChannelSet, ch dt.Channel) dt.ChannelSet {
    retchans := dt.ChannelSet(map[dt.Channel]interface{}{})
    for ch,_ := range chans {
        retchans[ch] = nil
    }
    retchans[ch] = nil
    return retchans
}

func SetUnion(chans1 dt.ChannelSet, chans2 dt.ChannelSet) dt.ChannelSet {
    retchans := dt.ChannelSet(map[dt.Channel]interface{}{})
    for ch,_ := range chans1 {
        retchans[ch] = nil
    }
    for ch,_ := range chans2 {
        retchans[ch] = nil
    }
    return retchans
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
