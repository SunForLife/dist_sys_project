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
	port := flag.Int("port", 8181, "port")
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
	db.AutoMigrate(&Product{})

	uploadHandler := UploadHandler{Db: db}

	http.HandleFunc("/upload", uploadHandler.Upload)

	log.Printf("Upload started on port: %d\n", *port)
	log.Fatal(http.ListenAndServe("upload:"+strconv.Itoa(*port), nil)) // "upload:"+strconv.Itoa(*port)
}
