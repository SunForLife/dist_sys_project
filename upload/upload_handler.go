package main

import (
	"encoding/csv"
	"fmt"
	"io"
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

type UploadHandler struct {
	Db *gorm.DB
}

func (uh *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		badRequest(w, fmt.Sprint("Error r.MultipartReader():", err))
		return
	}
	part, err := reader.NextPart()
	if err != nil {
		badRequest(w, fmt.Sprint("Error reader.NextPart():", err))
		return
	}
	if part.FormName() != "file" {
		badRequest(w, "Got something different from file.")
		return
	}

	csvReader := csv.NewReader(part)
	for i := 0; i != i+1; i++ {
		if i%1000 == 0 {
			fmt.Println("Reading line number:", i)
		}

		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		product := Product{
			Name:     line[0],
			Code:     line[1],
			Category: line[2],
		}

		if err := uh.Db.Where("Name = ?", product.Name).First(&Product{}).Error; !gorm.IsRecordNotFoundError(err) {
			if err != nil {
				badRequest(w, "Product existance check problem.")
				return
			}
			continue
		}

		uh.Db.Create(&product)
	}
	fmt.Println()
}

func badRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
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
