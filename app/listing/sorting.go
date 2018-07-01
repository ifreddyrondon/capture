package listing

import (
	"fmt"
	"net/url"
)

const errSortKeyNotAvailable = "There's no order criteria with the id %v"

// Sort criteria.
type Sort struct {
	ID, Name string
}

// NewSort returns a new instance of Sort
func NewSort(id, name string) Sort {
	return Sort{ID: id, Name: name}
}

// Sorting struct allows to sort a collection.
type Sorting struct {
	Sort      Sort
	Available []Sort
}

// NewSorting returns a new instance of Sorting
// with the sort criteria args as availables.
//
// The first sort criteria if present will be
// the default Sort when decode url query params
// and not params present.
func NewSorting(availables ...Sort) Sorting {
	s := Sorting{Available: availables}
	if len(availables) > 0 {
		s.Sort = availables[0]
	}

	return s
}

// Decode gets url.Values (with query params) and fill Sort and
// Available with these. If a value is missing from the params then
// it'll be filled by their equivalent default value.
func (s *Sorting) Decode(params url.Values, defaults Sorting) error {
	s.Available = defaults.Available
	s.Sort.ID = defaults.Sort.ID
	s.Sort.Name = defaults.Sort.Name

	sortStr, ok := params["sort"]
	if ok {
		sort := paramsInAvailables(sortStr[0], defaults.Available)
		if sort == nil {
			return fmt.Errorf(errSortKeyNotAvailable, sortStr[0])
		}
		s.Sort.ID = sort.ID
		s.Sort.Name = sort.Name
	}
	return nil
}

func paramsInAvailables(sortKey string, availables []Sort) *Sort {
	for _, sort := range availables {
		if sortKey == sort.ID {
			return &sort
		}
	}
	return nil
}
