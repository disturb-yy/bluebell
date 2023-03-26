package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) // 设置为测试模式
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id": 1,
		"title": "test",
		"content": "just a test",
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// 判断相应内容是否按照预期返回了需要登录的错误

	// 1. 判断响应内容是否包含指定字符串
	//assert.Contains(t, w.Body.String(), "需要登录")

	// 2. 将响应的内容反序化到ResponseData，然后判断与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err: %v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)
}
