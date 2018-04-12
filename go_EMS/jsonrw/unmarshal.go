package jsonrw

import (
    "time"
	dt "github.com/elastest/elastest-monitoring-service/go_EMS/datatypes"
	sets "github.com/elastest/elastest-monitoring-service/go_EMS/setoperators"
)

func ReadEvent (rawEvent map[string]interface{}) dt.Event {
    var channels dt.ChannelSet
    if evchannels, ok := rawEvent["channels"].([]interface{}); ok {
        strchans := make([]string, len(evchannels))
        for i:= range evchannels {
            strchans[i] = evchannels[i].(string)
        }
        channels = sets.SetFromList(strchans)
		delete(rawEvent, "channels")
    } else {
        channels = make(dt.ChannelSet) // default channels: none
    }

    // timestamp inference
    var timestamp string = time.Now().UTC().Format("2006-01-02T15:04:05.000") + "Z" // default timestamp: now
    if ts, ok := rawEvent["@timestamp"].(string); ok {
        timestamp = ts;
    }

    return dt.Event{channels, rawEvent, timestamp}
}
