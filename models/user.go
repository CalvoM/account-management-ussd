package models

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/CalvoM/account-management-ussd/db"
	log "github.com/sirupsen/logrus"
)

//User holds content for a single user
type User struct {
	ID         int64
	Name       string
	Email      string
	Password   string
	Registered bool
	Activated  bool
}

//IsUserInDb check if user exists in Db
func IsUserInDb(u *User) (isUserPresent bool) {
	rows := db.DbClient.QueryRow(`select * from reg_users where username=$1 and email=$2;`, u.Name, u.Email)
	var user = User{}
	isUserPresent = false
	switch err := rows.Scan(&user.ID, &user.Name, &user.Email); err {
	case sql.ErrNoRows:
		log.Warn("User not found")
	case nil:
		isUserPresent = true
	default:
		log.Error("Error getting user ", err)
		isUserPresent = true
	}
	return

}

//IsUserEmailInDb check if email in db
func IsUserEmailInDb(u *User) bool {
	rows := db.DbClient.QueryRow(`select * from reg_users where email=$1;`, u.Email)
	var user = User{}
	var isUser bool
	isUser = false
	switch err := rows.Scan(&user.ID, &user.Name, &user.Email); err {
	case sql.ErrNoRows:
		log.Warn("User not found")
	case nil:
		log.Warn("User found")
		isUser = true
	default:
		log.Error("Error getting user ", err)
		isUser = true
	}
	return isUser
}

//AddUser adds user to the databse
func (user *User) AddUser() (id int64, err error) {
	if IsUserInDb(user) {
		err = errors.New("User in DB")
		return
	}
	stmt, err := db.DbClient.Prepare("INSERT INTO reg_users(username,email,password) VALUES($1,$2,$3) RETURNING user_id;")
	if err != nil {
		log.Error(err)
		return
	}
	defer stmt.Close()
	user.HashUserPassword()
	err = stmt.QueryRow(user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		log.Error(err)
		return
	}
	user.ID = id
	log.Info("Added User with ID ", user.ID)
	return
}

//RemoveUser all users with given details
func (user *User) RemoveUser() (err error) {
	stmt, err := db.DbClient.Prepare("delete from reg_users where username=$1 and email=$2;")
	if err != nil {
		log.Error(err)
		return err
	}
	rows, err := stmt.Query(user.Name, user.Email)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(rows.Columns())
	return nil
}

//GetAllUsers Get all users in the system
func GetAllUsers() (users []User, err error) {
	rows, err := db.DbClient.Query("Select * from reg_users;")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		u := User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			log.Error(err)
			return users, err
		}
		users = append(users, u)
	}
	return
}

func (user *User) GetUserByID(id int64) (_user User, err error) {

	return
}

//HashUserPassword hash the password for the user
func (user *User) HashUserPassword() {
	if b, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		log.Error(err)
	} else {
		user.Password = string(b)
	}

}

//CheckUserPassword check the user hashed password against the one supplied
func (user *User) CheckUserPassword(suppliedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(suppliedPassword))
	return err == nil
}
