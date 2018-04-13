package data

type TaggerDefinition struct {
    InChannels []string `json:"inchannels"`
    Filter string `json:"filter"`
    OutChannel string `json:"outchannel"`
}

type TagCondition struct {
    InChannels ChannelSet
    EventCondition func(ev Event) bool
    OutChannel Channel
}
