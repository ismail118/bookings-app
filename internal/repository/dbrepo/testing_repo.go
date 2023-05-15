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

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "ismail@here.com" {
		return 0, "", errors.New("some error")
	}
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var r []models.Reservation
	return r, nil
}

func (m *testDBRepo) NewReservations() ([]models.Reservation, error) {
	var r []models.Reservation
	return r, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	rooms = append(rooms, models.Room{
		ID:        1,
		RoomName:  "room test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return rooms, nil
}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	if start.Format("2006-01-02") == "2050-05-01" {
		return []models.RoomRestriction{
			{
				ID:            1,
				StartDate:     start,
				EndDate:       start.AddDate(0, 0, 1),
				RoomID:        roomID,
				RestrictionID: 2,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}, nil
	}
	return []models.RoomRestriction{}, nil
}

func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	return nil
}

func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
