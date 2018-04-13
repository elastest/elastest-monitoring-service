package data

type Channel string

type ChannelSet map[Channel]interface{}

type Event struct {
    Channels ChannelSet
    Payload map[string]interface{}
    Timestamp string
}
