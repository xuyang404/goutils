package HttpClient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)
type Data map[string]interface{}
type File map[string]string

var client = &http.Client{}

type Request struct {
	client            *http.Client
	transport         *http.Transport
	debug             bool
	url               string
	method            string
	time              int64
	timeout           time.Duration
	proxy             string
	username          string
	password          string
	data              interface{}
	disableKeepAlives bool
	tlsClientConfig   *tls.Config
	jar               http.CookieJar
	headers           map[string]string
	cookies           map[string]string
	checkRedirect     func(req *http.Request, via []*http.Request) error
}

func NewRequest() *Request {
	return &Request{
		timeout: 30,
		headers: map[string]string{},
		cookies: map[string]string{},
	}
}

func (r *Request) Debug(debug bool) *Request {
	r.debug = debug
	return r
}

func (r *Request) Proxy(proxy string) *Request {
	r.proxy = proxy
	return r
}

func (r *Request) DisableKeepAlives(b bool) *Request {
	r.disableKeepAlives = b
	return r
}

func (r *Request) SetCheckRedirect(f func(req *http.Request, via []*http.Request) error) *Request {
	r.checkRedirect = f
	return r
}

func (r *Request) SetTlsClient(tls *tls.Config) *Request {
	r.tlsClientConfig = tls
	return r
}

func (r *Request) SetBasicAuth(username string, password string) *Request {
	r.username = username
	r.password = password
	return r
}

func (r *Request) initBasicAuth(req *http.Request) *Request {
	if r.username != "" || r.password != "" {
		req.SetBasicAuth(r.username, r.password)
	}
	return r
}

func (r *Request) SetCookieJar(jar http.CookieJar) *Request {
	r.jar = jar
	return r
}

func (r *Request) SetTimeout(t int) *Request {
	r.timeout = time.Duration(t)
	return r
}

func (r *Request) SetTransport(t *http.Transport) *Request {
	r.transport = t
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

func (r *Request) initHeaders(req *http.Request) *Request {
	if r.headers != nil {
		for k, v := range r.headers {
			req.Header.Set(k, v)
		}
	}
	return r
}

func (r *Request) SetCookies(cookies map[string]string) *Request {
	r.cookies = cookies
	return r
}

func (r *Request) AddCookies(cookies map[string]string) *Request {
	for k, v := range cookies {
		r.cookies[k] = v
	}
	return r
}

func (r *Request) initCookies(req *http.Request) *Request {
	if r.cookies != nil {
		for k, v := range r.cookies {
			req.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}
	return r
}

func (r *Request) Json() *Request {
	r.SetHeaders(map[string]string{"Content-Type": "application/json;charset=utf-8"})
	return r
}

func (r *Request) isJson() bool {
	for _, v := range r.headers {
		if strings.Contains(v, "application/json;charset=utf-8") {
			return true
		}
	}

	return false
}

func (r *Request) getTransport() (http.RoundTripper, error) {
	if r.transport == nil {
		r.transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	if r.proxy != "" {
		purl, err := url.Parse(r.proxy)
		if err != nil {
			return nil, err
		}
		r.transport.Proxy = http.ProxyURL(purl)
	}

	r.transport.DisableKeepAlives = r.disableKeepAlives

	if r.tlsClientConfig != nil {
		r.transport.TLSClientConfig = r.tlsClientConfig
	}

	return http.RoundTripper(r.transport), nil
}

func (r *Request) buildClient() (*Request, error) {
	if r.client == nil {
		t, err := r.getTransport()
		if err != nil {
			return r, err
		}
		//r.client
		//client = &http.Client{
		//	Transport:     t,
		//	CheckRedirect: r.checkRedirect,
		//	Jar:           r.jar,
		//	Timeout:       time.Second * r.timeout,
		//}
		client.Transport = t
		client.CheckRedirect = r.checkRedirect
		client.Jar = r.jar
		client.Timeout = time.Second * r.timeout
		r.client = client
	}
	return r, nil
}

func (r *Request) elapsedTime(t int64, resp *Response) *Request {
	end := time.Now().UnixNano() / 1e6
	resp.time = end - t
	return r
}

func (r *Request) log() {
	if r.debug {
		fmt.Printf("[goutils.HttpClient.Request]\n")
		fmt.Printf("-------------------------------------------------------------------\n")
		fmt.Printf("Request: %s %s\nHeaders: %v\nCookies: %v\nTimeout: %ds\nReqBody: %v\n", r.method, r.url, r.headers, r.cookies, r.timeout, r.data)
		fmt.Printf("-------------------------------------------------------------------\n\n")
	}
}

func parseQuery(reqUrl string) ([]string, error) {
	urlList := strings.Split(reqUrl, "?")
	if len(urlList) < 2 {
		return make([]string, 0), nil
	}

	query := make([]string, 0)
	for _, val := range strings.Split(urlList[1], "&") {
		vals := strings.Split(val, "=")
		if len(vals) < 2 {
			return make([]string, 0), errors.New("query parameter error")
		}

		query = append(query, fmt.Sprintf("%s=%s", vals[0], vals[1]))
	}

	return query, nil
}

func (r *Request) buildUrl(reqUrl string, data Data) (string, error) {
	query, err := parseQuery(reqUrl)
	if err != nil {
		return reqUrl, err
	}

	s := ""
	for k, v := range data {
		if val, ok := v.(string); ok {
			s = val
		} else {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			s = string(b)
		}
		query = append(query, fmt.Sprintf("%s=%s", k, s))
	}

	list := strings.Split(reqUrl, "?")

	if len(query) > 0 {
		return fmt.Sprintf("%s?%s", list[0], strings.Join(query, "&")), nil
	}

	return list[0], nil
}

func (r *Request) buildBody(data Data) (io.Reader, error) {
	if r.method == "GET" || r.method == "DELETE" {
		return nil, nil
	}

	if data == nil {
		return strings.NewReader(""), nil
	}

	if r.isJson() {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return strings.NewReader(string(b)), nil
	}

	body := make([]string, 0)
	s := ""
	for k, v := range data {
		if val, ok := v.(string); ok {
			s = val
		} else {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			s = string(b)
		}

		body = append(body, fmt.Sprintf("%s=%s", k, string(s)))
	}

	return strings.NewReader(strings.Join(body, "&")), nil
}

func (r *Request) request(method string, reqUrl string, data Data) (*Response, error) {
	if method == "" || reqUrl == "" {
		return nil, errors.New("method and url is required")
	}

	resp := &Response{}
	start := time.Now().UnixNano() / 1e6
	defer r.elapsedTime(start, resp)
	defer r.log()

	r.data = data
	r.url = reqUrl
	_, err := r.buildClient()
	if err != nil {
		return nil, err
	}

	r.method = strings.ToUpper(method)
	if r.method == "GET" || r.method == "DELETE" {
		reqUrl, err := r.buildUrl(reqUrl, data)
		if err != nil {
			return nil, err
		}

		r.url = reqUrl
	}

	body, err := r.buildBody(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.method, r.url, body)

	if err != nil {
		return nil, err
	}

	r.initHeaders(req)
	r.initCookies(req)
	r.initBasicAuth(req)

	res, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	resp.url = reqUrl
	resp.Resp = res
	return resp, nil
}

func (r *Request) sendFile(reqUrl string, files File, data Data) (*Response, error) {
	if reqUrl == "" {
		return nil, errors.New("parameter url is required")
	}

	bodyBuffer := &bytes.Buffer{}
	bodyWrite := multipart.NewWriter(bodyBuffer)
	for fieldname, filename := range files {
		fileWrite, err := bodyWrite.CreateFormFile(fieldname, filename)
		if err != nil {
			return nil, err
		}

		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(fileWrite, file)
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	if data != nil {
		for key, value := range data {
			if v, ok := value.(string); ok {
				err := bodyWrite.WriteField(key, v)
				if err != nil {
					return nil, err
				}
			} else {
				b, err := json.Marshal(value)
				if err != nil {
					return nil, err
				}
				err = bodyWrite.WriteField(key, string(b))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	contentType := bodyWrite.FormDataContentType()
	err := bodyWrite.Close()

	if err != nil {
		return nil, err
	}

	resp := &Response{}
	start := time.Now().UnixNano() / 1e6
	defer r.elapsedTime(start, resp)
	defer r.log()

	r.url = reqUrl
	r.data = data
	_, err = r.buildClient()
	if err != nil {
		return nil, err
	}

	r.method = "POST"

	req, err := http.NewRequest(r.method, r.url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	r.initHeaders(req)
	r.initCookies(req)
	r.initBasicAuth(req)
	req.Header.Set("Content-Type", contentType)

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	resp.url = reqUrl
	resp.Resp = res
	return resp, nil
}

func (r *Request) GET(reqUrl string, data Data) (*Response, error) {
	return r.request(http.MethodGet, reqUrl, data)
}

func (r *Request) POST(reqUrl string, data Data) (*Response, error) {
	if _, ok := r.headers["Content-Type"]; !ok {
		r.AddHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	}
	return r.request(http.MethodPost, reqUrl, data)
}

func (r *Request) PUT(reqUrl string, data Data) (*Response, error) {
	return r.request(http.MethodPut, reqUrl, data)
}

func (r *Request) DELETE(reqUrl string, data Data) (*Response, error) {
	return r.request(http.MethodDelete, reqUrl, data)
}

func (r *Request) Upload(reqUrl string, files File, data Data) (*Response, error) {
	return r.sendFile(reqUrl, files, data)
}