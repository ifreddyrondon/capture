package timestamp

import "encoding/json"

type jsonTimestamp struct {
	Date      json.Number `json:"date"`
	Timestamp json.Number `json:"timestamp"`
}
