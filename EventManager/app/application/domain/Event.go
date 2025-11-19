package domain

type Event struct {
	ID          int
	OwnerID     int
	Name        string
	Location    *string
	Description *string
	Seats       *int
}

// GetState returns the computed state for HATEOAS link generation
// For now, all events are "active" - customize this based on your business logic
func (e *Event) GetState() string {
	return "active"
}

type EventFilter struct {
	Location    *string
	Name        *string
	Description *string
	MinSeats    *int
	MaxSeats    *int
	Page        *int
	PerPage     *int
	OrderBy     *string
}

func (filter *EventFilter) Default() {
	if filter.Page == nil {
		filter.Page = new(int)
		*filter.Page = 1
	}
	if filter.PerPage == nil {
		filter.PerPage = new(int)
		*filter.PerPage = 10
	}
}

func (filter *EventFilter) Validate() error {
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
