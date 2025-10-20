package Model

type Event struct {
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID       int     `gorm:"column:id_owner;not null" json:"owner_id"`
	Name          string  `gorm:"unique;not null;size:255" json:"name"`
	Location      *string `gorm:"size:255" json:"location,omitempty"`
	Description   *string `gorm:"size:255" json:"description,omitempty"`
	NumberOfSeats *int    `gorm:"column:numberofseats" json:"number_of_seats,omitempty"`
}

func (Event) TableName() string {
	return "event"
}

type EventPacket struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID     int     `gorm:"column:id_owner;not null" json:"owner_id"`
	Name        string  `gorm:"unique;not null;size:255" json:"name"`
	Location    *string `gorm:"size:255" json:"location,omitempty"`
	Description *string `gorm:"size:255" json:"description,omitempty"`
}

func (EventPacket) TableName() string {
	return "event_packets"
}

type PacketEventRelation struct {
	CODE     string `gorm:"primaryKey;size:64" json:"code"`
	PacketID int    `gorm:"not null" json:"packet_id"`
	EventID  int    `gorm:"not null" json:"event_id"`

	Packet EventPacket `gorm:"foreignKey:PacketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"packet,omitempty"`
	Event  Event       `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"event,omitempty"`
}

func (PacketEventRelation) TableName() string {
	return "packet_event_relation"
}

type PacketEventSeats struct {
	PacketID      int  `gorm:"primaryKey" json:"packet_id"`
	EventID       int  `gorm:"primaryKey" json:"event_id"`
	NumberOfSeats *int `json:"number_of_seats,omitempty"`

	Packet EventPacket `gorm:"foreignKey:PacketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"packet,omitempty"`
	Event  Event       `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"event,omitempty"`
}

func (PacketEventSeats) TableName() string {
	return "packet_event_seats"
}
