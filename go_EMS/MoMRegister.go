package main

import(
	"fmt"
	"encoding/json"
)

func readAndRegister(dasmap map[string]interface{}) {
	fmt.Printf("JSON read: %v\n", dasmap)
	fmt.Printf("defd: %v\n", dasmap["definition"])
	thedef := []byte(dasmap["definition"].(string))
	switch dasmap["momType"] {
	case "sampledSignal" :
		var sampledSig SampledSignalDefinition
		if err := json.Unmarshal(thedef, &sampledSig); err != nil {
			panic("wrong type for unmarshal")
		}
		theGlobalSampledSignalDefs = append(theGlobalSampledSignalDefs, sampledSig)
	case "aggregatedSignal":
		var aggregatedSig AggregatedSignalDefinition
		if err := json.Unmarshal(thedef, &aggregatedSig); err != nil {
			panic("wrong type for unmarshal")
		}
		theGlobalAggregatedSignalDefs[aggregatedSig.Name] = aggregatedSig
	case "writeSignal":
		var writeDef SignalWriteDefinition
		if err := json.Unmarshal(thedef, &writeDef); err != nil {
			panic(err)
		}
		theGlobalWriteDefs = append(theGlobalWriteDefs, writeDef)
	}
}
