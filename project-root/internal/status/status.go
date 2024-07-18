package status

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
)

// StatusList represents a list of boolean statuses stored in a byte slice.
type StatusList struct {
	statuses []byte
}

// NewStatusList creates a new StatusList.
func NewStatusList() *StatusList {
	return &StatusList{
		statuses: []byte{},
	}
}

// SetStatus sets the status at the given index to the specified value (true or false).
func (sl *StatusList) SetStatus(index int, value bool) error {
	byteIndex := index / 8
	bitIndex := index % 8

	if byteIndex >= len(sl.statuses) {
		return fmt.Errorf("index out of range")
	}

	if value {
		sl.statuses[byteIndex] |= (1 << bitIndex)
	} else {
		sl.statuses[byteIndex] &^= (1 << bitIndex)
	}

	return nil
}

// AddStatus adds a new status to the list and returns its index.
func (sl *StatusList) AddStatus(value bool) int {
	index := len(sl.statuses) * 8
	sl.statuses = append(sl.statuses, 0)

	if value {
		sl.SetStatus(index, true)
	}

	return index
}

// Encode encodes the status list into a gzipped base64 encoded string.
func (sl *StatusList) Encode() (string, error) {
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	_, err := gzipWriter.Write(sl.statuses)
	if err != nil {
		return "", fmt.Errorf("failed to write to gzip writer: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %v", err)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
