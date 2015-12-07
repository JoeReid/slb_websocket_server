package schema

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
)

var (
	NotMessageError   = errors.New("Json is not a message type")
	NotSubscribeError = errors.New("Json is not a subscribe type")
)

type MessageJson struct {
	Messages []map[string]*json.RawMessage `json:"data"`
}

func (m *MessageJson) ToMessageArray() ([]SingleMessage, error) {
	var rtn []SingleMessage

	for _, elm := range m.Messages {
		var msg SingleMessage

		group, ok := elm["deviceid"]
		if ok {
			var gStr string

			err := json.Unmarshal(*group, &gStr)
			if err != nil {
				// bad json
				log.WithFields(log.Fields{
					"json": elm,
				}).Warn("Dropping malformed json")

				continue
			}

			msg.Group = gStr
		}

		msg.WholeMessage = elm

		rtn = append(rtn, msg)
	}
	return rtn, nil
}

type SingleMessage struct {
	Group        string
	WholeMessage map[string]*json.RawMessage
	ConnectionID int
}

type SubscribeJson struct {
	Groups []string `json:"data"`
}

type GenericJson struct {
	Action string           `json:"action"`
	Data   *json.RawMessage `json:"data"`
}

func (g *GenericJson) ToMessage() (MessageJson, error) {
	var rtn MessageJson

	if g.Action != "message" {
		return rtn, NotMessageError
	}

	err := json.Unmarshal(*g.Data, &rtn.Messages)
	if err != nil {
		return rtn, err
	}
	return rtn, nil
}

func (g *GenericJson) ToSubscribe() (SubscribeJson, error) {
	var rtn SubscribeJson

	if g.Action != "subscribe" {
		return rtn, NotSubscribeError
	}

	err := json.Unmarshal(*g.Data, &rtn.Groups)
	if err != nil {
		return rtn, err
	}
	return rtn, nil
}
