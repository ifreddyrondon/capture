package listing_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/stretchr/testify/assert"
)

func TestNewPaging(t *testing.T) {
	result := listing.NewPaging()
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 100, result.MaxAllowedLimit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, int64(0), result.Total)
}

func TestPagingDecodeOK(t *testing.T) {
	tt := []struct {
		name   string
		params url.Values
		result listing.Paging
	}{
		{
			"decode with no params and defaults from paging.NewDefaults",
			map[string][]string{},
			listing.NewPaging(),
		},
		{
			"decode with new limit and default offset from paging.NewDefaults",
			map[string][]string{"limit": []string{"1"}},
			func() listing.Paging {
				p := listing.NewPaging()
				p.Limit = 1
				return p
			}(),
		},
		{
			"decode with new offset and default limit from paging.NewDefaults",
			map[string][]string{"offset": []string{"1"}},
			func() listing.Paging {
				p := listing.NewPaging()
				p.Offset = 1
				return p
			}(),
		},
		{
			"decode with new offset and limit",
			map[string][]string{"offset": []string{"1"}, "limit": []string{"1"}},
			func() listing.Paging {
				p := listing.NewPaging()
				p.Offset = 1
				p.Limit = 1
				return p
			}(),
		},
		{
			"decode with new offset and limit > maxAllowed should return limit = maxAllowed",
			map[string][]string{"offset": []string{"1"}, "limit": []string{"101"}},
			func() listing.Paging {
				p := listing.NewPaging()
				p.Offset = 1
				p.Limit = 100
				return p
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var p listing.Paging
			err := p.Decode(tc.params, listing.NewPaging())
			assert.Nil(t, err)
			assert.Equal(t, p.Limit, tc.result.Limit)
			assert.Equal(t, p.Offset, tc.result.Offset)
			assert.Equal(t, p.MaxAllowedLimit, tc.result.MaxAllowedLimit)
		})
	}
}

func TestPagingDecodeBad(t *testing.T) {
	tt := []struct {
		name   string
		params url.Values
		err    string
	}{
		{
			"decode with invalid limit",
			map[string][]string{"limit": []string{"a"}},
			"invalid limit value, must be a number",
		},
		{
			"decode with invalid limit",
			map[string][]string{"limit": []string{"-1"}},
			"invalid limit value, must be greater than zero",
		},
		{
			"decode with invalid offset (not a number)",
			map[string][]string{"offset": []string{"a"}},
			"invalid offset value, must be a number",
		},
		{
			"decode with invalid offset (less than 0)",
			map[string][]string{"offset": []string{"-1"}},
			"invalid offset value, must be greater than zero",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var p listing.Paging
			err := p.Decode(tc.params, listing.NewPaging())
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
