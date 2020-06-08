package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

const dbFileName = "test.db"

func startEnv() *AuthHandler {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	db.AutoMigrate(&User{})

	return &AuthHandler{Db: db}
}

func onDefer(ah *AuthHandler) {
	ah.Db.Close()
	os.Remove(dbFileName)
}

func TestSimple(t *testing.T) {
	ah := startEnv()
	defer onDefer(ah)

	rr := httptest.NewRecorder()

	user := User{
		Email:    "dany@gmail.com",
		Password: "123abc",
	}

	jsonUser, err := json.Marshal(&user)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "/sign-up", bytes.NewBuffer(jsonUser))
	assert.NoError(t, err)

	handler := http.HandlerFunc(ah.SignUp)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tokens Tokens
	err = json.Unmarshal(rr.Body.Bytes(), &tokens)
	assert.NoError(t, err)

	require.NotEmpty(t, tokens.Access)
	require.NotEmpty(t, tokens.Refresh)
}

func TestTokensLogic(t *testing.T) {
	ah := startEnv()
	defer onDefer(ah)

	user := User{
		Email:    "dany@gmail.com",
		Password: "123abc",
	}

	jsonUser, err := json.Marshal(&user)
	assert.NoError(t, err)

	// Making sign-up request.
	req, err := http.NewRequest(http.MethodGet, "/sign-up", bytes.NewBuffer(jsonUser))
	assert.NoError(t, err)

	rrSignUp := httptest.NewRecorder()
	handler := http.HandlerFunc(ah.SignUp)
	handler.ServeHTTP(rrSignUp, req)

	assert.Equal(t, http.StatusOK, rrSignUp.Code)

	var tokensSignUp Tokens
	err = json.Unmarshal(rrSignUp.Body.Bytes(), &tokensSignUp)
	assert.NoError(t, err)

	require.NotEmpty(t, tokensSignUp.Access)
	require.NotEmpty(t, tokensSignUp.Refresh)

	time.Sleep(time.Second)
	// Making sign-in request.
	req, err = http.NewRequest(http.MethodPost, "/sign-in", bytes.NewBuffer(jsonUser))
	assert.NoError(t, err)

	rrSignIn := httptest.NewRecorder()
	handler = http.HandlerFunc(ah.SignIn)
	handler.ServeHTTP(rrSignIn, req)

	assert.Equal(t, http.StatusOK, rrSignIn.Code)

	var tokensSignIn Tokens
	err = json.Unmarshal(rrSignIn.Body.Bytes(), &tokensSignIn)
	assert.NoError(t, err)

	require.NotEmpty(t, tokensSignIn.Access)
	require.NotEmpty(t, tokensSignIn.Refresh)

	// Checking that we can't Refresh tokens with refresh one from SignUp request.
	req, err = http.NewRequest(http.MethodGet, "/refresh", nil)
	assert.NoError(t, err)

	req.Header.Set("refresh", tokensSignUp.Refresh)

	rrRefresh := httptest.NewRecorder()
	handler = http.HandlerFunc(ah.Refresh)
	handler.ServeHTTP(rrRefresh, req)

	assert.Equal(t, http.StatusUnauthorized, rrRefresh.Code)

	// And can with refresh one from SignIn.
	req, err = http.NewRequest(http.MethodGet, "/refresh", nil)
	assert.NoError(t, err)

	req.Header.Set("refresh", tokensSignIn.Refresh)

	rrRefresh2 := httptest.NewRecorder()
	handler = http.HandlerFunc(ah.Refresh)
	handler.ServeHTTP(rrRefresh2, req)

	assert.NotEqual(t, http.StatusUnauthorized, rrRefresh2.Code)

	var tokensRefresh Tokens
	err = json.Unmarshal(rrRefresh2.Body.Bytes(), &tokensRefresh)
	assert.NoError(t, err)

	require.NotEmpty(t, tokensRefresh.Access)
	require.NotEmpty(t, tokensRefresh.Refresh)

	// Checking validation.
	req, err = http.NewRequest(http.MethodGet, "/validate", nil)
	assert.NoError(t, err)

	req.Header.Set("access", tokensRefresh.Access)

	// TODO: switch tests to gRPC Validation.
	// rrValidate := httptest.NewRecorder()
	// handler = http.HandlerFunc(ah.ValidateREST)
	// handler.ServeHTTP(rrValidate, req)

	// assert.Equal(t, http.StatusOK, rrValidate.Code)
}
