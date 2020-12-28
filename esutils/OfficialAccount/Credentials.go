package OfficialAccount

import (
	"github.com/xuyang404/goutils/esutils"
	"strings"
)

type Credentials struct {
	GrantType string `json:"grant_type"`
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	tokenKey  string
	keyNotExistExpireHandler func(credentials *Credentials)
}

func NewCredentials(appid, secret string) *Credentials {
	return &Credentials{
		Appid:     appid,
		Secret:    secret,
		GrantType: "client_credential",
		tokenKey: "access_token",
	}
}

func (a *Credentials) SetTokenKeyNotExistExpireHandler(f func(credentials *Credentials))  {
	a.keyNotExistExpireHandler = f
}

func (a *Credentials) GetToken() (token esutils.KernelAccessToken, err error) {
	tk,err := esutils.GetToken(a)
	if err != nil && strings.Contains(err.Error(), "not exist") {
		if a.keyNotExistExpireHandler != nil {
			a.keyNotExistExpireHandler(a)
			return esutils.GetToken(a)
		}
	}
	return tk,err
}

func (a *Credentials) GetCacheKey() (string, error) {
	return esutils.GetCacheKey(a)
}

func (a *Credentials) GetTokenKey() interface{} {
	return a.tokenKey
}
