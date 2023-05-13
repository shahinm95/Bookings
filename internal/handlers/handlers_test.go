package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key string
	value string
}

var theTests = []struct {
	name string
	url string
	method string
	params []postData
	expectedStatusCode int
}{
	{"home", "/","GET", []postData{}, http.StatusOK},
}


func TestHandlers (t *testing.T) {
	routes := getRoutes();
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()
	
	for _, test := range theTests {
		if test.url == "GET" {
			resp , err := testServer.Client().Get(testServer.URL + test.url)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}
				if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("status code for %s url got %d , expected %d status code", test.url, resp.StatusCode, test.expectedStatusCode)
				}
		}
	}
}