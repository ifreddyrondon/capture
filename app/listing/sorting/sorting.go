package sorting

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
