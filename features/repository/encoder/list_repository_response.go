package encoder

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/features"
)

type ListRepositoryResponse struct {
	Results []features.Repository `json:"results"`
	Listing *listing.Listing      `json:"listing"`
}
