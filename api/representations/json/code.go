package json

import (
	"encoding/json"
	r "tri-fitness/genesis/api/representations"
)

type Code struct {
	r.Representation `json:"-"`

	Code int `json:"code"`
}

// AsBytes provides the representation as bytes.
func (c Code) AsBytes() []byte {
	bytes, err := json.Marshal(&c)
	if err != nil {
		panic(err)
	}
	return bytes
}
