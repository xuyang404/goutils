package esutils

type KernelAccessToken struct {
	AccessToken string
	ExpiresIn   int64 `json:"expires_in"`
}

type Token interface {
	GetToken() (token KernelAccessToken, err error)
	GetCacheKey() (string, error)
}

type Key interface {
	GetTokenKey() (a interface{})
}