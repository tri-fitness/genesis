package json

import (
	"encoding/json"
	r "tri-fitness/genesis/api/representations"
)

type Error struct {
	r.Representation

	// Message represents a user friendly message.
	Message string `json:"message"`

	// DetailedMessage represents an engineer friendly message.
	DetailedMessage string `json:"detailedMessage"`
}

// AsBytes provides the representation as bytes.
func (e *Error) AsBytes() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return bytes
}
