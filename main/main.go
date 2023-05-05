package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"spectre-go/controllers"
	"spectre-go/repo"
)

var rdb *redis.Client

func handleRequests(h *controllers.BaseHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/getPassword", h.ProcessGenPasswd).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       0,
	})
	log.Printf("Redis client initialized")

	siteResultRepo := repo.NewSiteResultRepo(client)

	h := controllers.NewBaseHandler(siteResultRepo)

	handleRequests(h)

}
