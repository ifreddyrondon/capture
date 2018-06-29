package listing

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	defaultMaxAllowedLimit = 100
	defaultLimit           = 10
	defaultOffset          = 0
)

var (
	// ErrInvalidOffsetValueNotANumber expected error when fails parsing the offset value to int.
	ErrInvalidOffsetValueNotANumber = errors.New("invalid offset value, must be a number")
	// ErrInvalidOffsetValueLessThanZero expected error when offset value is less than zero.
	ErrInvalidOffsetValueLessThanZero = errors.New("invalid offset value, must be greater than zero")
	// ErrInvalidLimitValueNotANumber expected error when fails parse limit value to int
	ErrInvalidLimitValueNotANumber = errors.New("invalid limit value, must be a number")
	// ErrInvalidLimitValueLessThanZero expected error when limit value is less than zero.
	ErrInvalidLimitValueLessThanZero = errors.New("invalid limit value, must be greater than zero")
)

// Paging struct allows to do pagination into a collection.
type Paging struct {
	maxAllowedLimit int
	Limit           int
	Total, Offset   int64
}

// NewPaging returns a new instance of Paging with defaults values.
func NewPaging() Paging {
	p := Paging{
		maxAllowedLimit: defaultMaxAllowedLimit,
		Offset:          defaultOffset,
		Limit:           defaultLimit,
	}

	return p
}

// Decode gets url.Values (with query params) and fill the Paging
// instance with these. If a value is missing from the params then
// it'll be filled by their equivalent default value.
func (p *Paging) Decode(params url.Values, defaults Paging) error {
	p.Offset = defaults.Offset
	p.Limit = defaults.Limit

	offsetStr, ok := params["offset"]
	if ok {
		off, err := strconv.ParseInt(offsetStr[0], 10, 64)
		if err != nil {
			return ErrInvalidOffsetValueNotANumber
		}
		if off < 0 {
			return ErrInvalidOffsetValueLessThanZero
		}
		p.Offset = off
	}
	limitStr, ok := params["limit"]
	if ok {
		l, err := strconv.Atoi(limitStr[0])
		if err != nil {
			return ErrInvalidLimitValueNotANumber
		}
		if l < 0 {
			return ErrInvalidLimitValueLessThanZero
		}
		if l > defaults.maxAllowedLimit {
			l = defaults.maxAllowedLimit
		}
		p.Limit = l
	}
	return nil
}
