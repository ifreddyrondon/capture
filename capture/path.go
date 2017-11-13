package capture

import (
	"encoding/json"
	"sort"
	"sync"
)

const WorkersNumber = 4

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

type indexCapture struct {
	index   int
	capture *Capture
}

type job struct {
	index int
	data  json.RawMessage
}

func worker(wg *sync.WaitGroup, jobs <-chan job, results chan<- indexCapture) {
	for job := range jobs {
		capture := new(Capture)
		if err := capture.UnmarshalJSON(job.data); err == nil {
			results <- indexCapture{index: job.index, capture: capture}
		}
		wg.Done()
	}
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

	var wg sync.WaitGroup
	wg.Add(len(pj))
	jobs := make(chan job, len(pj))
	results := make(chan indexCapture, len(pj))

	for w := 0; w < WorkersNumber; w++ {
		go worker(&wg, jobs, results)
	}

	for i, v := range pj {
		jobs <- job{index: i, data: v}
	}
	close(jobs)
	wg.Wait()
	close(results)

	processed := make([]indexCapture, len(results))
	idx := 0
	for data := range results {
		processed[idx] = data
		idx++
	}
	sort.Slice(processed, func(i, j int) bool {
		return processed[i].index < processed[j].index
	})
	for _, v := range processed {
		p.AddCapture(v.capture)
	}

	return nil
}
