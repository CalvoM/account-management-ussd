package main

import (
	"fmt"
	"github.com/AndroidStudyOpenSource/africastalking-go/sms"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func SMSHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.PostForm)
	fmt.Fprintf(w, "Good")
}

//go:generate go run update.go .env AT_USERNAME AT_APIKEY AT_ENV

//NotifyByATSMS send notification using AT's SMS service
func NotifyByATSMS(session SessionDetails, message string) {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	userName := viper.GetString("AT_USERNAME")
	apiKey := viper.GetString("AT_APIKEY")
	env := viper.GetString("AT_ENV")
	smsService := sms.NewService(userName, apiKey, env)
	smsResponse, err := smsService.Send("", session.PhoneNumber, message)
	if err != nil {
		log.Error(err)
	}
	log.Info(smsResponse)
}
