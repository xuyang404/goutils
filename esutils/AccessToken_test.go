package esutils_test

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xuyang404/goutils/esutils"
	"github.com/xuyang404/goutils/esutils/OfficialAccount"
	"github.com/xuyang404/goutils/esutils/OpenPlatform"
	"log"
	"testing"
)

func getRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB,
		MinIdleConns: 5,
		PoolSize:     10,
	})
	return rdb
}

func TestOpenPlatform(t *testing.T) {
	rdb := getRedis()
	cache := esutils.NewDefaultCache(rdb)
	esutils.SetCache(cache)
	manager := OpenPlatform.NewCredentials("", "", "")
	fmt.Println(manager.GetCacheKey())
	ticket, err := manager.GetTicket()
	if err != nil {
		log.Fatal(err)
	}

	token, err := manager.GetToken()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ticket", ticket)
	fmt.Println("token", token.AccessToken)
	fmt.Println("expire_in", token.ExpiresIn)
}

func TestOfficialAccount(t *testing.T) {
	rdb := getRedis()
	cache := esutils.NewDefaultCache(rdb)
	esutils.SetCache(cache)
	manager := OfficialAccount.NewCredentials("", "")
	manager.SetTokenKeyNotExistExpireHandler(func(credentials *OfficialAccount.Credentials) {
		fmt.Println("刷新key")
	})
	token, err := manager.GetToken()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("token", token.AccessToken)
}
