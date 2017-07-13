package main

type Event struct {
		Channel Channel
		Payload map[string]interface{}
		Timestamp string
}
