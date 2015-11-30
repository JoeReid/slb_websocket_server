package server

import (
	"encoding/json"
)

type genericJson struct {
	Action string                        `json:"action"`
	Data   []map[string]*json.RawMessage `json:"data"`
}
