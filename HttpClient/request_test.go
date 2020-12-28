package HttpClient_test

import (
	"crypto/md5"
	"fmt"
	"github.com/xuyang404/goutils/HttpClient"
	"github.com/xuyang404/goutils/gpool"
	"sync"
	"testing"
)

func TestRequest_SendField(t *testing.T) {
	req := HttpClient.NewRequest()

	resp, err := req.Debug(true).
		Upload("http://127.0.0.1:8000/wechat/test/test",
			HttpClient.File{
				"file[0]": "../go.mod",
				"file[1]": "../README.md",
			},
			HttpClient.Data{
				"a[0]": []string{"1", "2", "3"},
				"a[1]": "2",
				"a[2]": "3",
			},
		)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content())
}

func TestRequest_GET(t *testing.T) {
	req := HttpClient.NewRequest()
	resp, err := req.Debug(true).GET(
		"http://127.0.0.1:8000/wechat/test/test?test=abc",
		HttpClient.Data{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content())
}

func TestRequest_POST(t *testing.T) {
	req := HttpClient.NewRequest()
	resp, err := req.Debug(true).POST(
		"http://127.0.0.1:8000/wechat/test/test",
		HttpClient.Data{
			"a[0]": 1,
			"a[1]": 2,
			"a[2]": 3,
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content())
}

func TestRequest_PUT(t *testing.T) {
	req := HttpClient.NewRequest()
	resp, err := req.Debug(true).PUT(
		"http://127.0.0.1:8000/wechat/test/test",
		HttpClient.Data{
			"a[0]": 1,
			"a[1]": 2,
			"a[2]": 3,
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content())
}

func TestRequest_DELETE(t *testing.T) {
	req := HttpClient.NewRequest()
	resp, err := req.Debug(true).DELETE(
		"http://127.0.0.1:8000/wechat/test/test?test=abc",
		HttpClient.Data{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content())
}

func TestRequest_Json(t *testing.T) {
	req := HttpClient.NewRequest()
	resp, err := req.Debug(true).Json().POST(
		"http://127.0.0.1:8000/wechat/test/test",
		HttpClient.Data{
			"a": []int{1, 2, 3, 4},
			"b": HttpClient.Data{"t1": "1", "t2": 2},
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Content())
}

func TestRequest_SetHeaders(t *testing.T) {
	req := HttpClient.NewRequest()

	req.SetHeaders(map[string]string{"1": "2"})
	resp, err := req.Debug(true).POST("http://127.0.0.1:8000/wechat/test/test", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Content())
}

func TestRequest_SetCookies(t *testing.T) {
	req := HttpClient.NewRequest()

	req.SetCookies(map[string]string{"1": "2"})
	resp, err := req.Debug(true).POST("http://127.0.0.1:8000/wechat/test/test", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Content())
}

func TestRequest_SetTimeout(t *testing.T) {
	req := HttpClient.NewRequest()

	resp, err := req.Debug(true).SetTimeout(10).POST("http://127.0.0.1:8000/wechat/test/test", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Content())
}

func TestRequest_Proxy(t *testing.T) {
	req := HttpClient.NewRequest()
	req.Proxy("http://127.0.0.1:8081")
	resp, err := req.Debug(true).POST("http://www.baidu.com", HttpClient.Data{"123": "123"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Content())
}

func Task(wg *sync.WaitGroup, i int) func() {
	return func() {
		defer wg.Done()
		req := HttpClient.NewRequest()
		resp, err := req.Debug(true).POST(
			"https://mp.weixin.qq.com/cgi-bin/bizlogin?action=startlogin",
			HttpClient.Data{
				"username":     i,
				"pwd":          fmt.Sprintf("%x", md5.Sum([]byte("12312"))),
				"imgcode":      "",
				"f":            "json",
				"userlang":     "zh_CN",
				"redirect_url": "",
				"token":        "",
				"lang":         "zh_CN",
				"ajax":         "1",
			},
		)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp.Content())
	}
}

func TestClient(t *testing.T) {
	pool := gpool.NewPool(3)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		pool.Add(Task(wg, i))
	}

	wg.Wait()

	fmt.Println(123)
}
