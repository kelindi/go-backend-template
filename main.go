package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(uuid.New())

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/Filess", createFile).Methods("POST")
	r.HandleFunc("/api/v1/files", getFiles).Methods("GET")
	r.HandleFunc("/api/v1/files/{id}", getFile).Methods("GET")
	r.HandleFunc("/api/v1/files/{id}", updateFile).Methods("PATCH")
	r.HandleFunc("/api/v1/files/{id}", deleteFile).Methods("DELETE")

	// db := connect()
	// defer db.Close()

	// file := File{MCID: uuid.New()}
	// fmt.Println(file)

	// Starting Server
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Create File
func createFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Creating File Instance
	file := &File{
		MCID: uuid.New().String(),
	}

	// Decoding Request
	_ = json.NewDecoder(r.Body).Decode(&file)

	// Inserting Into Database
	_, err := db.Model(file).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning File
	json.NewEncoder(w).Encode(file)
}

// Get Files
func getFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Creating Files Slice
	var files []File
	if err := db.Model(&files).Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning Files
	json.NewEncoder(w).Encode(files)
}

// Get File
func getFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Get MCID
	params := mux.Vars(r)
	fileMCID := params["mcid"]

	// Creating file Instance
	file := &File{MCID: fileMCID}
	if err := db.Model(file).WherePK().Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning file
	json.NewEncoder(w).Encode(file)
}

// Update file
func updateFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Get MCID
	params := mux.Vars(r)
	fileMCID := params["mcid"]

	// Creating file Instance
	file := &File{MCID: fileMCID}

	_ = json.NewDecoder(r.Body).Decode(&file)

	
	_, err := db.Model(file).WherePK().Set("MCID = ?, CID = ?, Name = ?, Collection = ?", file.MCID, file.CID, file.Name, file.Collection).Update()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning File
	json.NewEncoder(w).Encode(file)
}

// Delete File
func deleteFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get Connect
	db := connect()
	defer db.Close()

	// Get ID
	// Get MCID
	params := mux.Vars(r)
	fileMCID := params["mcid"]

	// Creating File Instance Alternative Way
	// file := &File{MCID: fileMCID}
	// result, err := db.Model(file).WherePK().Delete()

	// Creating File Instance
	file := &File{}
	result, err := db.Model(file).Where("id = ?", fileMCID).Delete()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning result
	json.NewEncoder(w).Encode(result)
}