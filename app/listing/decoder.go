package listing

import (
	"net/url"

	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
)

// DecodeLimit set the paging limit default.
func DecodeLimit(limit int) func(*Decoder) {
	return func(dec *Decoder) {
		o := paging.Limit(limit)
		dec.optionsPagingDecoder = append(dec.optionsPagingDecoder, o)
	}
}

// DecodeMaxAllowedLimit set the max allowed limit default.
func DecodeMaxAllowedLimit(maxAllowed int) func(*Decoder) {
	return func(dec *Decoder) {
		o := paging.MaxAllowedLimit(maxAllowed)
		dec.optionsPagingDecoder = append(dec.optionsPagingDecoder, o)
	}
}

// DecodeSortCriterias set criterias to sort
func DecodeSortCriterias(criterias ...sorting.Sort) func(*Decoder) {
	return func(dec *Decoder) {
		dec.sortCriterias = append(dec.sortCriterias, criterias...)
	}
}

// A Decoder reads and decodes Listing values from url.Values.
type Decoder struct {
	pagingDecoder        *paging.Decoder
	optionsPagingDecoder []paging.Option
	sortingDecoder       *sorting.Decoder
	sortCriterias        []sorting.Sort
}

// NewDecoder returns a new decoder that reads from params.
func NewDecoder(params url.Values, opts ...func(*Decoder)) *Decoder {
	d := &Decoder{}
	for _, o := range opts {
		o(d)
	}

	d.pagingDecoder = paging.NewDecoder(params, d.optionsPagingDecoder...)
	d.sortingDecoder = sorting.NewDecoder(params, d.sortCriterias...)

	return d
}

// Decode reads the Params values from url params and
// stores it in the value pointed to by v.
func (dec *Decoder) Decode(v *Listing) error {
	if err := dec.pagingDecoder.Decode(&v.Paging); err != nil {
		return err
	}

	if err := dec.sortingDecoder.Decode(&v.Sorting); err != nil {
		return err
	}

	return nil
}
