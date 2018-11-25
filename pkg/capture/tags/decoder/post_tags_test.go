package decoder_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/capture/tags/decoder"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDecodePostTagsOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		body     string
		expected decoder.PostTags
	}{
		{
			"decode timestamp without data ({})",
			`{}`,
			decoder.PostTags{Tags: nil},
		},
		{
			"decode timestamp without data ({tags:[]})",
			`{"tags":[]}`,
			decoder.PostTags{Tags: []string{}},
		},
		{
			"decode timestamp with data",
			`{"tags":["at night", "ism"]}`,
			decoder.PostTags{Tags: []string{"at night", "ism"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var tags decoder.PostTags
			err := decoder.Decode(r, &tags)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, tags)
		})
	}
}

func TestDecodePostTagsError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			"decode timestamp when invalid json",
			".",
			"cannot unmarshal json into tags value",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var tags decoder.PostTags
			err := decoder.Decode(r, &tags)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestTagsFromPostTagsOK(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		postTags decoder.PostTags
		expected pq.StringArray
	}{
		{
			"get pq.StringArray from PostTags when empty",
			decoder.PostTags{Tags: nil},
			pq.StringArray{},
		},
		{
			"get pq.StringArray from PostTags with data",
			decoder.PostTags{Tags: []string{"at night"}},
			pq.StringArray{"at night"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			date := tc.postTags.GetTags()
			assert.Equal(t, tc.expected, date)
		})
	}
}
