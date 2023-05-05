package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"spectre-go/controllers"
	"spectre-go/repo"
)

var rdb *redis.Client

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       0,
	})
	log.Printf("Redis client initialized")

	siteResultRepo := repo.NewSiteResultRepo(client)

	h := controllers.NewBaseHandler(siteResultRepo)

	log.Fatal(http.ListenAndServe(":10000", controllers.HandleRequests(h)))
}
