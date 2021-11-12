package repository

import (
	"github.com/DavidHODs/bookings/models"
)

type DatabaseRepo interface {
	AllUsers() 		bool
	InsertReservation(res models.Reservation) error
}