package main

import "time"

func getEvent (rawEvent map[string]interface{}) Event {


	// channel inference
	var channel Channel = "undefined" // default channel: #undefined
	if evchannel, ok := rawEvent["channel"].(string); ok {
		channel = Channel(evchannel)
	} else {
		// Inference rule: if system.load is defined, then it goes through channel #in
		if sysmap, ok := rawEvent["system"].(map[string]interface{}); false && ok {
			if _, ok := sysmap["load"]; ok {
				channel = "in";
			}
		}
	}

	// timestamp inference
	var timestamp string = time.Now().UTC().Format("2006-01-02T15:04:05.000") + "Z" // default timestamp: now
	if ts, ok := rawEvent["@timestamp"].(string); ok {
		ts = ts;
	}

	return Event{channel, rawEvent, timestamp}
}
