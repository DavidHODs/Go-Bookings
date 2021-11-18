package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DavidHODs/bookings/config"
	"github.com/DavidHODs/bookings/driver"
	"github.com/DavidHODs/bookings/forms"
	"github.com/DavidHODs/bookings/helpers"
	"github.com/DavidHODs/bookings/models"
	"github.com/DavidHODs/bookings/render"
	"github.com/DavidHODs/bookings/repository"
	"github.com/DavidHODs/bookings/repository/dbrepo"
	"github.com/go-chi/chi"
)

// struct for responses in json format
type jsonResponse struct {
	OK 			bool	`json:"ok"`
	Message		string	`json:"message"`
	RoomID		string	`json:"room_id"`
	StartDate	string	`json:"start_date"`
	EndDate		string	`json:"end_date"`
}

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB	repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository {
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository for testing
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository {
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandler sets the repositories for the handlers
func NewHandlers(r *Repository){
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home_page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// sends the data to the template
	render.Template(w, r, "about_page.gohtml", &models.TemplateData{}) 
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Rooms.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "make-reservation_page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get reservation from session"))
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Has("first_name")
	form.Has("last_name")
	form.Has("email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation_page.gohtml", &models.TemplateData{
		Form: form,
		Data: data,
	})
	return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomsRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	htmlMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong>
		<br>
		Dear %s, <br>
		This is to confirm your %s reservation from %s to %s.
	`, reservation.FirstName, reservation.Rooms.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	ownerMessage := fmt.Sprintf(
		`<strong>Reservation Notification</strong>
		<br>
		Dear %s, <br>
		This is to inform you that your property, %s has been booked from %s to %s.
		`, "Owner", reservation.Rooms.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	// send notification to the guest 
	msg := models.MailData{
		To:      	reservation.Email,
		From:    	"me@here.com",
		Subject: 	"Reservation Confirmation",
		Content: 	htmlMessage,
		Template:	"basic.html",
	}

	m.App.MailChan <- msg

	// send notification to property owner 
	msg = models.MailData{
		To:      "owner@property.com",
		From:    "me@here.com",
		Subject: "Reservation Notification",
		Content: ownerMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the General's quarters page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals_page.gohtml", &models.TemplateData{})
}

// Generals renders the Majors's suites page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors_page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "reserve_page.gohtml", &models.TemplateData{})
}

// PostAvailability posts to availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// if startDate != time.Now().Local() {
	// 	m.App.Session.Put(r.Context(), "error", "arrival date can't be less than current date")
	// 	http.Redirect(w, r, "/", http.StatusBadRequest)
	// 	return
	// }

	// if endDate.Before(startDate) {
	// 	m.App.Session.Put(r.Context(), "error", "departure date can't be less than arrival date")
	// 	http.Redirect(w, r, "/", http.StatusBadRequest)
	// 	return
	// }

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	for _, i := range rooms {
		m.App.InfoLog.Println("Room:", i.ID, i.RoomName)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room-page.gohtml", &models.TemplateData{Data: data})
}

// AvailabilityJSON sends json response on availability request
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	
	resp := jsonResponse{
		OK: false,
		Message: "Internal Server Error",
	}

	out, err := json.MarshalIndent(resp, "", "	")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact_page.gohtml", &models.TemplateData{})
}

// ReservationSummary renders the summary of the reservation details
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get error from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary-page.gohtml", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}

func(m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}