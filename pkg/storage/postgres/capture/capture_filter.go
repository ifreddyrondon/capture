package capture

import (
	"github.com/go-pg/pg/orm"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

type filter domain.Listing

func (f *filter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.Owner != nil {
		q = q.Where("repository_id = ?", *f.Owner)
	}
	return q.Order(f.SortKey).
		Offset(int(f.Offset)).
		Limit(f.Limit), nil
}
