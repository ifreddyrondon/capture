package paging

// Paging struct allows to do pagination into a collection.
type Paging struct {
	MaxAllowedLimit int
	Limit           int
	Offset          int64
	Total           int64
}
