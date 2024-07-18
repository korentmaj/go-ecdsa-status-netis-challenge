package api

import (
	"project-root/internal/api/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/status/{statusId}", handlers.GetStatus).Methods("GET")
	r.HandleFunc("/api/status/{statusId}/{index}", handlers.SetStatus).Methods("PUT")
	r.HandleFunc("/api/status/{statusId}/{index}", handlers.DeleteStatus).Methods("DELETE")
	r.HandleFunc("/api/status/{statusId}", handlers.CreateStatus).Methods("POST")
	r.HandleFunc("/api/status", handlers.GetAllStatuses).Methods("GET")
	r.HandleFunc("/api/status", handlers.CreateNewStructure).Methods("POST")

	// Applying BasicAuth middleware to all PUT, POST, DELETE methods
	r.HandleFunc("/api/status/{statusId}/{index}", handlers.SetStatus).Methods("PUT").Use(BasicAuth)
	r.HandleFunc("/api/status/{statusId}/{index}", handlers.DeleteStatus).Methods("DELETE").Use(BasicAuth)
	r.HandleFunc("/api/status/{statusId}", handlers.CreateStatus).Methods("POST").Use(BasicAuth)
	r.HandleFunc("/api/status", handlers.CreateNewStructure).Methods("POST").Use(BasicAuth)

	return r
}
