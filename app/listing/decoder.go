package listing

import (
	"net/url"

	"github.com/ifreddyrondon/capture/app/listing/filtering"

	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
)

// DecodeLimit set the paging limit default.
func DecodeLimit(limit int) func(*Decoder) {
	return func(dec *Decoder) {
		o := paging.Limit(limit)
		dec.pagingOpts = append(dec.pagingOpts, o)
	}
}

// DecodeMaxAllowedLimit set the max allowed limit default.
func DecodeMaxAllowedLimit(maxAllowed int) func(*Decoder) {
	return func(dec *Decoder) {
		o := paging.MaxAllowedLimit(maxAllowed)
		dec.pagingOpts = append(dec.pagingOpts, o)
	}
}

// DecodeSort set criterias to sort
func DecodeSort(criterias ...sorting.Sort) func(*Decoder) {
	return func(dec *Decoder) {
		dec.sortCriterias = append(dec.sortCriterias, criterias...)
	}
}

// DecodeFilter set criterias to filter
func DecodeFilter(criterias ...filtering.FilterDecoder) func(*Decoder) {
	return func(dec *Decoder) {
		dec.filteringCriterias = append(dec.filteringCriterias, criterias...)
	}
}

// A Decoder reads and decodes Listing values from url.Values.
type Decoder struct {
	pagingDecoder      *paging.Decoder
	pagingOpts         []paging.Option
	sortingDecoder     *sorting.Decoder
	sortCriterias      []sorting.Sort
	filteringDecoder   *filtering.Decoder
	filteringCriterias []filtering.FilterDecoder
}

// NewDecoder returns a new decoder that reads from params.
func NewDecoder(params url.Values, opts ...func(*Decoder)) *Decoder {
	d := &Decoder{}
	for _, o := range opts {
		o(d)
	}

	d.pagingDecoder = paging.NewDecoder(params, d.pagingOpts...)
	d.sortingDecoder = sorting.NewDecoder(params, d.sortCriterias...)
	d.filteringDecoder = filtering.NewDecoder(params, d.filteringCriterias...)

	return d
}

// Decode reads the Params values from url params and
// stores it in the value pointed to by v.
func (dec *Decoder) Decode(v *Listing) error {
	if err := dec.pagingDecoder.Decode(&v.Paging); err != nil {
		return err
	}

	if err := decodeSorting(dec, v); err != nil {
		return err
	}

	if err := decodeFiltering(dec, v); err != nil {
		return err
	}

	return nil
}

func decodeSorting(dec *Decoder, v *Listing) error {
	if len(dec.sortCriterias) < 1 {
		return nil
	}

	if v.Sorting == nil {
		v.Sorting = &sorting.Sorting{}
	}

	if err := dec.sortingDecoder.Decode(v.Sorting); err != nil {
		return err
	}
	return nil
}

func decodeFiltering(dec *Decoder, v *Listing) error {
	if len(dec.filteringCriterias) < 1 {
		return nil
	}

	if v.Filtering == nil {
		v.Filtering = &filtering.Filtering{}
	}

	if err := dec.filteringDecoder.Decode(v.Filtering); err != nil {
		return err
	}
	return nil
}
