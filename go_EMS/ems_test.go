package main

import (
	"testing"
	"encoding/json"
)

func TestChannelInference(t *testing.T) {

	tables := []struct {
        json string
        channel string
    }{
        {"{\"channel\":\"algo\"}", "algo"},
        {"{\"otherfields\":\"somevals\"}", "undefined"},
    }

    for _, table := range tables {
		var rawEvent map[string]interface{} = nil
		thejson := []byte(table.json)
		json.Unmarshal(thejson, &rawEvent)
		inferredchan := string(getEvent(rawEvent).Channel)
		if inferredchan != table.channel {
			t.Errorf("Wrong inferred channel, got: %s, want: %s.", inferredchan, table.channel)
		}
	}
}
