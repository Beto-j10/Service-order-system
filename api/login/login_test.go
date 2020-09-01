package login_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestLogin(t *testing.T) {
	formData := url.Values{
		"email":    {"linda@dominio.com"},
		"password": {"234"},
	}

	resp, err := http.PostForm("http://localhost:8080/login", formData)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
