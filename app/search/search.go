package search

// Filter struct allows to filter a collection given a filter id.
type Filter struct {
	id, name string
}

// NewFilter returns a new instance of Filter
func NewFilter(id, name string) *Filter {
	return &Filter{id: id, name: name}
}

// Sort struct allows to sort a collection given a sort id.
type Sort struct {
	id, name string
}

// NewSort returns a new instance of Sort
func NewSort(id, name string) *Sort {
	return &Sort{id: id, name: name}
}

// Paging struct allows to do pagination into a collection.
type Paging struct {
	total, offset, limit int
}

// NewPaging returns a new instance of Paging
func NewPaging(total, offset, limit int) *Paging {
	return &Paging{total: total, offset: offset, limit: limit}
}

// Params struct containst the data to perform filter sort and paging over a collection.
type Params struct {
	Paging
	AvailableSort []Sort
	Sort
	AvailableFilter []Filter
	Filter          Filter
}
