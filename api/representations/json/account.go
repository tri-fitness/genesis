package json

import (
	"encoding/json"
	r "tri-fitness/genesis/api/representations"

	u "github.com/gofrs/uuid"
)

type Account struct {
	r.Representation `json:"-"`

	UUID                u.UUID         `json:"uuid"`
	PrimaryCredential   string         `json:"primaryCredential"`
	SecondaryCredential string         `json:"secondaryCredential"`
	Type                string         `json:"type"`
	GivenName           string         `json:"givenName"`
	Surname             string         `json:"surname"`
	Bio                 string         `json:"bio"`
	Email               string         `json:"emailAddress"`
	Phone               string         `json:"phoneNumber"`
	Confirmations       []Confirmation `json:"confirmations"`
}

// AsBytes provides the representation as bytes.
func (a Account) AsBytes() []byte {
	bytes, err := json.Marshal(&a)
	if err != nil {
		panic(err)
	}
	return bytes
}
