package gormrepository

import (
	"errors"
	"eventManager/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"gorm.io/gorm"
)

type GormEventPacketInclusionRepository struct {
	DB *gorm.DB
}

func (r *GormEventPacketInclusionRepository) Create(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	gormModel := gormmodel.FromEventPacketInclusion(event)
	if err := r.DB.Create(gormModel).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, errors.New("invalid event_id or packet_id: foreign key constraint failed")
		}
		// Check for duplicate key error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("this event-packet inclusion already exists")
		}
		return nil, err
	}

	// Preload relationships after creation
	if err := r.DB.Preload("Event").Preload("Packet").
		First(gormModel, "event_id = ? AND packet_id = ?", gormModel.EventID, gormModel.PacketID).Error; err != nil {
		return nil, err
	}

	return gormModel.ToDomain(), nil
}

func (r *GormEventPacketInclusionRepository) GetEventsByPacketID(packetID int) ([]*domain.Event, error) {
	var gormInclusions []gormmodel.GormEventPacketInclusion

	if err := r.DB.Preload("Event").
		Where("packet_id = ?", packetID).
		Find(&gormInclusions).Error; err != nil {
		return nil, err
	}

	// Extract just the Event objects
	result := make([]*domain.Event, len(gormInclusions))
	for i, inclusion := range gormInclusions {
		result[i] = inclusion.Event.ToDomain()
	}

	return result, nil
}

func (r *GormEventPacketInclusionRepository) GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error) {
	var gormInclusions []gormmodel.GormEventPacketInclusion

	// Query inclusions and preload Packet objects
	if err := r.DB.Preload("Packet").
		Where("event_id = ?", eventID).
		Find(&gormInclusions).Error; err != nil {
		return nil, err
	}

	// Extract just the Packet objects
	result := make([]*domain.EventPacket, len(gormInclusions))
	for i, inclusion := range gormInclusions {
		result[i] = inclusion.Packet.ToDomain()
	}

	return result, nil
}

func (r *GormEventPacketInclusionRepository) Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	// Update the record
	result := r.DB.Model(&gormmodel.GormEventPacketInclusion{}).
		Where("event_id = ? AND packet_id = ?", eventID, packetID).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	// Check if record exists
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Fetch the updated record with relationships
	var gormInclusion gormmodel.GormEventPacketInclusion
	if err := r.DB.Preload("Event").Preload("Packet").
		Where("event_id = ? AND packet_id = ?", eventID, packetID).
		First(&gormInclusion).Error; err != nil {
		return nil, err
	}

	return gormInclusion.ToDomain(), nil
}

func (r *GormEventPacketInclusionRepository) Delete(eventID, packetID int) (*domain.EventPacketInclusion, error) {
	// Fetch the record before deleting with relationships
	var gormInclusion gormmodel.GormEventPacketInclusion
	if err := r.DB.Preload("Event").Preload("Packet").
		Where("event_id = ? AND packet_id = ?", eventID, packetID).
		First(&gormInclusion).Error; err != nil {
		return nil, err
	}

	// Delete the record
	if err := r.DB.Delete(&gormInclusion).Error; err != nil {
		return nil, err
	}

	return gormInclusion.ToDomain(), nil
}
