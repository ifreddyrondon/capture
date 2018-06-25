package listing

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
