package handlers

import (
	"context"
	"fmt"
	"github.com/ismail118/bookings-app/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	//{"post-search-availability", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2023-01-10"},
	//	{key: "end", value: "2023-01-11"},
	//},
	//	http.StatusOK},
	//{"post-search-availability-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "2023-01-10"},
	//	{key: "end", value: "2023-01-11"},
	//},
	//	http.StatusOK},
	//{"post-make-reservation", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value: "ismail"},
	//	{key: "last_name", value: "alfiyasin"},
	//	{key: "email", value: "ismail@gmail.com"},
	//	{key: "phone_number", value: "021-555-555"},
	//}, http.StatusOK},
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=2050-01-02",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=1",
	)

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=00-00-00",
		"end_date=2050-01-02",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=1",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=00-00-00",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=1",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid room id
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=2050-01-01",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=invalid",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code for invalid room_id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for not found room
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=2050-01-01",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=3",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code can't find room: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid data
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=2050-01-01",
		"first_name=i",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=1",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=2050-01-01",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=2",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code for fail insert reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for fail insert restriction
	reqBody = fmt.Sprintf("&%v&%v&%v&%v&%v&%v&%v",
		"start_date=2050-01-01",
		"end_date=2050-01-01",
		"first_name=ismail",
		"last_name=alfiyasin",
		"email=contoh@gmail.com",
		"phone_number=021111111",
		"room_id=0",
	)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returend wrong response code for fail insert restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
