package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RequestPayload is the type describing the data that we received
// from the user's browser. We embed a custom type for each of the
// possible payloads (mail, auth, and log).
type RequestPayload struct {
	Action string `json:"action"`
	//Mail   MailPayload `json:"mail,omitempty"`
	Auth AuthPayload `json:"auth,omitempty"`
	//Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.errorJSON(w, err)
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	authServiceURL := fmt.Sprintf("http://%s/authenticate", "authentication-service")

	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode == http.StatusUnauthorized {
		_ = app.errorJSON(w, errors.New("invalid credentials ???"), http.StatusUnauthorized)
		return
	}
	//else if response.StatusCode != http.StatusAccepted {
	//	_ = app.errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
	//	return
	//}

	// create variable we'll read the response.Body from the authentication-service into
	var jsonFromService jsonResponse

	// decode the json we get from the authentication-service into our variable
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// did not authenticate successfully
	if jsonFromService.Error {
		// send error JSON back
		_ = app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// send json back to our end user, with user info embedded
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}
