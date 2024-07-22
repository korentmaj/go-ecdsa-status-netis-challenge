package models

type StatusList struct {
	ID          int
	EncodedList []byte
	CreatedAt   string
}

type Status struct {
	ID     int
	ListID int
	Status bool
}
