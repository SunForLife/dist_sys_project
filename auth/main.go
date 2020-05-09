package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	port := flag.Int("port", 9191, "port")
	flag.Parse()

	time.Sleep(5 * time.Second)
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"),
	))
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}
	defer db.Close()
	db.AutoMigrate(&User{})

	authHandler := AuthHandler{Db: db}

	log.Printf("Auth started on port: %d\n", *port)
	http.HandleFunc("/sign-up", authHandler.SignUp)
	http.HandleFunc("/sign-in", authHandler.SignIn)
	http.HandleFunc("/refresh", authHandler.Refresh)
	http.HandleFunc("/validate", authHandler.Validate)
	log.Fatal(http.ListenAndServe("auth:"+strconv.Itoa(*port), nil))
}
