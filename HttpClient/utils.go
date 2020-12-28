package HttpClient

import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Json() jsoniter.API {
	return json
}