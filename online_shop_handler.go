package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

// GetProductList handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerGetProductList(w http.ResponseWriter, r *http.Request) {
	log.Println("Got get-product-list request")

	products := []Product{}
	osh.Db.Find(&products)

	json, err := json.Marshal(products)
	if err != nil {
		ErrorRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

// GetProductInfo handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerGetProductInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Got get-product-info request")

	if len(r.URL.Query()["name"]) == 0 {
		BadRequest(w, "name param not found")
		return
	}
	name := r.URL.Query()["name"][0]

	var product Product
	if err := osh.Db.Where("Name = ?", name).First(&product).Error; err != nil {
		NotFoundRequest(w, fmt.Sprint("Not found product with name:", name))
		return
	}
	json, err := json.Marshal(product)
	if err != nil {
		ErrorRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(json))
}

// CreateNewProduct handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerCreateNewProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Got create-new-product request")

	if len(r.URL.Query()["name"]) == 0 {
		BadRequest(w, "name param not found")
		return
	}
	name := r.URL.Query()["name"][0]
	if len(r.URL.Query()["code"]) == 0 {
		BadRequest(w, "code param not found")
		return
	}
	code := r.URL.Query()["code"][0]
	if len(r.URL.Query()["category"]) == 0 {
		BadRequest(w, "category param not found")
		return
	}
	category := r.URL.Query()["category"][0]

	osh.Db.Create(&Product{Name: name, Code: code, Category: category})
}

// ChangeProductByName handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerChangeProductByName(w http.ResponseWriter, r *http.Request) {
	log.Println("Got change-product-by-name request")

	if len(r.URL.Query()["old-name"]) == 0 {
		BadRequest(w, "old-name param not found")
		return
	}
	oldName := r.URL.Query()["old-name"][0]
	if len(r.URL.Query()["name"]) == 0 {
		BadRequest(w, "name param not found")
		return
	}
	name := r.URL.Query()["name"][0]
	if len(r.URL.Query()["code"]) == 0 {
		BadRequest(w, "code param not found")
		return
	}
	code := r.URL.Query()["code"][0]
	if len(r.URL.Query()["category"]) == 0 {
		BadRequest(w, "category param not found")
		return
	}
	category := r.URL.Query()["category"][0]

	var product Product
	if err := osh.Db.Where("Name = ?", oldName).First(&product).Error; err != nil {
		NotFoundRequest(w, fmt.Sprint("Not found product with name:", oldName))
		return
	}

	osh.Db.Delete(&product, "Name = ?", oldName)
	osh.Db.Create(&Product{Name: name, Code: code, Category: category})
}

// DeleteProduct handler method of OnlineShopHandler.
func (osh *OnlineShopHandler) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Got delete-product request")

	if len(r.URL.Query()["name"]) == 0 {
		BadRequest(w, "name param not found")
		return
	}
	name := r.URL.Query()["name"][0]

	var product Product
	if err := osh.Db.Where("Name = ?", name).First(&product).Error; err != nil {
		NotFoundRequest(w, fmt.Sprint("Not found product with name:", name))
		return
	}

	osh.Db.Delete(&product, "Name = ?", name)
}

func BadRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	if err != "" {
		log.Println("error: ", err)
	} else {
		log.Println("Bad request")
	}
}

func NotFoundRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusNotFound)
	if err != "" {
		log.Println("error: ", err)
	} else {
		log.Println("Bad request")
	}
}

func ErrorRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		log.Println("error: ", err)
	} else {
		log.Println("Error while handling request")
	}
}
