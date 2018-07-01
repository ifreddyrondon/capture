package listing

import (
	"net/url"
)

// Params containst the info to perform filter sort and paging over a collection.
type Params struct {
	Paging Paging
	Sorting
	// AvailableFilter []Filter
	// Filter          Filter
}

// NewParams returns a new instance of Params with defaul values.
func NewParams() *Params {
	return &Params{
		Paging: NewPaging(),
	}
}

// Decode gets url.Values (with query params) and fill the Params
// instance with these. If a value is missing from the params then
// it'll be filled by their equivalent default value.
func (p *Params) Decode(params url.Values, defaults Params) error {
	if err := p.Paging.Decode(params, defaults.Paging); err != nil {
		return err
	}

	return nil
}

// // FilterValue defines a value that a Filter can have.
// type FilterValue struct {
// 	ID, Name string
// }

// // Filter struct allows to filter a collection by an identifier.
// type Filter struct {
// 	ID, Name string
// 	values   []FilterValue
// }
