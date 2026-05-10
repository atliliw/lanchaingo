package openai

import (
	"bufio"
	"io"
	"strings"
)

// StreamEvent represents a single Server-Sent Event.
type StreamEvent struct {
	ID    string
	Event string
	Data  string
}

// ParseSSE reads an SSE stream and returns a channel of StreamEvents.
// The channel is closed when the stream ends or an error occurs that
// prevents further reading.
func ParseSSE(reader io.Reader) <-chan StreamEvent {
	ch := make(chan StreamEvent)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(reader)
		var current StreamEvent

		for scanner.Scan() {
			line := scanner.Text()

			if line == "" {
				if current.Data != "" {
					ch <- current
				}
				current = StreamEvent{}
				continue
			}

			if strings.HasPrefix(line, "id:") {
				current.ID = strings.TrimSpace(line[3:])
			} else if strings.HasPrefix(line, "event:") {
				current.Event = strings.TrimSpace(line[6:])
			} else if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(line[5:])
				if current.Data != "" {
					current.Data += "\n"
				}
				current.Data += data
			}
		}
	
		if current.Data != "" {
			ch <- current
		}
	}()
	return ch
}

// isDoneEvent returns true if the event data signals stream completion.
func isDoneEvent(data string) bool {
	return strings.TrimSpace(data) == "[DONE]"
}
