package listing_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOK(t *testing.T) {
	tt := []struct {
		name      string
		urlParams url.Values
		opts      []func(*listing.Decoder)
		result    listing.Listing
	}{
		{
			"given none query params and non options should decode paging with defaults",
			map[string][]string{},
			[]func(*listing.Decoder){},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given none query params and limit option should decode paging defaults with new limit",
			map[string][]string{},
			[]func(*listing.Decoder){listing.DecodeLimit(50)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           50,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given offset and limit when limit > maxAllowed default and maxAllowed option should decode paging with offset and limit upper the default",
			map[string][]string{"offset": []string{"1"}, "limit": []string{"105"}},
			[]func(*listing.Decoder){listing.DecodeMaxAllowedLimit(110)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           105,
						Offset:          1,
						MaxAllowedLimit: 110,
					},
				}
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var l listing.Listing
			err := listing.NewDecoder(tc.urlParams, tc.opts...).Decode(&l)
			assert.Nil(t, err)
			assert.Equal(t, l.Paging.Limit, tc.result.Paging.Limit)
			assert.Equal(t, l.Paging.Offset, tc.result.Paging.Offset)
			assert.Equal(t, l.Paging.MaxAllowedLimit, tc.result.Paging.MaxAllowedLimit)
		})
	}
}

func TestDecodeFails(t *testing.T) {
	tt := []struct {
		name      string
		urlParams url.Values
		err       string
	}{
		{
			"given a not number limit param should return an error when decode paging",
			map[string][]string{"limit": []string{"a"}},
			"invalid limit value, must be a number",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var l listing.Listing
			err := listing.NewDecoder(tc.urlParams).Decode(&l)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
