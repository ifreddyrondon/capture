package domain_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestAllowedVisibility(t *testing.T) {
	tt := []struct {
		name     string
		given    string
		expected bool
	}{
		{
			"empty visibility",
			"",
			false,
		},
		{
			"not allowed visibility",
			"protected",
			false,
		},
		{
			"public visibility",
			"public",
			true,
		},
		{
			"private visibility",
			"private",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.AllowedVisibility(tc.given)
			assert.Equal(t, tc.expected, result)
		})
	}
}
