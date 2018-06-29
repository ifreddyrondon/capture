package listing

// Option are function to modify the defaults listing values
type Option func(*Listing)

// Limit set the paging limit default for the middleware.
func Limit(limit int) Option {
	return func(l *Listing) {
		l.defautls.Paging.Limit = limit
	}
}

// MaxAllowedLimit set the max allowed limit paging default for the middleware.
func MaxAllowedLimit(maxAllowed int) Option {
	return func(l *Listing) {
		l.defautls.Paging.maxAllowedLimit = maxAllowed
	}
}
