package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/jinzhu/gorm"
)

// var authWaiters = make(map[string]*User)
// var authTimeouts = make(map[string]time.Time)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type AuthHandler struct {
	Db *gorm.DB
	// Ch *amqp.Channel
	// Mq *amqp.Queue
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

		// authCode := randSeq(6)
		// sms := "localhost:9191/approve?code=" + authCode
		// err = ah.Ch.Publish(
		// 	"",         // exchange
		// 	ah.Mq.Name, // routing key
		// 	false,      // mandatory
		// 	false,      // immediate
		// 	amqp.Publishing{
		// 		ContentType: "text/plain",
		// 		Body:        []byte(sms),
		// 	})
		// if err != nil {
		// 	badRequest(w, fmt.Sprint("error sending sms:", err))
		// 	return
		// }
		// authWaiters[authCode] = &user
		// authTimeouts[authCode] = time.Now().Add(5 * time.Minute)
		ah.Db.Create(&user)
	}
}

// func (ah *AuthHandler) Approve(w http.ResponseWriter, r *http.Request) {
// 	if len(r.URL.Query()["code"]) == 0 {
// 		badRequest(w, "code param not found")
// 		return
// 	}
// 	authCode := r.URL.Query()["code"][0]

// 	if _, ok := authTimeouts[authCode]; !ok {
// 		badRequest(w, "incorrect authCode")
// 		return
// 	}

// 	if authTimeouts[authCode].After(time.Now()) {
// 		ah.Db.Create(authWaiters[authCode])
// 	} else {
// 		badRequest(w, "timeout exceeded")
// 		return
// 	}
// 	delete(authWaiters, authCode)
// 	delete(authTimeouts, authCode)
// }

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
