package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type USSDState int
type VaultSetState int
type VaultGetState int

const (
	Begin USSDState = iota
	BeforeEmail
	BeforeName
	DoneOK
	Error
)

const (
	VBegin VaultGetState = iota
	BeforePassword
	BeforeContent
	VDoneOK
	VError
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
	resp = "END Test Response"
	return
}

func getUserChoice(text string) string {
	vals := strings.Split(text, "*")
	return vals[len(vals)-1]
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
		fmt.Fprintf(w, "%s", generateUSSDResponse(getUserChoice(sessionDet.Text), sessionDet))
	}

}

//USSDEndNotificationHandler gets details of just ended session
func USSDEndNotificationHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Compare(contentType, "application/x-www/form-urlencoded") == 0 {
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
