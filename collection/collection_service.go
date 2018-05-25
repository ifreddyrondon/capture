package collection

// Service is the interface to be implemented by collection services
// It's the layer between HTTP server and Repositories.
type Service interface {
	// Save a collection.
	Save(*Collection) error
}

// REPOService implementation of Service for repository
type REPOService struct {
	repository Repository
}

// NewService creates a new collection Service
func NewService(repository Repository) *REPOService {
	return &REPOService{repository: repository}
}

// Save capture into the database.
func (s *REPOService) Save(c *Collection) error {
	c.fillIfEmpty()
	return s.repository.Save(c)
}
