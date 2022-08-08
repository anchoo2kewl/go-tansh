package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type EventsResource struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Address     string `json:"address"`
	Id          int    `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	EventTime   string `json:"event_time"`
}

type EventList struct {
	Events []EventsResource `json:"events"`
}

func (es EventsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", es.List)    // GET /events - Read a list of events.
	r.Post("/", es.Create) // POST /events - Create a new event.

	r.Route("/{id}", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", es.Get)       // GET /events/{id} - Read a single event by :id.
		r.Put("/", es.Update)    // PUT /events/{id} - Update a single event by :id.
		r.Delete("/", es.Delete) // DELETE /events/{id} - Delete a single event by :id.
	})

	return r
}

// Request Handler - GET /events - Read a list of events.
func (rs EventsResource) List(w http.ResponseWriter, r *http.Request) {
	list := EventList{}

	rows, err := DB.Query("SELECT * FROM events ORDER BY ID DESC")
	if err != nil {
		return
	}

	for rows.Next() {
		var event EventsResource
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Address, &event.CreatedAt, &event.EventTime)
		if err != nil {
			return
		}
		list.Events = append(list.Events, event)
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(string(b)))

}

// Request Handler - POST /events - Create a new event.
func (rs EventsResource) Create(w http.ResponseWriter, r *http.Request) {
	var es EventsResource
	err := json.NewDecoder(r.Body).Decode(&es)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Event Request Details:\n")
	fmt.Printf("Name: %+s\n", es.Name)
	fmt.Printf("Description: %+s\n", es.Description)
	fmt.Printf("Location: %+s\n", es.Location)
	fmt.Printf("Address: %+s\n", es.Address)
	fmt.Printf("Event Time: %+s\n", es.EventTime)

	eventTimeConverted, _ := time.Parse(time.RFC3339, es.EventTime)

	var id int
	var createdAt string

	query := `INSERT INTO events (name, description, location, address, event_time) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	row := DB.QueryRow(query, es.Name, es.Description, es.Location, es.Address, eventTimeConverted.Format(time.RFC3339))
	err = row.Scan(&id, &createdAt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Event Created with ID: %+d\n", id)

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

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Request Handler - GET /events/{id} - Read a single event by :id.
func (rs EventsResource) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %v\n", r.URL)

	eventId := r.Context().Value("id").(string)

	fmt.Printf("Event Id: %s\n", eventId)

	var es EventsResource

	query := `SELECT * FROM events WHERE id = $1`
	row := DB.QueryRow(query, eventId)
	err := row.Scan(&es.Id, &es.Name, &es.Description, &es.Location, &es.Address, &es.CreatedAt, &es.EventTime)

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

// Request Handler - PUT /events/{id} - Update a single event by :id.
func (rs EventsResource) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %v\n", r.URL)

	eventId := r.Context().Value("id").(string)

	fmt.Printf("Event Id: %s\n", eventId)

	var es EventsResource

	err := json.NewDecoder(r.Body).Decode(&es)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Request Details:\n")
	fmt.Printf("Name: %+s\n", es.Name)
	fmt.Printf("Description: %+s\n", es.Description)
	fmt.Printf("Location: %+s\n", es.Location)
	fmt.Printf("Address: %+s\n", es.Address)

	query := `UPDATE events SET name=$1, description=$2, location=$3, address=$4, event_time=$5 WHERE id=$6 RETURNING id, name, description, location, address, created_at, event_time;`
	row := DB.QueryRow(query, es.Name, es.Description, es.Address, es.Location, es.EventTime, eventId)
	err = row.Scan(&es.Id, &es.Name, &es.Description, &es.Location, &es.Address, &es.CreatedAt, &es.EventTime)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", es.Name)

	b, err := json.Marshal(es)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(string(b)))
}

// Request Handler - DELETE /events/{id} - Delete a single event by :id.
func (rs EventsResource) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL: %v\n", r.URL)

	eventId := r.Context().Value("id").(string)

	fmt.Printf("Event Id: %s\n", eventId)

	query := `DELETE FROM events WHERE id = $1`
	_, err := DB.Exec(query, eventId)

	if err != nil {
		panic(err)
	}

	w.Write([]byte(`{"status":"Event Record Deleted"}`))
}
