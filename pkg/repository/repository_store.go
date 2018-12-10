package repository

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save a repository.
	Save(*pkg.User, *pkg.Repository) error
	// List retrieve repositories from start index to count.
	List(ListingRepo) ([]pkg.Repository, error)
	// Get a repo by id
	Get(string) (*pkg.Repository, error)
	// Drop register if it is exist
	Drop()
}

type ListingRepo struct {
	SortKey    string
	Visibility *pkg.Visibility
	Owner      *pkg.User
	Offset     int64
	Limit      int
}

func newListingRepo(l listing.Listing) ListingRepo {
	lrepo := ListingRepo{
		SortKey: l.Sorting.Sort.Value,
		Offset:  l.Paging.Offset,
		Limit:   l.Paging.Limit,
	}

	if l.Filtering == nil {
		return lrepo
	}

	for i := range l.Filtering.Filters {
		if l.Filtering.Filters[i].ID == "visibility" {
			visibility := pkg.Visibility(l.Filtering.Filters[i].Values[0].ID)
			lrepo.Visibility = &visibility
			break
		}
	}

	return lrepo
}

// PGStore implementation of repository.Store for Postgres database.
type PGStore struct {
	db *gorm.DB
}

// NewPGStore creates a PGStore
func NewPGStore(db *gorm.DB) *PGStore {
	return &PGStore{db: db}
}

// Migrate (panic) runs schema migration for repository table.
func (pgs *PGStore) Migrate() {
	pgs.db.AutoMigrate(pkg.Repository{})
}

// Drop (panic) delete the repository schema.
func (pgs *PGStore) Drop() {
	pgs.db.DropTableIfExists(pkg.Repository{})
}

// Save a repository into the database.
func (pgs *PGStore) Save(owner *pkg.User, r *pkg.Repository) error {
	// FIXME: handler err
	id, _ := kallax.NewULIDFromText(owner.ID)
	r.UserID = id
	return pgs.db.Create(r).Error
}

func (pgs *PGStore) List(l ListingRepo) ([]pkg.Repository, error) {
	var results []pkg.Repository
	f := &pkg.Repository{}
	if l.Owner != nil {
		// FIXME: handler err
		id, _ := kallax.NewULIDFromText(l.Owner.ID)
		f.UserID = id
	}
	if l.Visibility != nil {
		f.Visibility = *l.Visibility
	}
	err := pgs.db.
		Where(f).
		Order(l.SortKey).
		Offset(l.Offset).
		Limit(l.Limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (pgs *PGStore) Get(idStr string) (*pkg.Repository, error) {
	var result pkg.Repository
	// FIXME: handler id err
	id, _ := kallax.NewULIDFromText(idStr)
	if pgs.db.Where(&pkg.Repository{ID: id}).First(&result).RecordNotFound() {
		return nil, ErrorNotFound
	}
	return &result, nil
}
