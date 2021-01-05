package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Structs for a basic to do note

type Note struct {
	User        string `json:"user"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Creates a variable for notes to be a slice of the Note Struct
var notes = []Note{}

func main() {

	// Sets a varibale router, that creates a new router
	router := mux.NewRouter()

	router.HandleFunc("/notes", addNote).Methods("POST")

	router.HandleFunc("/notes", getAllNotes).Methods("GET")

	router.HandleFunc("/notes/{id}", getNote).Methods("GET")

	//  Starts the server at port 8080
	fmt.Println("Now listening on ")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Posts a note to
func addNote(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	json.NewDecoder(r.Body).Decode(&newNote)

	notes = append(notes, newNote)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes)
}

// Gets all data
func getAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// Gets a single notes
func getNote(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("ID could not be converted into an integer"))
	}

	if id >= len(notes) {
		w.WriteHeader(400)
		w.Write([]byte("404: No Notes found"))
	}

	notes := notes[id]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
