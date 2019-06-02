package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"tri-fitness/genesis/config"
	"tri-fitness/genesis/domain"

	u "github.com/gofrs/uuid"
)

type ShortMessageService struct {
	httpClient *http.Client
	domain     string
}

func NewShortMessageService(
	configuration config.Configuration) ShortMessageService {
	host := configuration.Dependencies.Noti.Host
	port := configuration.Dependencies.Noti.Port
	domain := fmt.Sprintf("http://%s:%d", host, port)
	client := http.DefaultClient
	return ShortMessageService{
		httpClient: client,
		domain:     domain,
	}
}

func (s *ShortMessageService) RegisterAccount(
	account domain.Account) (u.UUID, error) {

	t := target{
		PhoneNumber: account.Phone,
		Name:        fmt.Sprintf("%s %s", account.GivenName, account.Surname),
	}

	mediaType := "application/json"
	j, err := json.Marshal(t)
	if err != nil {
		return u.UUID{}, err
	}

	// create the notification target.
	url := fmt.Sprintf("%s/targets", s.domain)
	request, err :=
		http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(j)))
	if err != nil {
		return u.UUID{}, err
	}
	request.Header.Add("Accept", mediaType)
	request.Header.Add("Content-Type", mediaType)
	pResponse, err := s.httpClient.Do(request)
	if err != nil {
		return u.UUID{}, err
	}
	defer pResponse.Body.Close()
	if pResponse.StatusCode/200 != 1 {
		return u.UUID{}, errors.New("failed request to create notification target")
	}

	// retrieve the newly created target.
	loc, err := pResponse.Location()
	if err != nil {
		return u.UUID{}, err
	}
	request, err = http.NewRequest(http.MethodGet, loc.String(), nil)
	if err != nil {
		return u.UUID{}, err
	}
	request.Header.Add("Accept", mediaType)
	response, err := s.httpClient.Do(request)
	if err != nil {
		return u.UUID{}, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return u.UUID{}, err
	}
	if err = json.Unmarshal(bytes, &t); err != nil {
		return u.UUID{}, err
	}
	return *t.UUID, nil
}

func (s *ShortMessageService) SendPhoneConfirmation(
	account domain.Account, confirmation domain.Confirmation) (u.UUID, error) {

	content :=
		fmt.Sprintf(
			"Welcome to TRI Fitness %s! Here is your phone number confirmation code: %d", account.GivenName, confirmation.Code)
	n := notification{
		Content: content,
		SendAt:  time.Now(),
		Status:  "PENDING",
		Targets: []target{
			{
				UUID:        &account.NotificationTargetUUID,
				PhoneNumber: account.Phone,
				Name:        fmt.Sprintf("%s %s", account.GivenName, account.Surname),
			},
		},
	}

	mediaType := "application/json"
	j, err := json.Marshal(n)
	if err != nil {
		return u.UUID{}, err
	}

	// create the notification target.
	url := fmt.Sprintf("%s/notifications", s.domain)
	request, err :=
		http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(j)))
	if err != nil {
		return u.UUID{}, err
	}
	request.Header.Add("Accept", mediaType)
	request.Header.Add("Content-Type", mediaType)
	pResponse, err := s.httpClient.Do(request)
	if err != nil {
		return u.UUID{}, err
	}
	defer pResponse.Body.Close()
	if pResponse.StatusCode/200 != 1 {
		return u.UUID{}, errors.New("failed request to create notification")
	}

	// retrieve the newly created notification.
	loc, err := pResponse.Location()
	if err != nil {
		return u.UUID{}, err
	}
	request, err = http.NewRequest(http.MethodGet, loc.String(), nil)
	if err != nil {
		return u.UUID{}, err
	}
	request.Header.Add("Accept", mediaType)
	response, err := s.httpClient.Do(request)
	if err != nil {
		return u.UUID{}, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return u.UUID{}, err
	}
	if err = json.Unmarshal(bytes, &n); err != nil {
		return u.UUID{}, err
	}
	return *n.UUID, nil
}

func (s *ShortMessageService) SendEmailConfirmation(
	account domain.Account, confirmation domain.Confirmation) (u.UUID, error) {
	return u.UUID{}, nil
}
