package listing

import "net/url"

// Sort struct allows to sort a collection.
type Sort struct {
	id, name string
}

// Decode gets url.Values (with query params) and fill the Sort
// instance with these. If a value is missing from the params then
// it'll be filled by their equivalent default value.
func (s *Sort) Decode(params url.Values, defaults Sort) error {
	return nil
}
