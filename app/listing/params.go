package listing

import (
	"net/url"

	"github.com/ifreddyrondon/capture/app/listing/paging"
)

// Params containst the info to perform filter sort and paging over a collection.
type Params struct {
	paging.Paging
	AvailableSort []Sort
	Sort
	AvailableFilter []Filter
	Filter          Filter
}

// NewParamsDefault returns a new instance of Params with defaul values.
func NewParamsDefault() *Params {
	return &Params{
		Paging: paging.NewDefaults(),
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

// FilterValue defines a value that a Filter can have.
type FilterValue struct {
	ID, Name string
}

// Filter struct allows to filter a collection by an identifier.
type Filter struct {
	ID, Name string
	values   []FilterValue
}

// NewFilter returns a new instance of Filter.
func NewFilter(id, name string) *Filter {
	return &Filter{ID: id, Name: name}
}

// Sort struct allows to sort a collection given a sort id.
type Sort struct {
	id, name string
}

// NewSort returns a new instance of Sort
func NewSort(id, name string) *Sort {
	return &Sort{id: id, name: name}
}
