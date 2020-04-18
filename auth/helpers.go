package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")                                                        // "VERYSTRONGSECRET"
var ACCESS_TIME_DURATION_MINUTES, _ = strconv.Atoi(os.Getenv("ACCESS_TIME_DURATION_MINUTES"))   // 5
var REFRESH_TIME_DURATION_MINUTES, _ = strconv.Atoi(os.Getenv("REFRESH_TIME_DURATION_MINUTES")) // 20

// Describes User.
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PassHash string `json:"passhash"`
}

// Describes Tokens.
type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

// Returns md5 hash from string.
func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Issues access and refresh JWT tokens from email.
func issueTokens(email string) (string, string, error) {
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":    "access",
		"email":   email,
		"expTime": time.Now().Add(time.Minute * time.Duration(ACCESS_TIME_DURATION_MINUTES)).Unix(),
	})
	access, err := accessJwt.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":    "refresh",
		"email":   email,
		"expTime": time.Now().Add(time.Minute * time.Duration(REFRESH_TIME_DURATION_MINUTES)).Unix(),
	})
	refresh, err := refreshJwt.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// Issues tokens from provided email
// and returns them into provided ResponseWriter.
func returnTokens(w http.ResponseWriter, email string) {
	var tokens Tokens
	var err error
	tokens.Access, tokens.Refresh, err = issueTokens(email)
	if err != nil {
		errorRequest(w, err)
		return
	}

	json, err := json.Marshal(tokens)
	if err != nil {
		errorRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

// Validates is provided token valid or not and returns token issued email.
func validateToken(w http.ResponseWriter, stoken, tokenType string) (string, bool) {
	token, err := jwt.Parse(stoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		badRequest(w, fmt.Sprint("Error parse token: ", err))
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		badRequest(w, "Invalid token.")
		return "", false
	}

	if int64(claims["expTime"].(float64)) < time.Now().Unix() {
		badRequest(w, "Expired token.")
		return "", false
	}

	if claims["type"] != tokenType {
		badRequest(w, "Unexpected token.")
		return "", false
	}

	return claims["email"].(string), true
}

func badRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	writeErrorBody(w, err)
}

func notFoundRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusNotFound)
	writeErrorBody(w, err)
}

func errorRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	writeErrorBody(w, err)
}

func writeErrorBody(w http.ResponseWriter, err interface{}) {
	var sErr string
	if err != nil {
		sErr = fmt.Sprint("error: ", err)
	} else {
		sErr = fmt.Sprint("Error while handling request.")
	}
	log.Println(sErr)
	w.Write([]byte(sErr))
}
