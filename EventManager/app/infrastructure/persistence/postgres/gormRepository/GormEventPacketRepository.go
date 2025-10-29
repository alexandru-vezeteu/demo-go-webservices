package gormrepository

import (
	"eventManager/domain"

	"gorm.io/gorm"
)

type GormEventPacketRepository struct {
	DB *gorm.DB
}

func (r *GormEventPacketRepository) Create(event *domain.EventPacket) (*domain.EventPacket, error) {

	return nil, nil
}

func (r *GormEventPacketRepository) GetByID(id int) (*domain.EventPacket, error) {

	return nil, nil
}

func (r *GormEventPacketRepository) Update(event *domain.EventPacket) (*domain.EventPacket, error) {
	return nil, nil
}

func (r *GormEventPacketRepository) Delete(event *domain.EventPacket) (*domain.EventPacket, error) {
	return nil, nil

}
