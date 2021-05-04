package Routes

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"../DB"
	"github.com/joho/godotenv"
)

func performRequest(r http.Handler, method, path string, param *url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(param.Encode()))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestWelcomePage(t *testing.T) {
	router := SetupRouter()
	param := url.Values{}
	w := performRequest(router, "GET", "/", &param)

	if http.StatusOK != w.Code {
		t.Error()
	}
}

func TestUserRoutes(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()

	router := SetupRouter()
	param := url.Values{}
	w := performRequest(router, "GET", "/", &param)
	if http.StatusOK != w.Code {
		t.Error()
	}

}
func TestUserRoutes1(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()

	router := SetupRouter()
	param := url.Values{}
	param.Add("id", "6090c46f5c814c0d846bdbac")
	w1 := performRequest(router, "GET", "/user", &param)
	if http.StatusOK != w1.Code {
		t.Error()
	}
}
