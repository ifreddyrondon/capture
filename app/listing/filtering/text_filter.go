package filtering

import "net/url"

const textFilterType = "text"

// TextDecoder validates text values and returns a Filter with the
// value found. If not filters were found, it returns nil.
type TextDecoder struct {
	FilterID
	values []FilterValue
}

// NewTextDecoder returns a new TextDecoder instance.
func NewTextDecoder(id, name string, values ...FilterValue) *TextDecoder {
	return &TextDecoder{
		FilterID: NewFilterID(id, name),
		values:   values,
	}
}

// Present gets the url params and check if a text filter is present,
// if it's present validates its value meets one of filter values options.
// Returns a Filter with the applied value or nil is not present.
func (b *TextDecoder) Present(keys url.Values) *Filter {
	for key, values := range keys {
		if key == b.ID {
			v := checkValues(b.values, values[0])
			if v != nil {
				return NewFilter(b.FilterID, textFilterType, *v)
			}
		}
	}
	return nil
}

func checkValues(availables []FilterValue, paramVal string) *FilterValue {
	for _, v := range availables {
		if v.ID == paramVal {
			return &v
		}
	}
	return nil
}

// WithValues returns the filter with all their values.
func (b *TextDecoder) WithValues() *Filter {
	return NewFilter(b.FilterID, textFilterType, b.values...)
}
