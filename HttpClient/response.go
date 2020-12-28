package HttpClient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	time int64
	url  string
	Resp *http.Response
	body []byte
}

func (r *Response) Time() string {
	if r != nil {
		return fmt.Sprintf("%dms", r.time)
	}
	return "0ms"
}

func (r *Response) Url() string {
	if r != nil {
		return r.url
	}
	return ""
}

func (r *Response) StatusCode() int {
	if r == nil || r.Resp == nil {
		return 0
	}
	return r.Resp.StatusCode
}

func (r *Response) Cookies() []*http.Cookie {
	if r == nil || r.Resp == nil {
		return nil
	}
	return r.Resp.Cookies()
}

func (r *Response) Headers() http.Header {
	if r == nil || r.Resp == nil {
		return nil
	}
	return r.Resp.Header
}

func (r *Response) Body() ([]byte, error) {

	if r == nil {
		return nil,errors.New("goutils.HttpClient.Response is nil")
	}

	if len(r.body) > 0 {
		return r.body, nil
	}

	if r.Resp == nil || r.Resp.Body == nil {
		return nil, errors.New("response or body is nil")
	}

	defer r.Resp.Body.Close()

	b, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		return nil, err
	}

	r.body = b

	return b, nil
}

func (r *Response) Content() (string, error) {
	b, err := r.Body()
	if err != nil {
		return "", err
	}

	return string(b), err
}

func (r *Response) Close() error {
	return r.Resp.Body.Close()
}

func (r *Response) Json(v interface{}) error {
	b, err := r.Body()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &v); err != nil {
		return err
	}
	return nil
}
