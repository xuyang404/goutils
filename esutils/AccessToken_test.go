package esutils_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xuyang404/goutils/esutils"
	"github.com/xuyang404/goutils/esutils/OfficialAccount"
	"github.com/xuyang404/goutils/esutils/OpenPlatform"
	"io/ioutil"
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
	b,err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Fatal(err)
	}
	data := make(map[string]string)
	json.Unmarshal(b, &data)
	fmt.Println(data)
	manager := OpenPlatform.NewCredentials(data["appid"], data["authAppid"], data["refreshToken"])
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
