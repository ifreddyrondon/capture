package capture

import (
	"encoding/json"
	"sort"
)

// Path represent an array of captures.
type Path struct {
	Captures []*Capture
}

// AddCapture add new captures into the path.
// The new captures will always be added at the end of the road.
// respecting their insertion order.
func (p *Path) AddCapture(captures ...*Capture) {
	p.Captures = append(p.Captures, captures...)
}

type channelData struct {
	index   int
	capture *Capture
}

// UnmarshalJSON decodes a path from a JSON body.
// Throws an error if the body of the path cannot be interpreted as JSON.
// Implements the json.Unmarshaler Interface
func (p *Path) UnmarshalJSON(data []byte) error {
	var pj []json.RawMessage
	if err := json.Unmarshal(data, &pj); err != nil {
		return err
	}

	if len(pj) == 0 {
		return nil
	}

	jobs := make(chan channelData, len(pj))
	done := make(chan bool)
	var readyCounter int

	for i, v := range pj {
		go func(index int, data json.RawMessage) {
			capture := new(Capture)
			if err := capture.UnmarshalJSON(data); err == nil {
				jobs <- channelData{index: i, capture: capture}
			}

			readyCounter++
			if readyCounter == len(pj) {
				close(jobs)
				done <- true
			}
		}(i, v)
	}

	<-done
	var processed []channelData
	for data := range jobs {
		processed = append(processed, data)
	}
	sort.Slice(processed, func(i, j int) bool {
		return processed[i].index < processed[j].index
	})
	for _, v := range processed {
		p.AddCapture(v.capture)
	}

	return nil
}
