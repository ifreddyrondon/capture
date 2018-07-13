package filtering

import (
	"net/url"
)

// FilterID is a helper struct with the necessaries fields to identify a filter.
// It should be used as embedded struct.
type FilterID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewFilterID returns a new FilterID instance.
func NewFilterID(id, name string) FilterID {
	return FilterID{ID: id, Name: name}
}

// FilterValue is the struct where the posibles filter values should be stored.
type FilterValue struct {
	FilterID
	Result int64 `json:"result,omitempty"`
}

// NewFilterValue returns a new FilterValue instance.
func NewFilterValue(id, name string) FilterValue {
	return FilterValue{FilterID: NewFilterID(id, name)}
}

// Filter struct that represent a filter
type Filter struct {
	FilterID
	Type   string        `json:"type"`
	Values []FilterValue `json:"values"`
}

// NewFilter returns a new Filter instance.
func NewFilter(id FilterID, typef string, values ...FilterValue) *Filter {
	return &Filter{
		FilterID: id,
		Type:     typef,
		Values:   values,
	}
}

// FilterBuilder interface to validate and returns Filter's.
type FilterBuilder interface {
	// Validate gets the url params and check if a filter is present within them,
	// if it's present validates if its value is valid.
	// Returns a Filter with the applied value or nil is not present.
	Validate(url.Values) *Filter
	// WithValues returns a filter with all their posible values.
	WithValues() *Filter
}

// Filtering allows to filter a collection with the selected Filters
// and their selected values. The Available are all the possible Filters
// with all their possible values.
type Filtering struct {
	Filters   []Filter `json:"filters,omitempty"`
	Available []Filter `json:"available,omitempty"`
}
