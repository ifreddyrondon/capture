package paging_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaults(t *testing.T) {
	result := paging.NewDefaults()
	assert.NotNil(t, result)
	assert.Equal(t, int64(10), result.Limit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, int64(0), result.Total)
}

func TestDecode(t *testing.T) {
	tt := []struct {
		name     string
		params   url.Values
		defaults paging.Paging
		result   paging.Paging
	}{
		{
			"decode with new limit and default offset from paging.NewDefaults",
			map[string][]string{"limit": []string{"1"}},
			paging.NewDefaults(),
			func() paging.Paging {
				p := paging.NewDefaults()
				p.Limit = 1
				return p
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var p paging.Paging
			p.Decode(tc.params, tc.defaults)
			assert.Equal(t, p.Limit, tc.result.Limit)
			assert.Equal(t, p.Offset, tc.result.Offset)
		})
	}
}
