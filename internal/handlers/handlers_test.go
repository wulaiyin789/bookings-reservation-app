package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var theTests = []struct {
	name   string
	url    string
	method string
	// params             []postData
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

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
			// } else {
			// 	values := url.Values{}
			// 	for _, x := range e.params {
			// 		values.Add(x.key, x.value)
			// 	}

			// 	resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			// 	if err != nil {
			// 		t.Log(err)
			// 		t.Fatal(err)
			// 	}

			// 	if resp.StatusCode != e.expectedStatusCode {
			// 		t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			// 	}
		}
	}
}

func TestRepository_GetReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}

	// get the context with session from the request
	ctx := getCtx(req)

	// Add to the request
	req = req.WithContext(ctx)

	// Simulate what we get from the request/response lifecycle
	rr := httptest.NewRecorder()

	// Put reservation to the session
	session.Put(ctx, "reservation", reservation)

	// Take the reservation variable and turn it to a function
	handler := http.HandlerFunc(Repo.Reservation)

	// Serve it as a HTTP
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	// overwrite the room id for test fail
	reservation.RoomID = 100

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	//* Appropriate date format (time.Time) to put into reservation struct
	sd, _ := time.Parse("2006-01-02", "2050-01-01")
	ed, _ := time.Parse("2006-01-02", "2050-01-02")

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
		StartDate: sd,
		EndDate:   ed,
	}

	// we only need these values since they come from form user inputs, and other ones are retrieved from session
	reqBody := fmt.Sprintf("%s&%s&%s&%s",
		"first_name=John",
		"last_name=Smith",
		"email=John@test.com",
		"phone=1234567890",
	)

	// postedData := url.Values{}
	// postedData.Add("start_date", "2050-01-01")
	// postedData.Add("end_date", "2050-01-02")
	// postedData.Add("first_name", "John")
	// postedData.Add("last_name", "Smith")
	// postedData.Add("email", "John@test.com")
	// postedData.Add("phone", "1234567890")

	// basic test
	// req, _ := http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(postedData.Encode()))
	req, _ := http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(reqBody))

	// get the context with session from the request
	ctx := getCtx(req)

	// Add to the request
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// simulate what we get from the request/response lifecycle
	rr := httptest.NewRecorder()

	// put reservation into session required to satisfy handler
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//* Test without request body (ParseForm fail)
	req, _ = http.NewRequest(http.MethodPost, "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//* Test without reservation struct in session
	req, _ = http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//* Test when violating one of validations (first_name < 3 chars)
	reqBody = fmt.Sprintf("%s&%s&%s&%s",
		"first_name=Jh",
		"last_name=Smith",
		"email=John@test.com",
		"phone=1234567890",
	)

	req, _ = http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//* Test when violating InsertReservation
	reqBody = fmt.Sprintf("%s&%s&%s&%s",
		"first_name=Invalid",
		"last_name=Smith",
		"email=John@test.com",
		"phone=1234567890",
	)

	req, _ = http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//* Test when violating InsertRoomRestriction
	reservation.RoomID = 50

	reqBody = fmt.Sprintf("%s&%s&%s&%s",
		"first_name=John",
		"last_name=Smith",
		"email=John@test.com",
		"phone=1234567890",
	)

	req, _ = http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-01")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2050-01-01")

	// create request
	req, _ := http.NewRequest(http.MethodPost, "/make-reservation", strings.NewReader(reqBody))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set hte request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// maker handler HandlerFun
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
