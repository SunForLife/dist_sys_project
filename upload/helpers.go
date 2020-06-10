package main

import (
	"errors"
	"net/http"
)

func checkUserRole(r *http.Request) error {
	token := r.Header.Get("access")
	if token == "" {
		return errors.New("access param not foundin checkUserRole.")
	}

	req, err := http.NewRequest(http.MethodGet, "http://auth:9191/check-user-role", nil)
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
