package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database.

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil {
		app.errorJSON(w, err, http.StatusServiceUnavailable)
		return
	}

	if !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.Marshal(entry)

	if err != nil {
		log.Println("Issues formatting logging json", err)
		return err
	}

	logService := "http://logger-service/log"

	request, err := http.NewRequest("POST", logService, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("Sending log to logger failed ", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("Logger failed to process the log", err)
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return errors.New("failed to process the log")
	}

	return nil
}
