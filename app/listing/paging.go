package listing

import (
	"net/url"
	"strconv"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

// Paging struct allows to do pagination into a collection.
type Paging struct {
	Total, Offset, Limit int64
}

func newPagingDefault() Paging {
	return Paging{Offset: defaultOffset, Limit: defaultLimit}
}

func (p *Paging) decode(params url.Values, defaults Paging) error {
	var err error
	offsetStr, ok := params["offset"]
	if ok {
		p.Offset, err = strconv.ParseInt(offsetStr[0], 10, 64)
		if err != nil {
			return err
		}
	} else {
		p.Offset = defaults.Offset
	}
	limitStr, ok := params["limit"]
	if ok {
		p.Limit, err = strconv.ParseInt(limitStr[0], 10, 64)
		if err != nil {
			return err
		}
	} else {
		p.Limit = defaults.Limit
	}
	return err
}
