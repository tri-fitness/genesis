package domain

import (
	"time"

	u "github.com/gofrs/uuid"
)

type Workout struct {
	UUID        u.UUID
	Exercises   []Exercise
	Description string
	Name        string
	Duration    time.Duration
	Difficulty  Difficulty
	Completions int
	Author      Author
}
