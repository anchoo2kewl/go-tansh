package models

import (
	"database/sql"
)

type Rsvps struct {
	ID           int
	EventId      int
	GuestName    string
	Email        string
	NumberGuests string
	CreatedAt    string
}

type RsvpsList struct {
	Rsvps []Rsvps
}

type RsvpService struct {
	DB *sql.DB
}

func (ss *RsvpService) GetRsvps() (RsvpsList, error) {
	list := RsvpsList{}

	rows, err := ss.DB.Query("SELECT * FROM rsvps ORDER BY ID DESC")
	if err != nil {
		return list, nil
	}

	for rows.Next() {

		var rsvp Rsvps
		err := rows.Scan(&rsvp.ID, &rsvp.EventId, &rsvp.GuestName, &rsvp.Email, &rsvp.NumberGuests, &rsvp.CreatedAt)
		if err != nil {
			panic(err)
		}

		list.Rsvps = append(list.Rsvps, rsvp)
	}

	return list, nil
}
