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

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access := r.Header.Get("access")
		if access == "" {
			badRequest(w, "access param not found in middleware")
			return
		}

		if err := validateAccess(access); err != nil {
			errorRequest(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	port := flag.Int("port", 7171, "port")
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

	onlineShopHandler := OnlineShopHandler{Db: db}

	mux := http.DefaultServeMux

	http.HandleFunc("/products", onlineShopHandler.handlerProducts)
	http.HandleFunc("/product-info", onlineShopHandler.handlerProductInfo)
	http.HandleFunc("/new-product", onlineShopHandler.handlerNewProduct)
	http.HandleFunc("/change-product-by-name", onlineShopHandler.handlerChangeProductByName)
	http.HandleFunc("/delete-product", onlineShopHandler.handlerDeleteProduct)

	log.Printf("App started on port: %d\n", *port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), middleware(mux)))
}
