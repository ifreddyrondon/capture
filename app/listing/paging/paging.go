package paging

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

var (
	// ErrInvalidOffsetValue expected error when fails parse offset value to int
	ErrInvalidOffsetValue = errors.New("invalid offset value")
	// ErrInvalidLimitValue expected error when fails parse limit value to int
	ErrInvalidLimitValue = errors.New("invalid limit value")
)

// Paging struct allows to do pagination into a collection.
type Paging struct {
	Total, Offset, Limit int64
}

// NewDefaults returns a new instance of Paging with defaults values.
func NewDefaults() Paging {
	return Paging{Offset: defaultOffset, Limit: defaultLimit}
}

// Decode gets url.Values (with query params) and fill the Paging
// instance with these. If a value is missing from the params then
// it'll be filled by their equivalent default value.
func (p *Paging) Decode(params url.Values, defaults Paging) error {
	var err error
	offsetStr, ok := params["offset"]
	if ok {
		p.Offset, err = strconv.ParseInt(offsetStr[0], 10, 64)
		if err != nil {
			return ErrInvalidOffsetValue
		}
	} else {
		p.Offset = defaults.Offset
	}
	limitStr, ok := params["limit"]
	if ok {
		p.Limit, err = strconv.ParseInt(limitStr[0], 10, 64)
		if err != nil {
			return ErrInvalidLimitValue
		}
	} else {
		p.Limit = defaults.Limit
	}
	return err
}
