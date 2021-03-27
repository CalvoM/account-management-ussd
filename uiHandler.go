package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CalvoM/account-management-ussd/models"
	log "github.com/sirupsen/logrus"
)

type userUpdate struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	ConfirmPasword string `json:"confirm"`
}

func serveUi(w http.ResponseWriter, r *http.Request) {
	basePath := filepath.Join("ui", "templates", "base.html")
	urlPath := r.URL.Path
	token := r.URL.Query().Get("token")
	uid := r.URL.Query().Get("uid")
	if token != "" && uid != "" {
		_uid, err := ValidateToken(token)
		if err != nil {
			http.Error(w, "Auth Error", http.StatusBadRequest)
			return
		}
		i_uid, err := strconv.Atoi(uid)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		if int64(i_uid) != _uid {
			http.Error(w, "Auth Error", http.StatusBadRequest)
			return
		}
		parts := strings.Split(urlPath, "/?")
		urlPath = parts[0] + "login.html"
		reqPath := filepath.Join("ui", "templates", filepath.Clean(urlPath))
		info, err := os.Stat(reqPath)
		if err != nil {
			log.Error("Error=>", err)
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
		}
		if info.IsDir() {
			http.NotFound(w, r)
			return
		}
		tmpl, err := template.ParseFiles(basePath, reqPath)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = tmpl.ExecuteTemplate(w, "base", token)
		if err != nil {
			log.Fatal(err.Error())
		}

	} else {
		fmt.Fprint(w, "Welcome to the this APP")
	}
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var u userUpdate
	token := r.Header.Get("Authorization")
	acc_token := strings.Split(token, " ")[1]
	uid, err := ValidateToken(acc_token)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	user := models.User{
		Email: u.Email,
		ID:    uid,
	}
	err = user.UpdateUserPassword(u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	fmt.Fprint(w, "Okay")
	return
}
