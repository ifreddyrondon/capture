package paging

// Paging struct allows to do pagination into a collection.
type Paging struct {
	MaxAllowedLimit int   `json:"max_allowed_limit"`
	Limit           int   `json:"limit"`
	Offset          int64 `json:"offset"`
	Total           int64 `json:"total,omitempty"`
}
