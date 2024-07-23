package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/status"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/pkg/models"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusId := vars["statusId"]
	indexStr := r.URL.Query().Get("index")
	log.Printf("Received request for statusId: %s, index: %s", statusId, indexStr)

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	status, err := models.GetStatus(statusId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Status not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to query status", http.StatusInternalServerError)
		return
	}

	encodedList, err := status.Encode()
	if err != nil {
		http.Error(w, "Failed to encode status list", http.StatusInternalServerError)
		return
	}

	iat := time.Now().Unix()
	exp := time.Now().Add(24 * time.Hour).Unix()
	domain := "localhost" // vstavi lastno domeno

	payload := map[string]interface{}{
		"iat": iat,
		"exp": exp,
		"iss": fmt.Sprintf("http://%s/api/status/%s", domain, statusId),
		"status": map[string]interface{}{
			"encodedList": encodedList,
			"index":       index,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func SetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusId := vars["statusId"]
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	status, err := models.GetStatus(statusId)
	if err != nil {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	if err := status.SetStatus(index, true); err != nil {
		http.Error(w, "Failed to set status", http.StatusInternalServerError)
		return
	}

	if err := models.SaveStatus(statusId, status); err != nil {
		http.Error(w, "Failed to save status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusId := vars["statusId"]
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	status, err := models.GetStatus(statusId)
	if err != nil {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	if err := status.SetStatus(index, false); err != nil {
		http.Error(w, "Failed to set status", http.StatusInternalServerError)
		return
	}

	if err := models.SaveStatus(statusId, status); err != nil {
		http.Error(w, "Failed to save status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusId := vars["statusId"]

	status, err := models.GetStatus(statusId)
	if err != nil {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	index := status.AddStatus(false)

	if err := models.SaveStatus(statusId, status); err != nil {
		http.Error(w, "Failed to save status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"index": index})
}

func GetAllStatuses(w http.ResponseWriter, r *http.Request) {
	statusIds, err := models.GetAllStatusIds()
	if err != nil {
		http.Error(w, "Failed to get status ids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusIds)
}

func CreateNewStructure(w http.ResponseWriter, r *http.Request) {
	status := status.NewStatusList()
	statusId, err := models.CreateNewStatus(status)
	if err != nil {
		http.Error(w, "Failed to create new status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"statusId": statusId})
}
