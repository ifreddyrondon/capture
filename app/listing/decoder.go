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
	params             url.Values
	pagingOpts         []paging.Option
	sortCriterias      []sorting.Sort
	filteringCriterias []filtering.FilterDecoder
}

// NewDecoder returns a new decoder that reads from params.
func NewDecoder(params url.Values, opts ...func(*Decoder)) *Decoder {
	d := &Decoder{params: params}
	for _, o := range opts {
		o(d)
	}

	return d
}

// Decode reads the Params values from url params and
// stores it in the value pointed to by v.
func (dec *Decoder) Decode(v *Listing) error {
	if err := dec.paging(v); err != nil {
		return err
	}

	if err := dec.sorting(v); err != nil {
		return err
	}

	dec.filtering(v)

	return nil
}

func (dec *Decoder) paging(v *Listing) error {
	decoder := paging.NewDecoder(dec.params, dec.pagingOpts...)
	if err := decoder.Decode(&v.Paging); err != nil {
		return err
	}
	return nil
}

func (dec *Decoder) sorting(v *Listing) error {
	if len(dec.sortCriterias) < 1 {
		return nil
	}

	if v.Sorting == nil {
		v.Sorting = &sorting.Sorting{}
	}

	decoder := sorting.NewDecoder(dec.params, dec.sortCriterias...)
	if err := decoder.Decode(v.Sorting); err != nil {
		return err
	}
	return nil
}

func (dec *Decoder) filtering(v *Listing) {
	if len(dec.filteringCriterias) < 1 {
		return
	}

	if v.Filtering == nil {
		v.Filtering = &filtering.Filtering{}
	}

	decoder := filtering.NewDecoder(dec.params, dec.filteringCriterias...)
	decoder.Decode(v.Filtering)
}
