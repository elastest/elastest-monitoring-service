package data

type TaggerDefinition struct {
    InChannels []string `json:"inchannels"`
    Filter string `json:"filter"`
    OutChannel string `json:"outchannel"`
}

type TagCondition struct {
    InChannels ChannelSet
    EventCondition func(payload map[string]interface{}) bool
    OutChannel Channel
}

// AST

type TagNode interface {
    Eval(evPayload map[string]interface{}) bool
}

type AndNode struct {
    Members []TagNode
}

type OrNode struct {
    Members []TagNode
}

type PathFunNode struct {
    Path []string
    Fun func(in interface{}) bool
}
