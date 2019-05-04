package domain

import "errors"

type Difficulty int

const (
	DifficultyBeginner = iota
	DifficultyIntermediate
	DifficultyAdvanced
)

func (d Difficulty) String() string {
	return []string{
		"BEGINNER",
		"INTERMEDIATE",
		"ADVANCED",
	}[d]
}

func NewDifficultyFromString(str string) (Difficulty, error) {
	valid := map[string]Difficulty{
		"BEGINNER":     DifficultyBeginner,
		"INTERMEDIATE": DifficultyIntermediate,
		"ADVANCED":     DifficultyAdvanced,
	}
	if d, ok := valid[str]; ok {
		return d, nil
	}
	return Difficulty(-1), errors.New("invalid difficulty")
}
