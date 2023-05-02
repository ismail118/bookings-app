package dbrepo

import (
	"errors"
	"fmt"
	"github.com/ismail118/bookings-app/internal/models"
	"time"
)

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 0 {
		return errors.New("some error")
	}
	return nil
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, the fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) SearchAvailabilityByRoomID(roomID int, start, end time.Time) (bool, error) {
	if roomID > 2 {
		return false, errors.New("some error")
	}
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	rooms := make([]models.Room, 0)

	if start.Format("2006-01-02") == "2050-01-01" {
		rooms = make([]models.Room, 1)
		return rooms, nil
	}
	if start.Format("2006-01-02") == "2050-01-02" {
		rooms = make([]models.Room, 1)
		return rooms, errors.New("some error")
	}
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, fmt.Errorf("can't find room_id:%d", id)
	}
	return room, nil
}
