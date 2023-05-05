package repo

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type SiteResultRepo struct {
	db *redis.Client
}

func NewSiteResultRepo(db *redis.Client) *SiteResultRepo {

	return &SiteResultRepo{
		db: db,
	}
}

func (repo *SiteResultRepo) Save(hashedKey string, password string) {
	ttl := time.Duration(24) * time.Hour

	op1 := repo.db.Set(context.Background(), hashedKey, password, ttl)
	if err := op1.Err(); err != nil {
		fmt.Printf("unable to SET data. error: %v", err)
		return
	}
	log.Println("set operation success")
}

func (repo *SiteResultRepo) FindSiteResult(hashedKey string) string {
	op2 := repo.db.Get(context.Background(), hashedKey)

	if err := op2.Err(); err != nil {
		log.Printf("unable to GET data. error: %v", err)
		return ""
	}
	res, err := op2.Result()
	if err != nil {
		log.Printf("unable to GET data. error: %v", err)
		return ""
	}
	return res
}
