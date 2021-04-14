package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllPages(t *testing.T) {
	baseUrl := "http://localhost:3000"

	var tests = []struct{
		method string
		url string
		code int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
        {"GET", "/notfound", 404
	},
        {"GET", "/articles", 200},
        {"GET", "/articles/add", 200},
        {"GET", "/articles/4", 200},
        {"GET", "/articles/4/edit", 200},
        {"POST", "/articles/4", 200},
        {"POST", "/articles", 200},
        {"POST", "/articles/1/delete", 404},
	}

	for _, test := range tests {
		t.Logf("请求URL：%v \n", test.url)
		var (
			reps *http.Response
			err error
		)

		switch test.method {
			case "GET":
				reps, err = http.Get(baseUrl + test.url)
			default:
				data := make(map[string][]string)
				reps, err = http.PostForm(baseUrl + test.url, data)
		}

		assert.NoError(t, err, "请求 "+test.url+" 时报错")
		assert.Equal(t, test.code, reps.StatusCode, test.url + "应该返回code:" + strconv.Itoa(test.code))
	}
}