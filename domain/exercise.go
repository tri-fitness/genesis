package domain

import (
	"time"

	u "github.com/gofrs/uuid"
)

type Exercise struct {
	UUID       u.UUID
	Difficulty Difficulty
	Watches    int
	Duration   time.Duration
	Author     Author
}
