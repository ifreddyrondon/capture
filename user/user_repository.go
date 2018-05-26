package user

import "github.com/jinzhu/gorm"

// Repository is the interface to be implemented by collection
// It make CRUD operations over a repository.
type Repository interface {
	// Save user into the database.
	Save(*User) error
	// Get a user by email from database
	Get(string) (*User, error)
}

// PGRepository implementation of user.Repository for Postgres database.
type PGRepository struct {
	db *gorm.DB
}

// NewPGRepository creates a PGRepository
func NewPGRepository(db *gorm.DB) *PGRepository {
	return &PGRepository{db: db}
}

// Migrate (panic) runs schema migration.
func (p *PGRepository) Migrate() {
	p.db.AutoMigrate(User{})
}

// Drop (panic) delete schema.
func (p *PGRepository) Drop() {
	p.db.DropTableIfExists(User{})
}

// Save capture into the database.
func (p *PGRepository) Save(user *User) error {
	return p.db.Create(user).Error
}

// Get a user by email, if not found user returns an error
func (p *PGRepository) Get(email string) (*User, error) {
	var result User
	if p.db.Where(&User{Email: email}).First(&result).RecordNotFound() {
		return nil, ErrNotFound
	}
	return &result, nil
}
