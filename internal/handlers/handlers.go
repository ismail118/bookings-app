package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ismail118/bookings-app/internal/config"
	"github.com/ismail118/bookings-app/internal/forms"
	"github.com/ismail118/bookings-app/internal/models"
	"github.com/ismail118/bookings-app/internal/render"
	"log"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

// RoomOne renders the room one page
func (m *Repository) RoomOne(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "room-one.page.gohtml", &models.TemplateData{})
}

// RoomTwo renders the room two page
func (m *Repository) RoomTwo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "room-two.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.PostForm.Get("start")
	end := r.PostForm.Get("end")

	_, err := w.Write([]byte(fmt.Sprintf("%s %s", start, end)))
	if err != nil {
		log.Println(err)
	}
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// PostAvailabilityJSON handle renders for availability and send JSON response
func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// Reservation renders the make reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName:   r.PostForm.Get("first_name"),
		LastName:    r.PostForm.Get("last_name"),
		Email:       r.PostForm.Get("email"),
		PhoneNumber: r.PostForm.Get("phone_number"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone_number")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary renders the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}
