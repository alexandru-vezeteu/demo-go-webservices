package domain

type EventPacket struct {
	ID             int
	OwnerID        int
	Name           string
	Location       *string
	Description    *string
	AllocatedSeats *int
}
