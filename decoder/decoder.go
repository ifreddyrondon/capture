package decoder

// OK represents types capable of validating themselves
type OK interface {
	OK() error
}
