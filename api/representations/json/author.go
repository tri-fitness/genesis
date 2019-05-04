package json

import (
	"encoding/json"

	u "github.com/gofrs/uuid"
)

type Author struct {
	UUID      u.UUID
	Type      string
	GivenName string
	Surname   string
}

// AsBytes provides the representation as bytes.
func (a *Author) AsBytes() []byte {
	bytes, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return bytes
}
