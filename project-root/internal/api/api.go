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
	return r
}
