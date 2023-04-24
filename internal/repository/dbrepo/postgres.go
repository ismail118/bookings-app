package dbrepo

import (
	"context"
	"github.com/ismail118/bookings-app/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
                               created_at, updated_at, restriction_id)
                               values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone,
                          start_date, end_date, room_id, created_at, updated_at)
                          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	row := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.PhoneNumber,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	err := row.Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDBRepo) SearchAvailabilityByRoomID(roomID int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRow int

	query := `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRow)
	if err != nil {
		return false, err
	}

	if numRow == 0 {
		return true, nil
	}

	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
select
    r.id,
    r.room_name
from
    rooms r
where
    r.id not in
    (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`

	rooms := make([]models.Room, 0)
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	err = rows.Err()
	if err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
select id, room_name, created_at, updated_at from rooms where id = $1
`

	var room models.Room

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}
