package esutils

import (
	"crypto/md5"
	"fmt"
	"github.com/faabiosr/cachego"
	jsoniter "github.com/json-iterator/go"
	"github.com/techoner/gophp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	ACCESS_TOKEN_KEY = "easywechat.kernel.access_token."
)

var cache cachego.Cache

func SetCache(cache2 cachego.Cache) {
	cache = cache2
}

func GetCache() cachego.Cache {
	return cache
}

func GetToken(a Key) (token KernelAccessToken, err error) {
	key, err := GetCacheKey(a)
	if err != nil {
		return token, err
	}

	tk, err := GetCache().Fetch(key)
	if err != nil {
		return token, err
	}

	real_token, err := gophp.Unserialize([]byte(tk))
	if err != nil {
		return token, err
	}

	b, err := json.Marshal(real_token)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(b, &token)

	token.AccessToken = json.Get(b, a.GetTokenKey()).ToString()

	if err != nil {
		return token, err
	}

	return token, nil
}

func GetCacheKey(a interface{}) (string, error) {
	cre, err := getCredentials(a)
	if err != nil {
		return "", err
	}
	return ACCESS_TOKEN_KEY + fmt.Sprintf("%x", md5.Sum(cre)), nil
}

func getCredentials(a interface{}) (cre []byte, err error) {
	cre, err = json.Marshal(a)
	fmt.Println("json", string(cre))
	if err != nil {
		return nil, err
	}

	return cre, nil
}
