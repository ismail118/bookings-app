package handlers

import (
	"context"
	"fmt"
	"github.com/ismail118/bookings-app/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var Tests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"room-one", "/room-one", "GET", http.StatusOK},
	{"room-two", "/room-two", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-exist-routes", "/this/not/exist", "GET", http.StatusNotFound},
	// new routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"all res", "/admin/reservations/new/1/show", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	route := getRoutes()
	ts := httptest.NewTLSServer(route)
	defer ts.Close()

	for _, e := range Tests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Room one",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test with non-existing room id
	reservation = models.Reservation{
		RoomID: 3,
		Room: models.Room{
			ID:       3,
			RoomName: "Room one",
		},
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	//reqBody := fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=2050-01-02",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=1",
	//)

	reqBody := url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid start date
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=00-00-00",
	//	"end_date=2050-01-02",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=1",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "00-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end date
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=00-00-00",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=1",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "00-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid room id
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=2050-01-01",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=invalid",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for invalid room_id: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for not found room
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=2050-01-01",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=3",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "3")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code can't find room: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid data
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=2050-01-01",
	//	"first_name=i",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=1",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "i")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for fail insert reservation
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=2050-01-01",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=2",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for fail insert reservation: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for fail insert restriction
	//reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
	//	"start_date=2050-01-01",
	//	"end_date=2050-01-01",
	//	"first_name=ismail",
	//	"last_name=alfiyasin",
	//	"email=contoh@gmail.com",
	//	"phone_number=021111111",
	//	"room_id=0",
	//)
	reqBody = url.Values{}
	reqBody.Add("start_date", "2050-01-01")
	reqBody.Add("end_date", "2050-01-02")
	reqBody.Add("first_name", "ismail")
	reqBody.Add("last_name", "alfiyasin")
	reqBody.Add("email", "alfiyasin@gmail.com")
	reqBody.Add("phone_number", "555-555-555")
	reqBody.Add("room_id", "0")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code for fail insert restriction: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	postData := fmt.Sprintf("&%s&%s", "start=2050-01-01", "end=2050-01-02")
	req, _ := http.NewRequest(http.MethodPost, "/search-availability", strings.NewReader(postData))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.PostAvailability)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test missing form
	req, _ = http.NewRequest(http.MethodPost, "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid start date
	postData = fmt.Sprintf("&%s&%s", "start=2050-01-01", "end=00-01-02")
	req, _ = http.NewRequest(http.MethodPost, "/search-availability", strings.NewReader(postData))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid start date
	postData = fmt.Sprintf("&%s&%s", "start=00-01-01", "end=2050-01-02")
	req, _ = http.NewRequest(http.MethodPost, "/search-availability", strings.NewReader(postData))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test fail search available room
	postData = fmt.Sprintf("&%s&%s", "start=2050-01-02", "end=2050-01-02")
	req, _ = http.NewRequest(http.MethodPost, "/search-availability", strings.NewReader(postData))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test no available room
	postData = fmt.Sprintf("&%s&%s", "start=2022-01-02", "end=2022-01-02")
	req, _ = http.NewRequest(http.MethodPost, "/search-availability", strings.NewReader(postData))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostAvailabilityJSON(t *testing.T) {
	postData := fmt.Sprintf("&%s&%s&%s", "start=2050-01-01", "end=2050-01-02", "room_id=1")
	req, _ := http.NewRequest(http.MethodPost, "/search-availability-json", strings.NewReader(postData))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.PostAvailabilityJSON)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test missing form
	req, _ = http.NewRequest(http.MethodPost, "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailabilityJSON)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test fail search availability
	postData = fmt.Sprintf("&%s&%s&%s", "start=2050-01-01", "end=2050-01-02", "room_id=3")
	req, _ = http.NewRequest(http.MethodPost, "/search-availability-json", strings.NewReader(postData))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.PostAvailabilityJSON)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returend wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test missing session reservation
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	uri := "/choose-room/1"
	req, _ := http.NewRequest("GET", uri, nil)
	req.RequestURI = uri
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", models.Reservation{})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid room id
	uri = "/choose-room/invalid"
	req, _ = http.NewRequest("GET", uri, nil)
	req.RequestURI = uri
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", models.Reservation{})

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test missing reservation session
	uri = "/choose-room/1"
	req, _ = http.NewRequest("GET", uri, nil)
	req.RequestURI = uri
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	url := fmt.Sprintf("/book-room?%s&%s&%s",
		"id=1",
		"s=2050-01-01",
		"e=2050-01-02",
	)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid room id
	url = fmt.Sprintf("/book-room?%s&%s&%s",
		"id=invalid",
		"s=2050-01-01",
		"e=2050-01-02",
	)

	req, _ = http.NewRequest(http.MethodGet, url, nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid start date
	url = fmt.Sprintf("/book-room?%s&%s&%s",
		"id=1",
		"s=invalid",
		"e=2050-01-02",
	)

	req, _ = http.NewRequest(http.MethodGet, url, nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid end date
	url = fmt.Sprintf("/book-room?%s&%s&%s",
		"id=1",
		"s=2050-01-01",
		"e=invalid",
	)

	req, _ = http.NewRequest(http.MethodGet, url, nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test can't get room
	url = fmt.Sprintf("/book-room?%s&%s&%s",
		"id=3",
		"s=2050-01-01",
		"e=2050-01-02",
	)

	req, _ = http.NewRequest(http.MethodGet, url, nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

var loginTest = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"ismail@here.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"invalid",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTest {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// make test request recorder
		rr := httptest.NewRecorder()

		// make handler
		handler := http.HandlerFunc(Repo.PostShowLogin)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s : got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			// get the url from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s : location got %s, wanted %s", e.name, actualLoc.String(), e.expectedLocation)
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}

}

var testResCalendars = []struct {
	name                string
	queryParams         []string
	expectationCode     int
	expectationHTML     string
	expectationLocation string
}{
	{
		"test-without-params",
		nil,
		http.StatusOK,
		fmt.Sprintf(`<h3>%s %d</h3>`, time.Now().Month(), time.Now().Year()),
		"",
	},
	{
		"test-with-params",
		[]string{"05", "2050"},
		http.StatusOK,
		fmt.Sprintf(`<h3>%s %d</h3>`, "May", 2050),
		"",
	},
}

func TestRepository_NewAdminReservationsCalendars(t *testing.T) {
	for _, e := range testResCalendars {
		var m, y string
		if e.queryParams != nil {
			m = e.queryParams[0]
			y = e.queryParams[1]
		}
		req, _ := http.NewRequest("GET", fmt.Sprintf("/adimin/reservations-calendar?y=%s&m=%s", y, m), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.NewAdminReservationsCalendars)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectationCode {
			t.Errorf("failed %s : wrong response code, got %d want %d", e.name, rr.Code, e.expectationCode)
		}

		if e.expectationLocation != "" {
			rrLoc, _ := rr.Result().Location()
			if rrLoc.String() != e.expectationLocation {
				t.Errorf("failed %s : wrong location, got %s want %s", e.name, rrLoc.String(), e.expectationLocation)
			}
		}

		if e.expectationHTML != "" {
			if !strings.Contains(rr.Body.String(), e.expectationHTML) {
				t.Errorf("failed %s: expected to find %s", e.name, e.expectationHTML)
			}
		}
	}
}

var testPostResCalendars = []struct {
	name            string
	bodyReq         string
	stringMap       map[string]int
	expectationCode int
}{
	{
		"test-delete-reservation",
		`m=05&y=2023`,
		map[string]int{"2023-05-1": 2},
		http.StatusSeeOther,
	},
	{
		"test-add-block",
		`m=05&y=2023&add_block_1_2023-05-01=1`,
		map[string]int{},
		http.StatusSeeOther,
	},
}

func TestRepository_NewAdminPostReservationsCalendars(t *testing.T) {
	for _, e := range testPostResCalendars {
		req, _ := http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(e.bodyReq))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		session.Put(req.Context(), "block_map_1", e.stringMap)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.NewAdminPostReservationsCalendars)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectationCode {
			t.Errorf("failed %s : wrong response code, got %d want %d", e.name, rr.Code, e.expectationCode)
		}
	}
}

var testProcRes = []struct {
	name                string
	params              []string
	expectationCode     int
	expectationLocation string
}{
	{
		"test-process-reservation-all",
		[]string{"all", "1", "", ""},
		http.StatusSeeOther,
		"/admin/reservations-all",
	},
	{
		"test-process-reservation-calendar",
		[]string{"calendar", "1", "2023", "05"},
		http.StatusSeeOther,
		"/admin/reservations-calendar",
	},
}

func TestRepository_AdminProcessReservation(t *testing.T) {
	for _, e := range testProcRes {
		uri := fmt.Sprintf("/admin/process-reservation/%s/%s/do??y=%s&m=%s", e.params[0], e.params[1], e.params[2], e.params[3])
		req, _ := http.NewRequest("GET", uri, nil)
		req.RequestURI = uri
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectationCode {
			t.Errorf("failed %s : wrong response code, got %d want %d", e.name, rr.Code, e.expectationCode)
		}

		if e.expectationLocation != "" {
			rrLoc, _ := rr.Result().Location()
			if rrLoc.String() != e.expectationLocation {
				t.Errorf("failed %s : wrong location, got %s want %s", e.name, rrLoc.String(), e.expectationLocation)
			}
		}
	}
}

var testPostShowRes = []struct {
	name                string
	params              []string
	formData            string
	expectationCode     int
	expectationLocation string
}{
	{
		"post-with-year",
		[]string{"cal", "5"},
		fmt.Sprintf("&%s&%s&%s&%s&%s&%s",
			"first_name=ismail",
			"last_name=alfiyasin",
			"email=me@here.com",
			"phone_number=0121121112",
			"year=2023",
			"month=05",
		),
		http.StatusSeeOther,
		"/admin/reservations-calendar?y=2023&m=05",
	},
	{
		"post-without-year",
		[]string{"all", "5"},
		fmt.Sprintf("&%s&%s&%s&%s&%s&%s",
			"first_name=ismail",
			"last_name=alfiyasin",
			"email=me@here.com",
			"phone_number=0121121112",
			"year=",
			"month=",
		),
		http.StatusSeeOther,
		"/admin/reservations-all",
	},
}

func TestRepository_AdminPostShowReservation(t *testing.T) {
	for _, e := range testPostShowRes {
		uri := fmt.Sprintf("/admin/reservations/%s/%s", e.params[0], e.params[1])
		req, _ := http.NewRequest("POST", uri,
			strings.NewReader(e.formData))
		req.RequestURI = uri
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminPostShowReservation)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectationCode {
			t.Errorf("failed %s : wrong response code, got %d want %d", e.name, rr.Code, e.expectationCode)
		}

		if e.expectationLocation != "" {
			rrLoc, _ := rr.Result().Location()
			if rrLoc.String() != e.expectationLocation {
				t.Errorf("failed %s : wrong location, got %s want %s", e.name, rrLoc.String(), e.expectationLocation)
			}
		}
	}
}

var testDeleteRes = []struct {
	name                string
	params              []string
	expectationCode     int
	expectationLocation string
}{
	{
		"delete-from-calendar",
		[]string{"cal", "1", "2023", "05"},
		http.StatusSeeOther,
		"/admin/reservations-calendar?y=2023&m=05",
	},
	{
		"delete-from-new",
		[]string{"new", "1", "", ""},
		http.StatusSeeOther,
		"/admin/reservations-new",
	},
}

func TestRepository_AdminDeleteReservation(t *testing.T) {
	for _, e := range testDeleteRes {
		uri := fmt.Sprintf("/admin/delete-reservation/%s/%s/do?y=%s&m=%s", e.params[0], e.params[1], e.params[2], e.params[3])
		req, _ := http.NewRequest("GET", uri, nil)
		req.RequestURI = uri
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectationCode {
			t.Errorf("failed %s : wrong response code, got %d want %d", e.name, rr.Code, e.expectationCode)
		}

		if e.expectationLocation != "" {
			rrLoc, _ := rr.Result().Location()
			if rrLoc.String() != e.expectationLocation {
				t.Errorf("failed %s : wrong location, got %s want %s", e.name, rrLoc.String(), e.expectationLocation)
			}
		}
	}
}
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
