package user

import (
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

type StoreFilter struct {
	Email *string
	ID    *kallax.ULID
}

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save user into the database.
	Save(*pkg.User) error
	// Get a user from database
	Get(StoreFilter) (*pkg.User, error)
}

// PGStore implementation of user.Store for Postgres database.
type PGStore struct {
	db *gorm.DB
}

// NewPGStore creates a PGStore
func NewPGStore(db *gorm.DB) *PGStore {
	return &PGStore{db: db}
}

// Migrate (panic) runs schema migration.
func (p *PGStore) Migrate() {
	p.db.AutoMigrate(pkg.User{})
}

// Drop (panic) delete schema.
func (p *PGStore) Drop() {
	p.db.DropTableIfExists(pkg.User{})
}

// Save capture into the database.
func (p *PGStore) Save(user *pkg.User) error {
	return p.db.Create(user).Error
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStore) Get(storeFilter StoreFilter) (*pkg.User, error) {
	f := &pkg.User{}
	if storeFilter.ID != nil {
		f.ID = *storeFilter.ID
	}
	if storeFilter.Email != nil {
		f.Email = *storeFilter.Email
	}

	var result pkg.User
	if p.db.Where(f).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}

type MockStore struct {
	User *pkg.User
	Err  error
}

func (m *MockStore) Save(user *pkg.User) error                      { return m.Err }
func (m *MockStore) Get(storeFilter StoreFilter) (*pkg.User, error) { return m.User, m.Err }
