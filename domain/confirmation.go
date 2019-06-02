package domain

import (
	"errors"
	"time"

	u "github.com/gofrs/uuid"
)

type ConfirmationType int

const (
	ConfirmationTypePhoneNumber = iota
	ConfirmationTypeEmailAddress
)

func (t ConfirmationType) String() string {
	return []string{
		"PHONE_NUMBER",
		"EMAIL_ADDRESS",
	}[t]
}

func NewConfirmationTypeFromString(str string) (ConfirmationType, error) {
	valid := map[string]ConfirmationType{
		"PHONE_NUMBER":  ConfirmationTypePhoneNumber,
		"EMAIL_ADDRESS": ConfirmationTypeEmailAddress,
	}
	if t, ok := valid[str]; ok {
		return t, nil
	}
	return ConfirmationType(-1), errors.New("invalid confirmation type")
}

type Confirmation struct {
	ID               int
	Code             int
	ExpiredAt        time.Time
	Type             ConfirmationType
	NotificationUUID *u.UUID
	ConfirmedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
