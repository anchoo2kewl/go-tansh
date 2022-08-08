package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RsvpsResource struct {
	GuestName    string `json:"guest_name"`
	EventId      int    `json:"event_id"`
	Email        string `json:"email"`
	Id           int    `json:"id"`
	NumberGuests int    `json:"number_guests"`
	CreatedAt    string `json:"created_at"`
}

type RsvpsList struct {
	Rsvps []RsvpsResource `json:"rsvps"`
}

func (rs RsvpsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)    // GET /rsvps - Read a list of rsvps.
	r.Post("/", rs.Create) // POST /rsvps - Create a new rsvp.

	r.Route("/{id}", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", rs.Get)       // GET /rsvps/{id} - Read a single rsvp by :id.
		r.Put("/", rs.Update)    // PUT /rsvps/{id} - Update a single rsvp by :id.
		r.Delete("/", rs.Delete) // DELETE /rsvps/{id} - Delete a single rsvp by :id.
	})

	return r
}

// Request Handler - GET /rsvps - Read a list of rsvps.
func (rs RsvpsResource) List(w http.ResponseWriter, r *http.Request) {
	list := RsvpsList{}

	rows, err := DB.Query("SELECT * FROM rsvps ORDER BY ID DESC")
	if err != nil {
		return
	}

	for rows.Next() {

		var rsvp RsvpsResource
		err := rows.Scan(&rsvp.Id, &rsvp.EventId, &rsvp.GuestName, &rsvp.Email, &rsvp.NumberGuests, &rsvp.CreatedAt)
		if err != nil {
			panic(err)
		}

		list.Rsvps = append(list.Rsvps, rsvp)
	}
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(string(b)))

}

// Request Handler - POST /rsvps - Create a new rsvp.
func (rs RsvpsResource) Create(w http.ResponseWriter, r *http.Request) {
	var es RsvpsResource
	err := json.NewDecoder(r.Body).Decode(&es)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("RSVP Request Details:\n")
	fmt.Printf("Name: %+s\n", es.GuestName)
	fmt.Printf("ID: %+d\n", es.EventId)
	fmt.Printf("Email: %+s\n", es.Email)
	fmt.Printf("Number of Guests: %+d\n", es.NumberGuests)

	var id int
	var createdAt string

	query := `INSERT INTO rsvps (guest_name, event_id, email, number_guests) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	row := DB.QueryRow(query, es.GuestName, es.EventId, es.Email, es.NumberGuests)
	err = row.Scan(&id, &createdAt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("RSVP Created with ID: %+d\n", id)

	es.Id = id
	es.CreatedAt = createdAt

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(es)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(string(b)))
}

// Request Handler - GET /rsvps/{id} - Read a single rsvp by :id.
func (rs RsvpsResource) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %v\n", r.URL)

	rsvpId := r.Context().Value("id").(string)

	fmt.Printf("rsvp Id: %s\n", rsvpId)

	var es RsvpsResource

	query := `SELECT * FROM rsvps WHERE id = $1`
	row := DB.QueryRow(query, rsvpId)
	err := row.Scan(&es.Id, &es.EventId, &es.GuestName, &es.Email, &es.NumberGuests, &es.CreatedAt)

	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(es)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(string(b)))
}

// Request Handler - PUT /rsvps/{id} - Update a single rsvp by :id.
func (rs RsvpsResource) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %v\n", r.URL)

	rsvpId := r.Context().Value("id").(string)

	fmt.Printf("rsvp Id: %s\n", rsvpId)

	var es RsvpsResource

	err := json.NewDecoder(r.Body).Decode(&es)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("RSVP Request Details:\n")
	fmt.Printf("Guest Name: %+s\n", es.GuestName)
	fmt.Printf("Event ID: %+d\n", es.EventId)
	fmt.Printf("Email: %+s\n", es.Email)
	fmt.Printf("Number of Guests: %+d\n", es.NumberGuests)

	query := `UPDATE rsvps SET guest_name=$1, event_id=$2, email=$3, number_guests=$4 WHERE id=$5 RETURNING id, event_id, guest_name, email, number_guests, created_at;`
	row := DB.QueryRow(query, es.GuestName, es.EventId, es.Email, es.NumberGuests, rsvpId)
	err = row.Scan(&es.Id, &es.EventId, &es.GuestName, &es.Email, &es.NumberGuests, &es.CreatedAt)

	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(es)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(string(b)))
}

// Request Handler - DELETE /rsvps/{id} - Delete a single rsvp by :id.
func (rs RsvpsResource) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %v\n", r.URL)

	rsvpId := r.Context().Value("id").(string)

	fmt.Printf("rsvp Id: %s\n", rsvpId)

	query := `DELETE FROM rsvps WHERE id = $1`
	_, err := DB.Exec(query, rsvpId)

	if err != nil {
		panic(err)
	}

	w.Write([]byte(`{"status":"RSVP Record Deleted"}`))
}
