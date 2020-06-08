package main

import (
	"context"
	"fmt"
	auth "lib/proto"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (ah *AuthHandler) Validate(ctx context.Context, request *auth.Request) (*auth.Response, error) {
	response := &auth.Response{}

	email, ok, err := validateTokenRPC(request.GetToken(), "access")
	if err != nil {
		return nil, err
	} else if !ok {
		response.Authorized = false
		return response, nil
	}

	user, err := ah.findUser(email)
	if err != nil {
		return nil, err
	}

	if user.Access != request.GetToken() {
		log.Println("Got valid access token, but it has already been refreshed.")
		response.Authorized = false
		return response, nil
	}

	response.Authorized = true
	return response, nil
}

// Validates if provided token valid or not and returns token issued email.
func validateTokenRPC(stoken, tokenType string) (string, bool, error) {
	token, err := jwt.Parse(stoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return "", false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", false, nil
	}

	if int64(claims["expTime"].(float64)) < time.Now().Unix() {
		return "", false, nil
	}

	if claims["type"] != tokenType {
		return "", false, nil
	}

	return claims["email"].(string), true, nil
}
