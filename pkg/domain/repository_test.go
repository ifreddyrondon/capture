package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestMarshalJSONRepository(t *testing.T) {
	t.Parallel()

	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := `{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","name":"test","current_branch":"","visibility":"public","createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","owner":"00000000-0000-0000-0000-000000000000"}`
	c := domain.Repository{
		Name: "test",
		ID: func() kallax.ULID {
			id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			return id
		}(),
		Visibility: domain.Public,
		CreatedAt:  d,
		UpdatedAt:  d,
	}

	result, err := json.Marshal(c)
	require.Nil(t, err)
	assert.Equal(t, expected, string(result))
}
