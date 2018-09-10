package user

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save user into the database.
	Save(*User) error
	// Get a user by email from database
	GetByEmail(string) (*User, error)
	// Get a user by id from database
	GetByID(kallax.ULID) (*User, error)
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
	p.db.AutoMigrate(User{})
}

// Drop (panic) delete schema.
func (p *PGStore) Drop() {
	p.db.DropTableIfExists(User{})
}

// Save capture into the database.
func (p *PGStore) Save(user *User) error {
	return p.db.Create(user).Error
}

// GetByEmail a user by email, if not found returns an error
func (p *PGStore) GetByEmail(email string) (*User, error) {
	var result User
	if p.db.Where(&User{Email: email}).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}

// GetByID a user by ID, if not found returns an error
func (p *PGStore) GetByID(id kallax.ULID) (*User, error) {
	var result User
	if p.db.Where(&User{ID: id}).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}

type MockStore struct {
	User *User
	Err  error
}

func (m *MockStore) Save(user *User) error                  { return m.Err }
func (m *MockStore) GetByEmail(email string) (*User, error) { return m.User, m.Err }
func (m *MockStore) GetByID(id kallax.ULID) (*User, error)  { return m.User, m.Err }
