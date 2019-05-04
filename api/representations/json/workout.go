package json

import (
	"encoding/json"
	"time"
	r "tri-fitness/genesis/api/representations"

	u "github.com/gofrs/uuid"
)

type Workout struct {
	r.Representation

	UUID        u.UUID
	Exercises   []Exercise
	Description string
	Name        string
	Duration    time.Duration
	Difficulty  string
	Completions int
	Author      u.UUID
}

// AsBytes provides the representation as bytes.
func (w Workout) AsBytes() []byte {
	bytes, err := json.Marshal(&w)
	if err != nil {
		panic(err)
	}
	return bytes
}
