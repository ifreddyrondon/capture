package timestamp

import "encoding/json"

type timestampJSON struct {
	Date      json.Number `json:"date"`
	Timestamp json.Number `json:"timestamp"`
}
