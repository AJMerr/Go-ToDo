package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Structs for a basic to do note

type Note struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User        string             `json:"user,omitempty" bson:"user,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
}

// Creates a variable for notes to be a slice of the Note Struct
var notes = []Note{}

func main() {

	// Sets a varibale router, that creates a new router
	clientOptions := options.Client().ApplyURI("mongodb://localhost/test")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()

	router.HandleFunc("/notes", addNote).Methods("POST")

	router.HandleFunc("/notes", getAllNotes).Methods("GET")

	router.HandleFunc("/notes/{id}", getNote).Methods("GET")

	router.HandleFunc("/notes/{id}", updateNote).Methods("PUT")

	router.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")

	//  Starts the server at port 8080
	fmt.Println("Now listening on PORT: 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// Posts a note
func addNote(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	_ = json.NewDecoder(r.Body).Decode(&newNote)
	collection := client.Database("GoToDo").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, newNote)

	// Reads the JSON request body
	json.NewDecoder(r.Body).Decode(&newNote)

	// Appends the Notes Struct
	notes = append(notes, newNote)

	// Sets the HTTP header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Sends a JSON response
	json.NewEncoder(w).Encode(result)
}

// Gets all data
func getAllNotes(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	collection := client.Database("GoToDo").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var note Note
		cursor.Decode(&note)
		notes = append(notes, note)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// Gets a single notes
func getNote(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted into an integer"))
	}

	if id >= len(notes) {
		w.WriteHeader(404)
		w.Write([]byte("404: No Notes found"))
	}

	notes := notes[id]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// Edit a note
func updateNote(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted into an integer"))
	}

	if id >= len(notes) {
		w.WriteHeader(404)
		w.Write([]byte("404 Item could not be updated"))
	}

	var updateNote Note
	json.NewDecoder(r.Body).Decode(&updateNote)

	notes[id] = updateNote

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(updateNote)
}

// Deletes a note
func deleteNote(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error 400: ID cannot be  converted into an integer"))
	}

	if id >= len(notes) {
		w.WriteHeader(494)
		w.Write([]byte("Error 404: ID not found"))
	}

	notes = append(notes[:id], notes[id+1:]...)

	w.WriteHeader(200)
}
