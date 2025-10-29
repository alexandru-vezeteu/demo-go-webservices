package gormrepository

import (
	"eventManager/domain"

	"gorm.io/gorm"
)

type GormEventPacketInclusionRepository struct {
	DB *gorm.DB
}

func (r *GormEventPacketInclusionRepository) Create(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {

	return nil, nil
}

func (r *GormEventPacketInclusionRepository) GetEventsInPacketbyID(id int) (*domain.EventPacketInclusion, error) {
	return nil, nil
}
func (r *GormEventPacketInclusionRepository) GetEventPacketsByEventID(id int) (*domain.EventPacketInclusion, error) {
	return nil, nil
}

func (r *GormEventPacketInclusionRepository) Delete(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	return nil, nil
}

//Update(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
