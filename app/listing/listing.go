package listing

import (
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// Listing containst the info to perform filter sort and paging over a collection.
type Listing struct {
	Paging paging.Paging `json:"paging"`
	sorting.Sorting
}

// MarshalJSON supports json.Marshaler interface
func (v Listing) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDe046902EncodeGithubComIfreddyrondonCaptureAppListing(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
