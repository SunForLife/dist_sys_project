package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", 9191, "port")
	flag.Parse()

	authHandler := AuthHandler{
		Users: make([]User, 0),
	}

	log.Printf("Auth started on port: %d\n", *port)
	http.HandleFunc("/sign-up", authHandler.SignUp)
	http.HandleFunc("/sign-in", authHandler.SignIn)
	http.HandleFunc("/refresh", authHandler.Refresh)
	http.HandleFunc("/validate", authHandler.Validate)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil))
}
