package capture

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
)

const (
	maxBulkPayload = 100
	workersNumber  = 4
)

var (
	// ErrorNoCapturesPayload expected error when not found captures when try to unmarshal it.
	ErrorNoCapturesPayload = errors.New("cannot unmarshal json into valid captures, it needs at least one valid capture")
	// ErrorMaxPayloadSize expected error when payload list is greater than 100.
	ErrorMaxPayloadSize = errors.New("limited to 100 calls in a single batch request. If it needs to make more calls than that, use multiple batch requests")
)

// Captures represent a collection of capture in any particular order
type Captures []*Capture

// UnmarshalJSON decodes a collection of captures from a JSON body.
// Throws an error if the body of the branch cannot be interpreted as JSON.
// Implements the json.Unmarshaler Interface
func (c *Captures) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		var capt Capture
		if err := capt.UnmarshalJSON(data); err != nil {
			return err
		}
		*c = append(*c, &capt)
		return nil
	}

	if len(raw) == 0 {
		return ErrorNoCapturesPayload
	} else if len(raw) > maxBulkPayload {
		return ErrorMaxPayloadSize
	}

	var wg sync.WaitGroup
	wg.Add(len(raw))
	jobs := make(chan job, len(raw))
	results := make(chan indexCapture, len(raw))

	for w := 0; w < workersNumber; w++ {
		go worker(&wg, jobs, results)
	}

	for i, v := range raw {
		jobs <- job{index: i, data: v}
	}
	close(jobs)
	wg.Wait()
	close(results)

	if len(results) == 0 {
		return ErrorNoCapturesPayload
	}

	processed := make([]indexCapture, 0, len(results))
	for data := range results {
		processed = append(processed, data)
	}
	sort.Slice(processed, func(i, j int) bool {
		return processed[i].index < processed[j].index
	})
	for _, v := range processed {
		*c = append(*c, v.capture)
	}

	return nil
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
		var c Capture
		if err := c.UnmarshalJSON(job.data); err == nil {
			results <- indexCapture{index: job.index, capture: &c}
		}
		wg.Done()
	}
}
