package repository

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/features"
	"github.com/jinzhu/gorm"
)

// Store is the interface to be implemented by any kind of store
// It make CRUD operations over a store.
type Store interface {
	// Save a repository.
	Save(*features.Repository) error
	// List retrieve repositories from start index to count.
	List(ListingRepo) ([]features.Repository, error)
}

type ListingRepo struct {
	SortKey    string
	Visibility *features.Visibility
	Offset     int64
	Limit      int
}

func NewListingRepo(l listing.Listing) ListingRepo {
	lrepo := ListingRepo{
		SortKey: l.Sorting.Sort.Value,
		Offset:  l.Paging.Offset,
		Limit:   l.Paging.Limit,
	}

	for i := range l.Filtering.Filters {
		if l.Filtering.Filters[i].ID == "visibility" {
			visibility := features.Visibility(l.Filtering.Filters[i].Values[0].ID)
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
	pgs.db.AutoMigrate(features.Repository{})
}

// Drop (panic) delete the repository schema.
func (pgs *PGStore) Drop() {
	pgs.db.DropTableIfExists(features.Repository{})
}

// Save a repository into the database.
func (pgs *PGStore) Save(r *features.Repository) error {
	return pgs.db.Create(r).Error
}

func (pgs *PGStore) List(l ListingRepo) ([]features.Repository, error) {
	var results []features.Repository
	q := pgs.db
	if l.Visibility != nil {
		q = pgs.db.Where(&features.Repository{Visibility: *l.Visibility})
	}
	q = q.Order(l.SortKey).Offset(l.Offset).Limit(l.Limit)
	if err := q.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
