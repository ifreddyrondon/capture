package filtering

import "net/url"

const (
	trueID            = "true"
	falseID           = "false"
	booleanFilterType = "boolean"
)

// BooleanDecoder validates boolean values and returns a Filter.
// It returns a Filter with two posible values true or false.
type BooleanDecoder struct {
	FilterID
	trueFilterValue  FilterValue
	falseFilterValue FilterValue
}

// NewBooleanDecoder returns a new BooleanDecoder instance.
func NewBooleanDecoder(id, name string, trueValName, falseValName string) *BooleanDecoder {
	return &BooleanDecoder{
		FilterID:         NewFilterID(id, name),
		trueFilterValue:  NewFilterValue(trueID, trueValName),
		falseFilterValue: NewFilterValue(falseID, falseValName),
	}
}

// Present gets the url params and check if a boolean filter is present,
// if it's present validates if its value are true or false.
// Returns a Filter with the applied value or nil is not present.
func (b *BooleanDecoder) Present(keys url.Values) *Filter {
	for key, values := range keys {
		if key == b.ID {
			v := values[0]
			if v == b.trueFilterValue.ID {
				return NewFilter(b.FilterID, booleanFilterType, b.trueFilterValue)
			}
			if v == b.falseFilterValue.ID {
				return NewFilter(b.FilterID, booleanFilterType, b.falseFilterValue)
			}
		}
	}
	return nil
}

// WithValues returns the filter with true and false values.
func (b *BooleanDecoder) WithValues() *Filter {
	return NewFilter(b.FilterID, booleanFilterType, b.trueFilterValue, b.falseFilterValue)
}
