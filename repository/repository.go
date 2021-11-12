package repository

import (
	"time"

	"github.com/DavidHODs/bookings/models"
)

type DatabaseRepo interface {
	AllUsers() 		bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomsRestriction) error
	SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error)
}