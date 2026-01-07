package domain

type EventPacket struct {
	ID             int
	OwnerID        int
	Name           string
	Location       *string
	Description    *string
	AllocatedSeats *int
}

type EventPacketFilter struct {
	Location    *string
	Name        *string
	Description *string
	MinSeats    *int
	MaxSeats    *int
	Page        *int
	PerPage     *int
	OrderBy     *string
}

func (filter *EventPacketFilter) Default() {
	if filter.Page == nil {
		filter.Page = new(int)
		*filter.Page = 1
	}
	if filter.PerPage == nil {
		filter.PerPage = new(int)
		*filter.PerPage = 10
	}
}

func (filter *EventPacketFilter) Validate() error {
	validOrderings := map[string]bool{
		"name_asc":   true,
		"name_desc":  true,
		"seats_asc":  true,
		"seats_desc": true,
	}

	if filter.OrderBy != nil && !validOrderings[*filter.OrderBy] {
		return &ValidationError{Reason: "invalid order by. valid options: name_asc/desc, seats_asc/desc"}
	}
	return nil
}
