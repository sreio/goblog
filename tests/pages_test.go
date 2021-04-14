package tests

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	baseUrl := "http://localhost:3000"

	var (
		resp *http.Response
		err error
	)
	resp, err = http.Get(baseUrl + "/")
	
	assert.NoError(t, err, "有错误发生时，err不为空")
	assert.Equal(t, 200, resp.StatusCode, "状态码应该为200")
}