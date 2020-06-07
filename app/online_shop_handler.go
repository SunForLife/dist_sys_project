package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Describes Product.
type Product struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Category string `json:"category"`
}

// Handler for HTTP requests, that stores all data.
type OnlineShopHandler struct {
	Db *gorm.DB
}

// Products handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Got products request")

	if len(r.URL.Query()["position"]) == 0 {
		badRequest(w, "position param not found")
		return
	}
	position, err := strconv.Atoi(r.URL.Query()["position"][0])
	if err != nil {
		badRequest(w, "incorrect position param")
		return
	}

	if len(r.URL.Query()["limit"]) == 0 {
		badRequest(w, "limit param not found")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query()["limit"][0])
	if err != nil {
		badRequest(w, "incorrect limit param")
		return
	}

	products := []Product{}
	osh.Db.Limit(position + limit).Find(&products)
	products = products[position:]

	json, err := json.Marshal(products)
	if err != nil {
		errorRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

// ProductInfo handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerProductInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Got get-product-info request")

	if len(r.URL.Query()["name"]) == 0 {
		badRequest(w, "name param not found")
		return
	}
	name := r.URL.Query()["name"][0]

	var product Product
	if err := osh.Db.Where("Name = ?", name).First(&product).Error; err != nil {
		notFoundRequest(w, fmt.Sprint("Not found product with name:", name))
		return
	}
	json, err := json.Marshal(product)
	if err != nil {
		errorRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

// NewProduct handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerNewProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Got create-new-product request")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRequest(w, "error ioutil.ReadAll(r.Body)")
		return
	}

	var product Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		badRequest(w, "error json.Unmarshal(body, &product)")
		return
	}

	osh.Db.Create(&product)
}

// ChangeProductByName handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerChangeProductByName(w http.ResponseWriter, r *http.Request) {
	log.Println("Got change-product-by-name request")

	if len(r.URL.Query()["old-name"]) == 0 {
		badRequest(w, "old-name param not found")
		return
	}
	oldName := r.URL.Query()["old-name"][0]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRequest(w, "error ioutil.ReadAll(r.Body)")
		return
	}

	var product Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		badRequest(w, "error json.Unmarshal(body, &product)")
		return
	}

	if err := osh.Db.Where("Name = ?", oldName).First(&Product{}).Error; err != nil {
		notFoundRequest(w, fmt.Sprint("Not found product with name:", oldName))
		return
	}

	osh.Db.Delete(&Product{}, "Name = ?", oldName)
	osh.Db.Create(&product)
}

// DeleteProduct handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Got delete-product request")

	if len(r.URL.Query()["name"]) == 0 {
		badRequest(w, "name param not found")
		return
	}
	name := r.URL.Query()["name"][0]

	var product Product
	if err := osh.Db.Where("Name = ?", name).First(&product).Error; err != nil {
		notFoundRequest(w, fmt.Sprint("Not found product with name:", name))
		return
	}

	osh.Db.Delete(&product, "Name = ?", name)
}
