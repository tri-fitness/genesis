package domain

import (
	"errors"
	"time"
)

type SubscriptionType int

const (
	SubscriptionTypeTrial = iota
	SubscriptionTypePremium
)

func (t SubscriptionType) String() string {
	return []string{
		"TRIAL",
		"PREMIUM",
	}[t]
}

func NewSubscriptionTypeFromString(str string) (SubscriptionType, error) {
	valid := map[string]SubscriptionType{
		"TRIAL":   SubscriptionTypeTrial,
		"PREMIUM": SubscriptionTypePremium,
	}
	if t, ok := valid[str]; ok {
		return t, nil
	}
	return SubscriptionType(-1), errors.New("invalid subscription type")
}

type SubscriptionStatus int

const (
	SubscriptionStatusPending = iota
	SubscriptionStatusActive
	SubscriptionStatusExpired
	SubscriptionStatusCancelled
)

func (t SubscriptionStatus) String() string {
	return []string{
		"TRIAL",
		"PREMIUM",
	}[t]
}

func NewSubscriptionStatusFromString(str string) (SubscriptionStatus, error) {
	valid := map[string]SubscriptionStatus{
		"PENDING":   SubscriptionStatusPending,
		"ACTIVE":    SubscriptionStatusActive,
		"EXPIRED":   SubscriptionStatusExpired,
		"CANCELLED": SubscriptionStatusCancelled,
	}
	if t, ok := valid[str]; ok {
		return t, nil
	}
	return SubscriptionStatus(-1), errors.New("invalid subscription status")
}

type SubscriptionInterval int

const (
	SubscriptionIntervalMonthly = iota
	SubscriptionIntervalAnnual
)

func (t SubscriptionInterval) String() string {
	return []string{
		"MONTHLY",
		"ANNUAL",
	}[t]
}

func NewSubscriptionInterval(str string) (SubscriptionInterval, error) {
	valid := map[string]SubscriptionInterval{
		"MONTHLY": SubscriptionIntervalMonthly,
		"ANNUAL":  SubscriptionIntervalAnnual,
	}
	if t, ok := valid[str]; ok {
		return t, nil
	}
	return SubscriptionInterval(-1), errors.New("invalid subscription interval")
}

type Subscription struct {
	ID          int
	StartedAt   time.Time
	EndedAt     time.Time
	CancelledAt *time.Time
	Type        SubscriptionType
	Status      SubscriptionStatus
	Interval    SubscriptionInterval
}
