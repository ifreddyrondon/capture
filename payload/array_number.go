package payload

import (
	"encoding/json"
	"log"

	"github.com/mailru/easyjson/jlexer"
)

type unmarshalMap struct {
	Cap      []float64 `json:"cap"`
	Captures []float64 `json:"captures"`
	Data     []float64 `json:"data"`
	Payload  []float64 `json:"payload"`
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *unmarshalMap) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComIfreddyrondonGocapturePayload(&r, v)
	return r.Error()
}

func (v *unmarshalMap) getPayload() []float64 {
	if v.Cap != nil {
		return v.Cap
	} else if v.Captures != nil {
		return v.Captures
	} else if v.Data != nil {
		return v.Data
	}
	return v.Payload
}

// ArrayNumberPayload represent an association of float numbers
type ArrayNumberPayload []float64

func (pp ArrayNumberPayload) Values() interface{} {
	return pp
}

func (pp *ArrayNumberPayload) UnmarshalJSON(data []byte) error {
	model := new(unmarshalMap)
	if err := json.Unmarshal(data, model); err != nil {
		log.Print(err)
		return ErrorUnmarshalPayload
	}
	*pp = model.getPayload()
	return nil
}
