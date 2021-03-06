package repository

import (
	"time"

	"github.com/DavidHODs/bookings/models"
)

type DatabaseRepo interface {
	AllUsers() 		bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomsRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, password string) (int, string, error)
}