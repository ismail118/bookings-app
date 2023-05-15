package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ismail118/bookings-app/internal/config"
	"github.com/ismail118/bookings-app/internal/driver"
	"github.com/ismail118/bookings-app/internal/forms"
	"github.com/ismail118/bookings-app/internal/models"
	"github.com/ismail118/bookings-app/internal/render"
	"github.com/ismail118/bookings-app/internal/repository"
	"github.com/ismail118/bookings-app/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send data to the template
	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

// RoomOne renders the room one page
func (m *Repository) RoomOne(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "room-one.page.gohtml", &models.TemplateData{})
}

// RoomTwo renders the room two page
func (m *Repository) RoomTwo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "room-two.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't search availability")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "chose-rooms.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// PostAvailabilityJSON handle renders for availability and send JSON response
func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal Server Error",
		}

		out, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)

	endDate, _ := time.Parse(layout, ed)

	rd := r.Form.Get("room_id")
	roomID, _ := strconv.Atoi(rd)

	available, err := m.DB.SearchAvailabilityByRoomID(roomID, startDate, endDate)
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Cant search availability on database",
		}

		out, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	res := jsonResponse{
		Ok:        available,
		Message:   "",
		RoomID:    rd,
		StartDate: sd,
		EndDate:   ed,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// Reservation renders the make reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	layout := "2006-01-02"
	sd := res.StartDate.Format(layout)
	ed := res.EndDate.Format(layout)
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	res.Room.RoomName = room.RoomName
	data["reservation"] = res

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse room_id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		FirstName:   r.Form.Get("first_name"),
		LastName:    r.Form.Get("last_name"),
		Email:       r.Form.Get("email"),
		PhoneNumber: r.Form.Get("phone_number"),
		StartDate:   startDate,
		EndDate:     endDate,
		RoomID:      roomID,
		Room:        room,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone_number")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		// add these lines to fix bad data error
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		http.Error(w, "my own error", http.StatusSeeOther)

		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br>
	Dear %s <br>
	This is confirm your reservation from %s to %s.

`, reservation.FirstName, sd, ed)

	// send email notification
	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	htmlMessage = fmt.Sprintf(`
	Dear Owner <br>
	You got new reservation from %s to %s.<br>
`, sd, ed)
	msg = models.MailData{
		To:       "owner@gmail.com",
		From:     "me@here.com",
		Subject:  "Reservation Coming Alert",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary renders the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")

	layout := "2006-01-02"
	sd := reservation.StartDate.Format(layout)
	ed := reservation.EndDate.Format(layout)
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid room id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid room id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "can't get room")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.IsEmail("email")
	form.Required("email", "password")
	if !form.Valid() {
		render.Template(w, r, "login.page.gohtml", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Success Login")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.gohtml", &models.TemplateData{})
}

func (m *Repository) NewAdminReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.NewReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get all reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewAdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get all reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewAdminReservationsCalendars(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	yr := r.URL.Query().Get("y")
	mo := r.URL.Query().Get("m")

	if yr != "" && mo != "" {
		year, _ := strconv.Atoi(yr)
		month, _ := strconv.Atoi(mo)
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get all rooms")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "can't get restrictions")
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// it's a reservations
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				// it's a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calendars.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

func (m *Repository) NewAdminPostReservationsCalendars(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	month, _ := strconv.Atoi(r.Form.Get("m"))
	year, _ := strconv.Atoi(r.Form.Get("y"))

	// process blocks
	rooms, err := m.DB.AllRooms()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get all rooms")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {
		// get the block map from the session. loop through entire map, if we have an entry in the map
		// that does not exist in our posted data, and if the restriction id > 0, then it is a block we need to remove
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			// ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				// only pay attention to values > 0, and that are not in the form post
				// the rest are just placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						// delete the restriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	// now handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			// insert new block
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}

func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	reservationId, err := strconv.Atoi(exploded[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse to int")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := exploded[3]

	err = m.DB.UpdateProcessedForReservation(reservationId, 1)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't update reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year != "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	reservationId, err := strconv.Atoi(exploded[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse to int")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["year"] = year
	stringMap["month"] = month

	reservation, err := m.DB.GetReservationByID(reservationId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "admin-reservation-show.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	reservationId, err := strconv.Atoi(exploded[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse to int")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationByID(reservationId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.PhoneNumber = r.Form.Get("phone_number")

	err = m.DB.UpdateReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	year := r.Form.Get("year")
	month := r.Form.Get("month")

	m.App.Session.Put(r.Context(), "flash", "Changes Save")

	if year != "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	reservationId, err := strconv.Atoi(exploded[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse to int")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := exploded[3]

	err = m.DB.DeleteReservation(reservationId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't update reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	month := r.URL.Query().Get("m")
	year := r.URL.Query().Get("y")

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	if year != "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
