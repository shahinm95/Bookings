package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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
	// {"post-search-availability", "/search-availability", "Post", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "Post", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
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

func TestRepository_PostReservation(t *testing.T) {

	requestBody := "start_date=2015-01-01"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2015-01-01")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1")

	req, err := http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.PostReservation)
	requestRecorder := httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusSeeOther,
			requestRecorder.Code)
	}

	// test case where there is no form to parse
	req, err = http.NewRequest("POST", "/make-reservation", nil)
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

	// test case where the start data is invalid
	requestBody = "start_date=invalid"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2015-01-01")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()

	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

	// test case where the end data is invalid
	requestBody = "start_date=2015-01-01"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=inavalid")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)
	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

	// test case where the room id is invalid
	requestBody = "start_date=2015-01-01"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2015-01-01")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=invalid")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

	// test case where form is invalid
	requestBody = "start_date=2015-01-01"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2015-01-01")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=J")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusSeeOther,
			requestRecorder.Code)
	}

	// test case where it fails to insert reservation to database
	requestBody = "start_date=2015-01-01"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2015-01-01")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=2")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

	// test case where it fails to insert restriction to database
	requestBody = "start_date=2015-01-01"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "end_date=2015-01-01")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "first_name=John")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Smith")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=john@smith.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=1234649")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "room_id=1000")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		log.Println("error is here", err)
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	requestRecorder = httptest.NewRecorder()
	handler.ServeHTTP(requestRecorder, req)

	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code it should've be %d but got : %d ",
			http.StatusTemporaryRedirect,
			requestRecorder.Code)
	}

}

func TestRepository_ResevationSummary(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Genaral's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(requestRecorder, req)
	if requestRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler retrun wrong response got %d , should be %d",
			requestRecorder.Code, http.StatusOK)
	}

	// test case there is no reservation
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	requestRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(requestRecorder, req)
	if requestRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retrun wrong response got %d , should be %d",
			requestRecorder.Code, http.StatusTemporaryRedirect)
	}
}


func GetCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
