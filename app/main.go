package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query()["access"]) == 0 {
			badRequest(w, "access param not found2")
			return
		}
		access := r.URL.Query()["access"][0]

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

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&Product{})

	onlineShopHandler := OnlineShopHandler{Db: db}

	mux := http.DefaultServeMux

	http.HandleFunc("/product-list", onlineShopHandler.handlerProductList)
	http.HandleFunc("/product-info", onlineShopHandler.handlerProductInfo)
	http.HandleFunc("/new-product", onlineShopHandler.handlerNewProduct)
	http.HandleFunc("/change-product-by-name", onlineShopHandler.handlerChangeProductByName)
	http.HandleFunc("/delete-product", onlineShopHandler.handlerDeleteProduct)

	log.Printf("App started on port: %d\n", *port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), middleware(mux)))
}
