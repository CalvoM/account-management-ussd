package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	DbClient *sql.DB
	Ctx      context.Context
)

func init() {
	Ctx = context.TODO()
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	host := viper.Get("DB_HOST").(string)
	port := viper.GetInt("PG_PORT")
	user := viper.Get("POSTGRES_USER").(string)
	password := viper.Get("POSTGRES_PASSWORD").(string)
	dbname := viper.Get("POSTGRES_DB")
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DbClient, err = sql.Open("postgres", pgInfo)
	if err != nil {
		panic(err)
	}
	err = DbClient.Ping()
	if err != nil {
		panic(err)
	}
	log.Info("DB connection OK")
}
