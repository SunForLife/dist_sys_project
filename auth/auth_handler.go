package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AuthHandler struct {
	Users []User
}

func (ah *AuthHandler) findUser(email string) (*User, error) {
	for _, user := range ah.Users {
		if email == user.Email {
			return &user, nil
		}
	}
	return nil, errors.New("User doesn't exist.")
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

	ah.Users = append(ah.Users, user)

	returnTokens(w, user.Email)
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

	returnTokens(w, user.Email)
}

func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["refresh"]) == 0 {
		badRequest(w, "refresh param not found.")
		return
	}
	refresh := r.URL.Query()["refresh"][0]

	email, ok := validateToken(w, refresh, "refresh")
	if !ok {
		return
	}

	returnTokens(w, email)
}

func (ah *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	access := r.Header.Get("access")
	if access == "" {
		if len(r.URL.Query()["access"]) == 0 {
			badRequest(w, "access param not found")
			return
		}
		access = r.URL.Query()["access"][0]
	}

	_, ok := validateToken(w, access, "access")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
