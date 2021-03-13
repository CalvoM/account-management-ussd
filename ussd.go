package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type USSDState int

const (
	Begin USSDState = iota
	GetOption
	RegGetEmail
	RegGetusername
	RegDoneOK
	VGGetPassword
	VGGetContent
	VGDoneOK
	VSGetPassword
	VSSetContent
	VSDoneOK
	Error
)

//SessionDetails store session details sent to callback
type SessionDetails struct {
	SessionID   string
	PhoneNumber string
	NetworkCode string
	ServiceCode string
	Text        string
}

//EndSessionDetails store end of session details
type EndSessionDetails struct {
	SessionID    string
	ServiceCode  string
	NetworkCode  string
	PhoneNumber  string
	Status       string
	Input        string
	ErrorMessage string
}

//generateUSSDResponse respond to USSD
func generateUSSDResponse(text string, session SessionDetails) (resp string) {
	data, err := redisClient.HGet(session.SessionID, "state").Result()
	if err != nil { //Session Handling err
		return "END Error Detected."
	}
	i, err := strconv.Atoi(data)
	if err != nil { //Session Handling err
		return "END Error Detected."
	}
	ussdState := USSDState(i)
	if text == "" && ussdState != Begin { //If user does not supply input and it is not start
		resp += "END Error Detected.\nPlease provide an input."
		ussdState = Begin
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return
		}
		return
	}

	switch {
	case ussdState == Begin:
		resp = "CON Welcome to the Bazenga Vault.\r\n"
		resp += "What would you like to do?\r\n"
		resp += "1. Register as a new user.\r\n"
		resp += "2. Add to the vault.\r\n"
		resp += "3. Get all items in the vault."
		ussdState = GetOption
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}

	case ussdState == GetOption:
		if text == "1" {
			ussdState = RegGetEmail
			if err := updateRedisSession(ussdState, session.SessionID); err != nil {
				return "END Error detected."
			}
			resp = handleRegistration(session, ussdState)
		} else if text == "2" {
			ussdState = VSGetPassword
			if err := updateRedisSession(ussdState, session.SessionID); err != nil {
				return "END Error detected."
			}
			resp = handleUpdateVault(session, ussdState)
		} else if text == "3" {
			ussdState = VGGetPassword
			if err := updateRedisSession(ussdState, session.SessionID); err != nil {
				return "END Error detected."
			}
			resp = handleGetVaultItems(session, ussdState)
		} else {
			resp = "END Not chosen"
		}
	case ussdState >= RegGetEmail && ussdState <= RegDoneOK:
		resp = handleRegistration(session, ussdState)

	case ussdState >= VSGetPassword && ussdState <= VSDoneOK:
		resp = handleUpdateVault(session, ussdState)
	case ussdState >= VGGetPassword && ussdState <= VGDoneOK:
		resp = handleGetVaultItems(session, ussdState)
	}

	return
}

func handleRegistration(session SessionDetails, state USSDState) (resp string) {
	switch state {
	case RegGetEmail:
		resp = "CON Please enter your email"
		ussdState := RegGetusername
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}
	case RegGetusername:
		resp = "CON Please enter your user name"
		ussdState := RegDoneOK
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}
	case RegDoneOK:
		resp = "END You will receive an SMS to complete registration."
	}
	return
}
func handleUpdateVault(session SessionDetails, state USSDState) (resp string) {
	switch state {
	case VSGetPassword:
		resp = "CON Enter password"
		ussdState := VSSetContent
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}
	case VSSetContent:
		resp = "CON Enter item to save to the vault"
		ussdState := VSDoneOK
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}
	case VSDoneOK:
		resp = "END Your vault has been updated"
	}
	return
}

func handleGetVaultItems(session SessionDetails, state USSDState) (resp string) {
	switch state {
	case VGGetPassword:
		resp = "CON Enter password"
		ussdState := VGGetContent
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}
	case VGGetContent:
		resp = "END Sending items via SMS"
		ussdState := VGDoneOK
		if err := updateRedisSession(ussdState, session.SessionID); err != nil {
			return "END Error detected."
		}
	}
	return
}
func getUserChoice(text string) string {
	vals := strings.Split(text, "*")
	return vals[len(vals)-1]
}

func updateRedisSession(state USSDState, sessionID string) error {
	stateStr := strconv.Itoa(int(state))
	err := redisClient.HSet(sessionID, "state", stateStr).Err()
	return err
}

//USSDHandler handle details of ussd sessions
func USSDHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Compare(contentType, "application/x-www-form-urlencoded") == 0 {
		r.ParseForm()
		sessionDet := SessionDetails{}
		sessionDet.SessionID = r.PostForm.Get("sessionId")
		sessionDet.PhoneNumber = r.PostForm.Get("phoneNumber")
		sessionDet.NetworkCode = r.PostForm.Get("networkCode")
		sessionDet.ServiceCode = r.PostForm.Get("serviceCode")
		sessionDet.Text = r.PostForm.Get("text")
		val, err := redisClient.Exists(sessionDet.SessionID).Result()
		if err != nil || val == 0 {
			log.Warn(sessionDet.SessionID, " Not found")
			state := Begin
			err = updateRedisSession(state, sessionDet.SessionID)
			if err != nil {
				log.Error(err)
			}
		}
		fmt.Fprintf(w, "%s", generateUSSDResponse(getUserChoice(sessionDet.Text), sessionDet))
	}

}

//USSDEndNotificationHandler gets details of just ended session
func USSDEndNotificationHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Compare(contentType, "application/x-www-form-urlencoded") == 0 {
		r.ParseForm()
		sessionDet := EndSessionDetails{}
		sessionDet.SessionID = r.PostForm.Get("sessionId")
		sessionDet.ServiceCode = r.PostForm.Get("serviceCode")
		sessionDet.NetworkCode = r.PostForm.Get("networkCode")
		sessionDet.PhoneNumber = r.PostForm.Get("phoneNumber")
		sessionDet.Input = r.PostForm.Get("input")
		sessionDet.Status = r.PostForm.Get("status")
		if strings.Compare(sessionDet.Status, "Failed") == 0 {
			sessionDet.ErrorMessage = r.PostForm.Get("errorMessage")
		}
		log.Info(sessionDet)
	}
}
