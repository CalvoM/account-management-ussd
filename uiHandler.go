package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
			log.Info(_uid, int64(i_uid))
			http.Error(w, "Auth Error", http.StatusBadRequest)
			return
		}
		parts := strings.Split(urlPath, "/?")
		urlPath = parts[0] + "login.html"
	} else {
		urlPath += "index.html"
	}
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
	err = tmpl.ExecuteTemplate(w, "base", "Test")
	if err != nil {
		log.Fatal(err.Error())
	}

}
