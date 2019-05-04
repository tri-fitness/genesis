package domain

import (
	u "github.com/gofrs/uuid"
)

type Author struct {
	UUID      u.UUID
	Type      AccountType
	GivenName string
	Surname   string
}
