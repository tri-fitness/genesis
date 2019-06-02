package infrastructure

import (
	"time"

	u "github.com/gofrs/uuid"
)

type notification struct {
	UUID    *u.UUID    `json:"uuid"`
	Content string     `json:"content"`
	Status  string     `json:"status"`
	SendAt  time.Time  `json:"sendAt"`
	SentAt  *time.Time `json:"sentAt"`
	Targets []target   `json:"targets"`
}

type target struct {
	UUID        *u.UUID `json:"uuid"`
	PhoneNumber string  `json:"phoneNumber"`
	Name        string  `json:"name"`
}
