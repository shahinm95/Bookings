package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shahinm95/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	// {"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"post-search-availability", "/search-availability", "Post", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "Post", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"make-reservation", "/make-reservation", "Post", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "555-555-5555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + test.url)
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

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Genaral's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	requestRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusOK {
		t.Error("Reservation handler retrun wrong response code it should've be 200 but got : ", requestRecorder.Code)
	}

	// test case where reservation is defined is session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	requestRecorder = httptest.NewRecorder()

	handler.ServeHTTP(requestRecorder, req)
	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retrun wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

	// text case where there is no existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	requestRecorder = httptest.NewRecorder()

	reservation.RoomID = 3
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retrun wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}
}

func GetCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
