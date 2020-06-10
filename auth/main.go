package main

import (
	"flag"
	"fmt"
	auth "lib/proto"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
)

func startRPCServer(port string, authHandler *AuthHandler) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthServer(s, authHandler)

	log.Printf("RPC serving on %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("RPC failed to serve: %v", err)
	}
}

func main() {
	port := flag.Int("port", 9191, "port")
	flag.Parse()

	time.Sleep(5 * time.Second)
	// Postgres part.
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

	// Rabbitmq part.
	// conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	// if err != nil {
	// 	log.Fatal("Failed to connect rabbit", err)
	// }
	// defer conn.Close()

	// ch, err := conn.Channel()
	// if err != nil {
	// 	log.Fatal("Failed to connect channel", err)
	// }
	// defer ch.Close()

	// mq, err := ch.QueueDeclare("sms", false, false, false, false, nil)
	// if err != nil {
	// 	log.Fatal("Failed to declare queue", err)
	// }

	authHandler := AuthHandler{
		Db: db,
		// Ch: ch,
		// Mq: &mq,
	}

	go startRPCServer(":6161", &authHandler)

	http.HandleFunc("/sign-up", authHandler.SignUp)
	// http.HandleFunc("/approve", authHandler.Approve)
	http.HandleFunc("/sign-in", authHandler.SignIn)
	http.HandleFunc("/refresh", authHandler.Refresh)
	http.HandleFunc("/check-user-role", authHandler.CheckUserRole)
	http.HandleFunc("/change-role", authHandler.ChangeRole)

	log.Printf("Auth started on port: %d\n", *port)
	log.Fatal(http.ListenAndServe("auth:"+strconv.Itoa(*port), nil))
}
