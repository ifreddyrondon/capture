package branch

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/ifreddyrondon/gocapture/capture"
)

const WorkersNumber = 4

// Branch represent a collection of captures.
type Branch struct {
	Captures []*capture.Capture `json:"captures"`
}

// AddCapture add new captures into the branch.
// The new captures will always be added at the end of the road.
// respecting their insertion order.
func (p *Branch) AddCapture(captures ...*capture.Capture) {
	p.Captures = append(p.Captures, captures...)
}

type indexCapture struct {
	index   int
	capture *capture.Capture
}

type job struct {
	index int
	data  json.RawMessage
}

func worker(wg *sync.WaitGroup, jobs <-chan job, results chan<- indexCapture) {
	for job := range jobs {
		capture := new(capture.Capture)
		if err := capture.UnmarshalJSON(job.data); err == nil {
			results <- indexCapture{index: job.index, capture: capture}
		}
		wg.Done()
	}
}

// UnmarshalJSON decodes a branch from a JSON body.
// Throws an error if the body of the branch cannot be interpreted as JSON.
// Implements the json.Unmarshaler Interface
func (p *Branch) UnmarshalJSON(data []byte) error {
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

	processed := make([]indexCapture, 0, len(results))
	for data := range results {
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
