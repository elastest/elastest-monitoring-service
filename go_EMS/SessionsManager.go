package main

import "time"
import "errors"
import "fmt"

type SignalIdToBaseSession struct {
	sigid SignalNameAndPars
	signal *BaseSessionSignal
}

type SessionParsAndSignal struct {
	params map[Param]string
	signal *SessionSignal
}

var baseSessionMan []*SignalIdToBaseSession

// stub
type StubSessSignal struct {}
func (ssignal StubSessSignal) getState() bool {
	return time.Now().Minute() %2 == 0
}

var stubSessionSignal SessionSignal = StubSessSignal {}

var stubSessionParsAndSignals []SessionParsAndSignal = []SessionParsAndSignal {
	SessionParsAndSignal {nil, &stubSessionSignal},
}
// end of stub

func getSessionSignals(sessionName SignalName, conditionBoundParams map[Param]string) []SessionParsAndSignal {
	var ret []SessionParsAndSignal = nil
	for _, sigidandsignal := range baseSessionMan {
		sigid := sigidandsignal.sigid
		add := false
		if (sigid.signalName == sessionName) {
			add = true
			for p,v := range conditionBoundParams {
				if (sigid.parameters[p] != v) {
					add = false
				}
			}
		}
		if (add) {
			var signal SessionSignal = *sigidandsignal.signal
			ret = append(ret, SessionParsAndSignal{sigid.parameters, &signal})
		}
	}
	return ret
}

func registerBaseSessionSignal(signalid SignalNameAndPars, signal *BaseSessionSignal) error {
	for _, entry := range baseSessionMan {
		if (signalid.equals(entry.sigid)) {
			return errors.New("entry already exists")
		}
	}
	baseSessionMan = append(baseSessionMan, &SignalIdToBaseSession{signalid, signal})
	return nil
}

func getBaseSession(sessionid SignalNameAndPars) (*BaseSessionSignal, error) {
	for _, entry := range baseSessionMan {
		if (sessionid.equals(entry.sigid)) {
			return entry.signal, nil
		}
	}
	return &BaseSessionSignal{}, errors.New("no such entry")
}

func updateBaseSession(signalpars SignalNameAndPars, value bool) {
	theSignal, err := getBaseSession(signalpars)
	if (err != nil) {
		theSignal = createBaseSession(signalpars)
	}
	fmt.Printf("session %v is now %v\n", signalpars, value)
	theSignal.state = value
	fmt.Printf("csignalupdated: %p\n", theSignal)
}

func createBaseSession(signalpars SignalNameAndPars) *BaseSessionSignal {
	ret := &BaseSessionSignal{false}
	err := registerBaseSessionSignal(signalpars, ret)
	if (err!= nil) {
		panic(err)
	}
	reportSessionSignalCreation(signalpars, ret)
	return ret
}


func reportSessionSignalCreation(srcSignalId SignalNameAndPars, srcSignal SessionSignal) {
	sName := srcSignalId.signalName
	sPars := srcSignalId.parameters

	// create conditional signals
	arr, ok := conditionalSignalCreationMap[sName]
	if (ok) {
		for _, inducedSignal := range arr {
			theDefinition,ok := theGlobalConditionalSignalDefs[inducedSignal.signalName]
			// assert ok
			if (!ok) {
				// error
				panic("nosuchcondsig")
			}

			conditionBoundParams := make(map[Param]string)
			for sessParam, myParam := range inducedSignal.reboundSessionParameters {
				conditionBoundParams[myParam] = sPars[sessParam]
			}

			signalBoundParams := make(map[Param]string)
			for srcParam, myParam := range inducedSignal.reboundMetricParameters {
				val, ok := signalBoundParams[myParam]
				if (ok) {
					signalBoundParams[srcParam] = val
				}
			}

			metricParsAndSignals := getSignals(theDefinition.condition, signalBoundParams)

			for _, metricParsAndSignal := range metricParsAndSignals {
				paramvals := make(map[Param]string)
				for k,v := range conditionBoundParams {
					paramvals[k] = v
				}
				for k,v := range metricParsAndSignal.params {
					paramvals[inducedSignal.reboundMetricParameters[k]] = v
				}
				// assert paramvals are all the parameters
				nameAndPars := SignalNameAndPars{inducedSignal.signalName, paramvals}
				createConditionalSignal(nameAndPars, &srcSignal, metricParsAndSignal.signal)
			}
		}
	}
}
