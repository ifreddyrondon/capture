package filtering

import "net/url"

const (
	trueID            = "true"
	falseID           = "false"
	booleanFilterType = "boolean"
)

// BooleanBuilder validates boolean values and return the Filter.
// It returns Filter with to posible values true and false.
type BooleanBuilder struct {
	FilterID
	trueFilterValue  FilterValue
	falseFilterValue FilterValue
}

// NewBooleanBuilder returns a new BooleanBuilder instance.
func NewBooleanBuilder(id, name string, trueValName, falseValName string) *BooleanBuilder {
	return &BooleanBuilder{
		FilterID:         NewFilterID(id, name),
		trueFilterValue:  NewFilterValue(trueID, trueValName),
		falseFilterValue: NewFilterValue(falseID, falseValName),
	}
}

// Validate gets the url params and check if a boolean filter is present,
// if it's present validates if its value are true or false.
// Returns a Filter with the applied value or nil is not present.
func (b *BooleanBuilder) Validate(keys url.Values) *Filter {
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
func (b *BooleanBuilder) WithValues() *Filter {
	return NewFilter(b.FilterID, booleanFilterType, b.trueFilterValue, b.falseFilterValue)
}
