package domain

import (
	"errors"
	"time"

	u "github.com/gofrs/uuid"
)

type AccountType int

const (
	AccountTypeStudent AccountType = iota
	AccountTypeInstructor
)

func (t AccountType) String() string {
	return []string{
		"STUDENT",
		"INSTRUCTOR",
	}[t]
}

func AccountTypeFromString(str string) (AccountType, error) {
	valid := map[string]AccountType{
		"STUDENT":    AccountTypeStudent,
		"INSTRUCTOR": AccountTypeInstructor,
	}
	if t, ok := valid[str]; ok {
		return t, nil
	}
	return AccountType(-1), errors.New("invalid account type")
}

type Account struct {
	UUID                u.UUID
	PrimaryCredential   string
	SecondaryCredential string
	Type                AccountType
	GivenName           string
	Surname             string
	Bio                 string
	Email               string
	Phone               string
	Subscriptions       []Subscription
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
