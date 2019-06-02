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
	UUID                   u.UUID
	PrimaryCredential      string
	SecondaryCredential    string
	Type                   AccountType
	GivenName              string
	Surname                string
	Bio                    string
	Email                  string
	Phone                  string
	Subscriptions          []Subscription
	Confirmations          []Confirmation
	NotificationTargetUUID u.UUID
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              *time.Time
}

var counter int = 0

func (a *Account) NewConfirmation(
	cType ConfirmationType,
	code int,
	expiresAt time.Time,
) Confirmation {
	now := time.Now()
	counter = counter + 1
	c := Confirmation{
		ID:        counter,
		Code:      code,
		ExpiredAt: expiresAt,
		Type:      cType,
		CreatedAt: now,
		UpdatedAt: now,
	}
	a.Confirmations = append(a.Confirmations, c)
	return c
}

func (a *Account) Confirm(code int) error {
	var match Confirmation
	var matched bool
	var matchedIdx int
	now := time.Now()
	for i, c := range a.Confirmations {
		if c.Code == code {
			match = c
			matched = true
			matchedIdx = i
			break
		}
	}

	if !matched {
		return errors.New("no matching confirmation")
	}

	if match.ExpiredAt.Before(now) {
		return errors.New("confirmation has expired")
	}

	match.ConfirmedAt = &now
	a.Confirmations[matchedIdx] = match
	return nil
}

func (a *Account) SetConfirmation(
	confirmation Confirmation) error {

	var matched bool
	var matchedIdx int
	for i, c := range a.Confirmations {
		if c.ID == confirmation.ID {
			matched = true
			matchedIdx = i
			break
		}
	}

	if !matched {
		return errors.New("no matching confirmation")
	}

	a.Confirmations[matchedIdx] = confirmation
	return nil
}
