package json

import (
	"encoding/json"
	"net/http"
	"time"
	r "tri-fitness/genesis/api/representations"
	"tri-fitness/genesis/domain"

	u "github.com/gofrs/uuid"
)

type jsonRepresentationFactory struct{}

func NewJSONRepresentationFactory() r.RepresentationFactory {
	return jsonRepresentationFactory{}
}

func (j jsonRepresentationFactory) decode(
	request *http.Request, representation interface{}) error {
	if err := json.NewDecoder(request.Body).Decode(&representation); err != nil {
		return err
	}
	return nil
}

func (j jsonRepresentationFactory) Account(
	account domain.Account) (r.Representation, error) {

	representation := Account{
		UUID:                account.UUID,
		PrimaryCredential:   account.PrimaryCredential,
		SecondaryCredential: account.SecondaryCredential,
		Type:                account.Type.String(),
		GivenName:           account.GivenName,
		Surname:             account.Surname,
		Bio:                 account.Bio,
		Email:               account.Email,
		Phone:               account.Phone,
	}

	confirmationRepresentations := []Confirmation{}
	for _, e := range account.Confirmations {
		confirmationRepresentations =
			append(confirmationRepresentations, j.confirmation(e))
	}
	representation.Confirmations = confirmationRepresentations

	return &representation, nil
}

func (j jsonRepresentationFactory) AccountFromRequest(
	request *http.Request) (r.Representation, error) {

	representation := Account{}
	return &representation, j.decode(request, &representation)
}

func (j jsonRepresentationFactory) AccountEntityFromRequest(
	request *http.Request) (domain.Account, error) {

	representation := Account{}
	if err := j.decode(request, &representation); err != nil {
		return domain.Account{}, err
	}

	t, err := domain.AccountTypeFromString(representation.Type)
	if err != nil {
		return domain.Account{}, err
	}

	var uuid u.UUID
	if representation.UUID == uuid {
		if uuid, err = u.NewV4(); err != nil {
			return domain.Account{}, err
		}
		representation.UUID = uuid
	}

	confirmations := []domain.Confirmation{}
	for _, c := range representation.Confirmations {
		confirmation, err := j.confirmationEntity(c)
		if err != nil {
			return domain.Account{}, err
		}
		confirmations = append(confirmations, confirmation)
	}

	account := domain.Account{
		UUID:                uuid,
		Type:                t,
		PrimaryCredential:   representation.PrimaryCredential,
		SecondaryCredential: representation.SecondaryCredential,
		GivenName:           representation.GivenName,
		Surname:             representation.Surname,
		Bio:                 representation.Bio,
		Email:               representation.Email,
		Phone:               representation.Phone,
		Confirmations:       confirmations,
	}
	now := time.Now()
	account.CreatedAt, account.UpdatedAt = now, now
	return account, nil
}

func (j jsonRepresentationFactory) Workout(
	workout domain.Workout) (r.Representation, error) {

	representation := Workout{
		UUID:        workout.UUID,
		Description: workout.Description,
		Name:        workout.Name,
		Duration:    workout.Duration,
		Difficulty:  workout.Difficulty.String(),
		Completions: workout.Completions,
		Author:      workout.Author.UUID,
	}

	exerciseRepresentations := []Exercise{}
	for _, e := range workout.Exercises {
		exerciseRepresentations =
			append(exerciseRepresentations, j.exercise(e))
	}
	representation.Exercises = exerciseRepresentations
	return &representation, nil
}

func (j jsonRepresentationFactory) WorkoutFromRequest(
	request *http.Request) (r.Representation, error) {

	representation := Workout{}
	return representation, j.decode(request, representation)
}

func (j jsonRepresentationFactory) WorkoutEntityFromRequest(
	request *http.Request) (domain.Workout, error) {

	representation := Workout{}
	if err := j.decode(request, &representation); err != nil {
		return domain.Workout{}, err
	}
	d, err := domain.NewDifficultyFromString(representation.Difficulty)
	if err != nil {
		return domain.Workout{}, err
	}
	exercises := []domain.Exercise{}
	for _, e := range representation.Exercises {
		exercise, err := j.exerciseEntity(e)
		if err != nil {
			return domain.Workout{}, err
		}
		exercises = append(exercises, exercise)
	}
	workout := domain.Workout{
		UUID:        representation.UUID,
		Description: representation.Description,
		Name:        representation.Name,
		Completions: representation.Completions,
		Duration:    representation.Duration,
		Difficulty:  d,
		Exercises:   exercises,
		Author: domain.Author{
			UUID: representation.Author,
		},
	}
	return workout, nil
}

func (j jsonRepresentationFactory) Exercise(
	exercise domain.Exercise) (r.Representation, error) {

	representation := j.exercise(exercise)
	return &representation, nil
}

func (j jsonRepresentationFactory) ExerciseFromRequest(
	request *http.Request) (r.Representation, error) {

	representation := Exercise{}
	return representation, j.decode(request, &representation)
}

func (j jsonRepresentationFactory) ExerciseEntityFromRequest(
	request *http.Request) (domain.Exercise, error) {

	representation := Exercise{}
	if err := j.decode(request, &representation); err != nil {
		return domain.Exercise{}, err
	}

	return j.exerciseEntity(representation)
}

func (j jsonRepresentationFactory) exercise(
	exercise domain.Exercise) Exercise {

	representation := Exercise{
		UUID:       exercise.UUID,
		Difficulty: exercise.Difficulty.String(),
		Watches:    exercise.Watches,
		Duration:   exercise.Duration,
		Author:     exercise.Author.UUID,
	}

	return representation
}

func (j jsonRepresentationFactory) exerciseEntity(
	representation Exercise) (domain.Exercise, error) {

	d, err := domain.NewDifficultyFromString(representation.Difficulty)
	if err != nil {
		return domain.Exercise{}, err
	}

	exercise := domain.Exercise{
		UUID:       representation.UUID,
		Difficulty: d,
		Watches:    representation.Watches,
		Duration:   representation.Duration,
		Author: domain.Author{
			UUID: representation.Author,
		},
	}
	return exercise, nil
}

func (j jsonRepresentationFactory) confirmation(
	confirmation domain.Confirmation) Confirmation {

	representation := Confirmation{
		ID:          confirmation.ID,
		ExpiredAt:   confirmation.ExpiredAt,
		Type:        confirmation.Type.String(),
		ConfirmedAt: confirmation.ConfirmedAt,
		CreatedAt:   confirmation.CreatedAt,
	}
	return representation
}

func (j jsonRepresentationFactory) confirmationEntity(
	representation Confirmation) (domain.Confirmation, error) {

	confirmationType, err :=
		domain.NewConfirmationTypeFromString(representation.Type)
	if err != nil {
		return domain.Confirmation{}, err
	}
	return domain.Confirmation{
		ID:          representation.ID,
		ConfirmedAt: representation.ConfirmedAt,
		CreatedAt:   representation.CreatedAt,
		ExpiredAt:   representation.ExpiredAt,
		Type:        confirmationType,
	}, nil
}

func (j jsonRepresentationFactory) CodeEntityFromRequest(
	request *http.Request) (domain.Code, error) {

	representation := Code{}
	if err := j.decode(request, &representation); err != nil {
		return domain.Code{}, err
	}

	code := domain.Code{
		Code: representation.Code,
	}
	return code, nil
}
