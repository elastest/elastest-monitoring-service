// Note that this version is completely sequential and doesn't use goroutines at all
package eventscounter

var processedEvents int = 0

func GetProcessedEvents() int {
    return processedEvents
}

func IncrementProcessedEvents() {
    processedEvents++
}
