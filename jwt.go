package main

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(user_id int64) (token string, err error) {
	viper.SetConfigFile("./.env")
	err = viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		return
	}

	claims := UserClaims{
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "ussd",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tok.SignedString([]byte(viper.GetString("SECRET")))
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(token)
	return
}

func ValidateToken(tokenString string) (uid int64, err error) {
	viper.SetConfigFile("./.env")
	err = viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("SECRET")), nil
	})
	claims, ok := token.Claims.(*UserClaims)
	if ok && token.Valid {
		uid = claims.UserID
	}
	return
}
