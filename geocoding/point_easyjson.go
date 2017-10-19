// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package geocoding

import (
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

func easyjson3844eb60DecodeGithubComIfreddyrondonGocaptureGeocoding(in *jlexer.Lexer, out *pointJSON) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "lat":
			out.Lat = float64(in.Float64())
		case "latitude":
			out.Latitude = float64(in.Float64())
		case "lng":
			out.Lng = float64(in.Float64())
		case "longitude":
			out.Longitude = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *pointJSON) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3844eb60DecodeGithubComIfreddyrondonGocaptureGeocoding(&r, v)
	return r.Error()
}

func easyjson3844eb60EncodeGithubComIfreddyrondonGocaptureGeocoding1(out *jwriter.Writer, in Point) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"lat\":")
	out.Float64(float64(in.Lat))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"lng\":")
	out.Float64(float64(in.Lng))
	out.RawByte('}')
}