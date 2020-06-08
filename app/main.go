package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	auth "lib/proto"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthClient
}

var authClient *AuthClient

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access := r.Header.Get("access")
		if access == "" {
			badRequest(w, "access param not found in middleware")
			return
		}

		request := &auth.Request{Token: access}
		resp, err := authClient.client.Validate(context.Background(), request)
		if err != nil {
			errorRequest(w, err)
			return
		} else if !resp.Authorized {
			errorRequest(w, errors.New("error unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getRPCClient() (*AuthClient, error) {
	conn, err := grpc.Dial("auth:6161", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	authClient := &AuthClient{
		conn:   conn,
		client: auth.NewAuthClient(conn),
	}

	return authClient, nil
}

func main() {
	port := flag.Int("port", 7171, "port")
	flag.Parse()

	time.Sleep(10 * time.Second)
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

	authClient, err = getRPCClient()
	if err != nil {
		log.Fatal("Failed to connect to RPC server", err)
	}
	defer authClient.conn.Close()

	mux := http.DefaultServeMux

	http.HandleFunc("/products", onlineShopHandler.handlerProducts)
	http.HandleFunc("/product-info", onlineShopHandler.handlerProductInfo)
	http.HandleFunc("/new-product", onlineShopHandler.handlerNewProduct)
	http.HandleFunc("/change-product-by-name", onlineShopHandler.handlerChangeProductByName)
	http.HandleFunc("/delete-product", onlineShopHandler.handlerDeleteProduct)

	log.Printf("App started on port: %d\n", *port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), middleware(mux)))
}
