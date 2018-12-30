package repo

import (
	"github.com/go-pg/pg/orm"
	"github.com/ifreddyrondon/capture/pkg/domain"
)

type filter domain.Listing

func (f *filter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.Owner != nil {
		q = q.Where("user_id = ?", *f.Owner)
	}
	if f.Visibility != nil {
		q = q.Where("visibility = ?", *f.Visibility)
	}

	return q.Order(f.SortKey).
		Offset(int(f.Offset)).
		Limit(f.Limit), nil
}
