package filtering

// import (
// 	"net/url"
// )

// const errFilterValueNotAvailable = "the filter %v doesn't have the value %v"

// // A Decoder reads and decodes Filtering values from url.Values.
// type Decoder struct {
// 	params url.Values

// 	available []Filter
// }

// // NewDecoder returns a new decoder that reads from query params.
// func NewDecoder(params url.Values, available ...Filter) *Decoder {
// 	d := &Decoder{params: params, available: available}

// 	return d
// }

// // Decode reads the filter-encoded values from params and stores it
// // in the value pointed to by v. If a value is missing from the params
// // it'll be filled by their equivalent default value.
// func (dec *Decoder) Decode(v *Filtering) error {
// 	dec.fillDefaults(v)

// 	// if ok {
// 	// 	f := paramsInAvailables(filterStr, dec.available)
// 	// 	if f == nil {
// 	// 		return fmt.Errorf(errSortKeyNotAvailable, sortStr[0])
// 	// 	}
// 	// 	// 	v.Sort.ID = sort.ID
// 	// 	// 	v.Sort.Name = sort.Name
// 	// }
// 	return nil
// }

// func paramsInAvailables(filterKeys url.Values, availables []Filter) []Filter {
// 	filters := []Filter{}
// 	for key, values := range filterKeys {
// 		for _, f := range availables {
// 			if key == f.ID {
// 				filters = append(filters, valuesInParams(values, f))
// 				break
// 			}
// 		}
// 	}

// 	return filters
// }

// func valuesInParams(values []string, filter Filter) Filter {
// 	for _, val := range values
// }

// func (dec *Decoder) fillDefaults(s *Filtering) {
// 	s.Available = dec.available
// }
