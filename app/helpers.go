package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func validateAccess(token string) error {
	req, err := http.NewRequest(http.MethodGet, "http://auth:9191/validate", nil)
	if err != nil {
		return err
	}
	req.Header.Set("access", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("unauthorized")
	}
	return nil
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
