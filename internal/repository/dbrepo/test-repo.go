package dbrepo

import (
	"errors"
	"time"

	"github.com/shahinm95/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if room id == 2 fail otherwise pass
	if res.RoomID== 2 {
		return 0 , errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inset a room restriction to database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID== 1000{
		return errors.New("some error")
	}
	return nil
}

// SearchAvailability return true room availability exits for given room ID otherswise it return false
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any , for given date
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID  gets a room by ID
func (m *testDBRepo) GetRoomByID (id int)(models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("can't get room by ID")
	}
	return room , nil
}