package capture

// Path represent an array of captures.
type Path struct {
	Captures []*Capture
}

// AddCapture add new captures into the path.
// The new captures will always be added at the end of the road.
// respecting their insertion order.
func (p *Path) AddCapture(captures ...*Capture) {
	p.Captures = append(p.Captures, captures...)
}
