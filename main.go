package main

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

var redisClient *redis.Client

func main() {
	log.Info("USSD app server starting")
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Failed to load config file")
		panic(err)
	}
	log.Info("Loaded the config settings")
	redisAddr := viper.Get("DB_HOST").(string) + ":" + viper.Get("REDIS_PORT").(string)
	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Error("Redis DB connection failed")
		panic(err)
	}
	log.Info("Started Redis successfully")
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./ui"))
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", fs))
	r.Path("/ussd/end_note/").HandlerFunc(USSDEndNotificationHandler).Methods("POST")
	r.Path("/ussd/").HandlerFunc(USSDHandler).Methods("POST")
	r.Path("/msg/").HandlerFunc(SMSHandler).Methods("POST")
	r.PathPrefix("/").HandlerFunc(serveUi)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8083",
	}
	log.Fatal(srv.ListenAndServe())
}
