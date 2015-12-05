package schema

import (
	"encoding/json"
)

type GenericJson struct {
	Action string                        `json:"action"`
	Data   []map[string]*json.RawMessage `json:"data"`
}

func (g *GenericJson) IsMessage() bool {
	return g.Action == "message"
}

func (g *GenericJson) IsSubscribe() bool {
	return g.Action == "subscribe"
}
