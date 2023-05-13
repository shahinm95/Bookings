package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "Post", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", "Post", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "Post", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}



func TestHandlers (t *testing.T) {
	routes := getRoutes();
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()
	
	for _, test := range theTests {
		if test.method == "GET" {
			resp , err := testServer.Client().Get(testServer.URL + test.url)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}
				if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("status code for %s url got %d , expected %d status code", test.url, resp.StatusCode, test.expectedStatusCode)
				}
		} else {
			values := url.Values{} // data type needed for PostForm()
			for _, x := range test.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL + test.url, values)
			if err != nil {
					t.Log(err)
					t.Fatal(err, resp)
				}
			if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("status code for %s url got %d , expected %d status code", test.url, resp.StatusCode, test.expectedStatusCode)
				}
		}
			
	}
}