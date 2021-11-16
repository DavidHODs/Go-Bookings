package dbrepo

import (
	"time"

	"github.com/DavidHODs/bookings/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database 
func(m *testDBRepo) InsertRoomRestriction(r models.RoomsRestriction) error {
	
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID and false if otherwise
func(m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for a given date range
func(m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID gets a room by ID
func(m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	return room, nil
}