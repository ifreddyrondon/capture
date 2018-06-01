package postgres

import (
	"fmt"

	"github.com/lib/pq"
)

// IsUniqueConstraintError is a helper function to check if the an error returned is for unique constrait
func IsUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		fmt.Println(pqErr.Constraint)
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}
