package models

import (
	"database/sql"
	"fmt"

	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/database"
	"github.com/korentmaj/go-ecdsa-status-netis-challenge/internal/status"
)

func GetStatus(statusId string) (*status.StatusList, error) {
	var encodedList []byte
	err := database.DB.QueryRow("SELECT encoded_list FROM statuses WHERE id = $1", statusId).Scan(&encodedList)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("status not found")
		}
		return nil, fmt.Errorf("failed to query status: %v", err)
	}

	status := status.NewStatusList()
	// Decode the status list from the stored encodedList
	// TODO: Add decoding logic
	return status, nil
}

func SaveStatus(statusId string, status *status.StatusList) error {
	encodedList, err := status.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode status list: %v", err)
	}

	_, err = database.DB.Exec("UPDATE statuses SET encoded_list = $1 WHERE id = $2", encodedList, statusId)
	if err != nil {
		return fmt.Errorf("failed to update status: %v", err)
	}

	return nil
}

func CreateNewStatus(status *status.StatusList) (string, error) {
	encodedList, err := status.Encode()
	if err != nil {
		return "", fmt.Errorf("failed to encode status list: %v", err)
	}

	var statusId string
	err = database.DB.QueryRow("INSERT INTO statuses (encoded_list) VALUES ($1) RETURNING id", encodedList).Scan(&statusId)
	if err != nil {
		return "", fmt.Errorf("failed to insert new status: %v", err)
	}

	return statusId, nil
}

func GetAllStatusIds() ([]string, error) {
	rows, err := database.DB.Query("SELECT id FROM statuses")
	if err != nil {
		return nil, fmt.Errorf("failed to query status ids: %v", err)
	}
	defer rows.Close()

	var statusIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan status id: %v", err)
		}
		statusIds = append(statusIds, id)
	}

	return statusIds, nil
}
