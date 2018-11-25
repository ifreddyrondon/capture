package decoder

import (
	"github.com/lib/pq"
)

type PostTags struct {
	Tags []string `json:"tags"`
}

func (t *PostTags) OK() error {
	return nil
}

func (t *PostTags) GetTags() pq.StringArray {
	if t.Tags == nil {
		return []string{}
	}
	return pq.StringArray(t.Tags)
}
