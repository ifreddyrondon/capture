package branch

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/markbates/going/defaults"

	"github.com/ifreddyrondon/gocapture/capture"
)

const (
	workersNumber     = 4
	defaultBranchName = "master"
)

// Branch represent a collection of captures.
type Branch struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Captures []*capture.Capture `json:"captures"`
}

func New(name string, captures ...*capture.Capture) *Branch {
	b := Branch{
		Name:     defaults.String(name, defaultBranchName),
		Captures: captures,
	}
	return &b
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
		var c capture.Capture
		if err := c.UnmarshalJSON(job.data); err == nil {
			results <- indexCapture{index: job.index, capture: &c}
		}
		wg.Done()
	}
}

// UnmarshalJSON decodes a branch from a JSON body.
// Throws an error if the body of the branch cannot be interpreted as JSON.
// Implements the json.Unmarshaler Interface
func (b *Branch) UnmarshalJSON(data []byte) error {
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

	for w := 0; w < workersNumber; w++ {
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
		b.Captures = append(b.Captures, v.capture)
	}

	return nil
}
