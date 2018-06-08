package main

type SessionSignal interface {
	getState() bool
	// duration?
}

type BaseSessionSignal struct {
	state bool
	// get duration?
}

func (ssignal BaseSessionSignal) getState() bool {
	return ssignal.state
}
