package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", 7171, "port")
	flag.Parse()

	onlineShopHandler := OnlineShopHandler{Products: make([]Product, 0)}

	log.Printf("Started on port: %d\n", *port)
	http.HandleFunc("/get-product-list", onlineShopHandler.handlerGetProductList)
	http.HandleFunc("/get-product-info", onlineShopHandler.handlerGetProductInfo)
	http.HandleFunc("/create-new-product", onlineShopHandler.handlerCreateNewProduct)
	http.HandleFunc("/change-product-by-name", onlineShopHandler.handlerChangeProductByName)
	http.HandleFunc("/delete-product", onlineShopHandler.handlerDeleteProduct)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil))
}
