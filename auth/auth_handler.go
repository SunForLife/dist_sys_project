package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

type AuthHandler struct {
	Db *gorm.DB
}

func (ah *AuthHandler) findUser(email string) (*User, error) {
	var user User
	if err := ah.Db.Where("Email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New(fmt.Sprint("Not found user with email:", email))
	}
	return &user, nil
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ACCESS_TIME_DURATION_MINUTES:", ACCESS_TIME_DURATION_MINUTES)
	fmt.Println("REFRESH_TIME_DURATION_MINUTES:", REFRESH_TIME_DURATION_MINUTES)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRequest(w, fmt.Sprint("Error reading request body: ", err))
		return
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		badRequest(w, fmt.Sprint("Error unmarshal request body: ", err))
		return
	}

	_, err = ah.findUser(user.Email)
	if err == nil {
		badRequest(w, fmt.Sprint("Email already exists."))
		return
	}

	user.PassHash = md5Hash(user.Password)
	user.Password = ""

	tokens := returnTokens(w, user.Email)
	if tokens != nil {
		user.Access = tokens.Access
		user.Refresh = tokens.Refresh
		ah.Db.Create(&user)
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRequest(w, fmt.Sprint("Error reading request body: ", err))
		return
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		badRequest(w, fmt.Sprint("Error unmarshal request body: ", err))
		return
	}

	dbUser, err := ah.findUser(user.Email)
	if err != nil {
		badRequest(w, fmt.Sprintf("%s", err))
		return
	}

	if dbUser.PassHash != md5Hash(user.Password) {
		badRequest(w, "Incorrect password.")
		return
	}

	tokens := returnTokens(w, user.Email)
	if tokens != nil {
		ah.Db.Delete(&dbUser, "Email = ?", dbUser.Email)
		dbUser.Access = tokens.Access
		dbUser.Refresh = tokens.Refresh
		ah.Db.Create(&dbUser)
	}
}

func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refresh := r.Header.Get("refresh")
	if refresh == "" {
		badRequest(w, "refresh param not found.")
		return
	}

	email, ok := validateToken(w, refresh, "refresh")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := ah.findUser(email)
	if err != nil {
		badRequest(w, "No user found in ah.Refresh.")
		return
	}
	if user.Refresh != refresh {
		log.Println("Got valid refresh token, but it has already been refreshed.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokens := returnTokens(w, email)
	if tokens != nil {
		ah.Db.Delete(user, "Email = ?", email)
		user.Access = tokens.Access
		user.Refresh = tokens.Refresh
		ah.Db.Create(user)
	}
}

func (ah *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	access := r.Header.Get("access")
	if access == "" {
		badRequest(w, "access param not found")
		return
	}

	email, ok := validateToken(w, access, "access")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := ah.findUser(email)
	if err != nil {
		badRequest(w, "No user found in ah.Validate.")
		return
	}

	if user.Access != access {
		log.Println("Got valid access token, but it has already been refreshed.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
