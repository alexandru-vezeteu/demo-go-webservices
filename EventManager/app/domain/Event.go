package domain

type Event struct {
	ID          int
	OwnerID     int
	Name        string
	Location    string
	Description string
	Seats       int
}
