package encoder

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg"
)

type ListRepositoryResponse struct {
	Results []pkg.Repository `json:"results"`
	Listing *listing.Listing `json:"listing"`
}
