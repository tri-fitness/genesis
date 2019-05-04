package representations

import (
	"net/http"
	"tri-fitness/genesis/domain"
)

type RepresentationFactory interface {
	Account(domain.Account) (Representation, error)
	AccountFromRequest(*http.Request) (Representation, error)
	AccountEntityFromRequest(*http.Request) (domain.Account, error)
	Workout(domain.Workout) (Representation, error)
	WorkoutFromRequest(*http.Request) (Representation, error)
	WorkoutEntityFromRequest(*http.Request) (domain.Workout, error)
	Exercise(domain.Exercise) (Representation, error)
	ExerciseFromRequest(*http.Request) (Representation, error)
	ExerciseEntityFromRequest(*http.Request) (domain.Exercise, error)
}
