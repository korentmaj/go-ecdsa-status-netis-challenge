package api

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/status/{statusId}", GetStatus).Methods("GET")
	r.HandleFunc("/api/status/{statusId}/{index}", SetStatus).Methods("PUT")
	r.HandleFunc("/api/status/{statusId}/{index}", DeleteStatus).Methods("DELETE")
	r.HandleFunc("/api/status/{statusId}", CreateStatus).Methods("POST")
	r.HandleFunc("/api/status", GetAllStatuses).Methods("GET")
	r.HandleFunc("/api/status", CreateNewStructure).Methods("POST")

	// Applying BasicAuth middleware to all PUT, POST, DELETE methods
	r.Use(BasicAuth)

	return r
}
