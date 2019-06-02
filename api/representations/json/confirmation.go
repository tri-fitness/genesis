package json

import (
	"encoding/json"
	"time"
	r "tri-fitness/genesis/api/representations"
)

type Confirmation struct {
	r.Representation `json:"-"`

	ID          int        `json:"id"`
	ExpiredAt   time.Time  `json:"expiredAt"`
	Type        string     `json:"type"`
	ConfirmedAt *time.Time `json:"confirmedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// AsBytes provides the representation as bytes.
func (c Confirmation) AsBytes() []byte {
	bytes, err := json.Marshal(&c)
	if err != nil {
		panic(err)
	}
	return bytes
}
