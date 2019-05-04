package json

import (
	"encoding/json"
	"time"
	r "tri-fitness/genesis/api/representations"

	u "github.com/gofrs/uuid"
)

type Exercise struct {
	r.Representation

	UUID       u.UUID
	VideoID    int
	Difficulty string
	Watches    int
	Duration   time.Duration
	Author     u.UUID
}

// AsBytes provides the representation as bytes.
func (e Exercise) AsBytes() []byte {
	bytes, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}
	return bytes
}
