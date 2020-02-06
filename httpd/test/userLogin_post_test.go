package test

import (
	"bytes"
	"inventory/httpd/initRouter"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLoginPost(t *testing.T) {
	email := "gheugrehiu@fwer.qq"
	value := url.Values{}
	value.Add("email", email)
	value.Add("password", "555")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/userLogin", bytes.NewBufferString(value.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	router := initRouter.SetupRouter()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, strings.Contains(w.Body.String(), email), true)
}
