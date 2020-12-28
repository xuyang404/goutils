package OpenPlatform

import (
	"errors"
	"github.com/techoner/gophp"
	"github.com/xuyang404/goutils/esutils"
	"strings"
)

type Credentials struct {
	Appid        string `json:"component_appid"`
	AuthAppid    string `json:"authorizer_appid"`
	RefreshToken string `json:"authorizer_refresh_token"`
	tokenKey     string
	keyNotExistExpireHandler func(credentials *Credentials)
	ticketKeyNotExistExpireHandler func(credentials *Credentials)
}

var (
	TICKET_NOT_EXIST = errors.New(`Credential "component_verify_ticket" does not exist in cache.`)
)

const TICKET_KEY = "easywechat.open_platform.verify_ticket."

func NewCredentials(appid, authAppid, refreshToken string) *Credentials {
	return &Credentials{
		Appid:        appid,
		AuthAppid:    authAppid,
		RefreshToken: refreshToken,
		tokenKey: "authorizer_access_token",
	}
}

func (a *Credentials) GetTicket() (ticket string, err error) {
	ticket,err = a.getTicket()
	if err != nil && strings.Contains(err.Error(), "not exist") {
		if a.ticketKeyNotExistExpireHandler != nil {
			a.ticketKeyNotExistExpireHandler(a)
			return a.getTicket()
		}
	}

	return ticket,err
}

func (a *Credentials)getTicket()(ticket string, err error)  {
	ticket, err = esutils.GetCache().Fetch(TICKET_KEY + a.Appid)
	if err != nil {
		return "", err
	}

	if ticket == "" {
		return "", TICKET_NOT_EXIST
	}

	real_ticket, err := gophp.Unserialize([]byte(ticket))
	if err != nil {
		return "", err
	}
	return real_ticket.(string), nil
}

func (a *Credentials) SetTokenKeyNotExistExpireHandler(f func(credentials *Credentials))  {
	a.keyNotExistExpireHandler = f
}

func (a *Credentials) SetTicketKeyNotExistExpireHandler(f func(credentials *Credentials))  {
	a.ticketKeyNotExistExpireHandler = f
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
