package schema

import (
	"encoding/json"
)

type GenericJson struct {
	Action string                        `json:"action"`
	Data   []map[string]*json.RawMessage `json:"data"`
}
