// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package capture

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonCbca9c40EncodeGithubComIfreddyrondonGocaptureCapture(out *jwriter.Writer, in Capture) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.RawText((in.ID).MarshalText())
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"payload\":")
	if in.Payload == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v1First := true
		for v1Name, v1Value := range in.Payload {
			if !v1First {
				out.RawByte(',')
			}
			v1First = false
			out.String(string(v1Name))
			out.RawByte(':')
			if m, ok := v1Value.(easyjson.Marshaler); ok {
				m.MarshalEasyJSON(out)
			} else if m, ok := v1Value.(json.Marshaler); ok {
				out.Raw(m.MarshalJSON())
			} else {
				out.Raw(json.Marshal(v1Value))
			}
		}
		out.RawByte('}')
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"timestamp\":")
	out.Raw((in.Timestamp).MarshalJSON())
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"createdAt\":")
	out.Raw((in.CreatedAt).MarshalJSON())
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"updatedAt\":")
	out.Raw((in.UpdatedAt).MarshalJSON())
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"lat\":")
	if in.LAT == nil {
		out.RawString("null")
	} else {
		out.Float64(float64(*in.LAT))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"lng\":")
	if in.LNG == nil {
		out.RawString("null")
	} else {
		out.Float64(float64(*in.LNG))
	}
	out.RawByte('}')
}
