package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Entity struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Email    string `json:"email"`
	Capacity int    `json:"capacity"`
}
type Booking struct {
	ID       string    `json:"id"`
	RoomID   string    `json:"roomId"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
	Size     int       `json:"size"`
}

func reply(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(value)
}
func fail(w http.ResponseWriter, code int, message string) {
	reply(w, code, map[string]string{"message": message})
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	kind := os.Getenv("ENTITY")
	if kind != "students" && kind != "teachers" && kind != "rooms" {
		log.Fatal("Unknown entity")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/"+kind, func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), "SELECT id,name,group_name,email,capacity FROM "+kind+" ORDER BY name")
		if err != nil {
			fail(w, 500, err.Error())
			return
		}
		defer rows.Close()
		out := []Entity{}
		for rows.Next() {
			var e Entity
			if err = rows.Scan(&e.ID, &e.Name, &e.Group, &e.Email, &e.Capacity); err != nil {
				fail(w, 500, err.Error())
				return
			}
			out = append(out, e)
		}
		if rows.Err() != nil {
			fail(w, 500, rows.Err().Error())
			return
		}
		reply(w, 200, out)
	})
	mux.HandleFunc("GET /api/"+kind+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		var e Entity
		err := db.QueryRowContext(r.Context(), "SELECT id,name,group_name,email,capacity FROM "+kind+" WHERE id=$1", r.PathValue("id")).Scan(&e.ID, &e.Name, &e.Group, &e.Email, &e.Capacity)
		if errors.Is(err, sql.ErrNoRows) {
			fail(w, 404, "Not found")
			return
		}
		if err != nil {
			fail(w, 400, "Invalid ID")
			return
		}
		reply(w, 200, e)
	})
	save := func(w http.ResponseWriter, r *http.Request) {
		var e Entity
		if json.NewDecoder(r.Body).Decode(&e) != nil || e.Name == "" || (kind == "rooms" && e.Capacity <= 0) || (kind == "students" && e.Group == "") {
			fail(w, 400, "Invalid entity")
			return
		}
		var err error
		code := 201
		if r.Method == "POST" {
			err = db.QueryRowContext(r.Context(), "INSERT INTO "+kind+" (name,group_name,email,capacity) VALUES ($1,$2,$3,$4) RETURNING id", e.Name, e.Group, e.Email, e.Capacity).Scan(&e.ID)
		} else {
			code = 200
			e.ID = r.PathValue("id")
			var result sql.Result
			result, err = db.ExecContext(r.Context(), "UPDATE "+kind+" SET name=$2,group_name=$3,email=$4,capacity=$5 WHERE id=$1", e.ID, e.Name, e.Group, e.Email, e.Capacity)
			if err == nil {
				n, _ := result.RowsAffected()
				if n == 0 {
					fail(w, 404, "Not found")
					return
				}
			}
		}
		if err != nil {
			fail(w, 400, err.Error())
			return
		}
		reply(w, code, e)
	}
	mux.HandleFunc("POST /api/"+kind, save)
	mux.HandleFunc("PUT /api/"+kind+"/{id}", save)
	mux.HandleFunc("DELETE /api/"+kind+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		result, err := db.ExecContext(r.Context(), "DELETE FROM "+kind+" WHERE id=$1", r.PathValue("id"))
		if err != nil {
			fail(w, 409, "Entity is referenced or ID is invalid")
			return
		}
		n, _ := result.RowsAffected()
		if n == 0 {
			fail(w, 404, "Not found")
			return
		}
		w.WriteHeader(204)
	})
	if kind == "rooms" {
		mux.HandleFunc("POST /api/bookings", func(w http.ResponseWriter, r *http.Request) {
			var b Booking
			if json.NewDecoder(r.Body).Decode(&b) != nil || b.Size <= 0 || !b.EndsAt.After(b.StartsAt) {
				fail(w, 400, "Invalid booking")
				return
			}
			tx, err := db.BeginTx(r.Context(), nil)
			if err != nil {
				fail(w, 500, err.Error())
				return
			}
			defer tx.Rollback()
			var capacity int
			err = tx.QueryRowContext(r.Context(), "SELECT capacity FROM rooms WHERE id=$1 FOR UPDATE", b.RoomID).Scan(&capacity)
			if err != nil {
				fail(w, 404, "Room not found")
				return
			}
			if b.Size > capacity {
				fail(w, 409, "Room capacity exceeded")
				return
			}
			_, err = tx.ExecContext(r.Context(), "INSERT INTO bookings(id,room_id,starts_at,ends_at) VALUES ($1,$2,$3,$4)", b.ID, b.RoomID, b.StartsAt, b.EndsAt)
			if err != nil {
				fail(w, 409, "Room is occupied or booking ID exists")
				return
			}
			if err = tx.Commit(); err != nil {
				fail(w, 500, err.Error())
				return
			}
			log.Printf("Booking created: %s", b.ID)
			reply(w, 201, b)
		})
		mux.HandleFunc("DELETE /api/bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
			_, err := db.ExecContext(r.Context(), "DELETE FROM bookings WHERE id=$1", r.PathValue("id"))
			if err != nil {
				fail(w, 400, "Invalid booking ID")
				return
			}
			log.Printf("Booking cancelled: %s", r.PathValue("id"))
			w.WriteHeader(204)
		})
	}
	port := os.Getenv("PORT")
	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}
	log.Printf("%s service :%s", kind, port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
