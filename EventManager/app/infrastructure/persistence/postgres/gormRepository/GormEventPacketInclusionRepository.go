package gormrepository

import (
	"context"
	"errors"
	"eventManager/application/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"gorm.io/gorm"
)

type GormEventPacketInclusionRepository struct {
	DB *gorm.DB
}

func (r *GormEventPacketInclusionRepository) Create(ctx context.Context, event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	gormModel := gormmodel.FromEventPacketInclusion(event)
	if err := r.DB.WithContext(ctx).Create(gormModel).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, &domain.ForeignKeyError{}
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &domain.AlreadyExistsError{Name: "event-packet inclusion"}
		}
		return nil, &domain.InternalError{Msg: "failed to create event-packet inclusion", Err: err}
	}

	if err := r.DB.WithContext(ctx).Preload("Event").Preload("Packet").
		First(gormModel, "event_id = ? AND packet_id = ?", gormModel.EventID, gormModel.PacketID).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to load created event-packet inclusion", Err: err}
	}

	return gormModel.ToDomain(), nil
}

func (r *GormEventPacketInclusionRepository) GetEventsByPacketID(ctx context.Context, packetID int) ([]*domain.Event, error) {
	var gormInclusions []gormmodel.GormEventPacketInclusion

	if err := r.DB.WithContext(ctx).Preload("Event").
		Where("packet_id = ?", packetID).
		Find(&gormInclusions).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to get events by packet ID", Err: err}
	}

	result := make([]*domain.Event, len(gormInclusions))
	for i, inclusion := range gormInclusions {
		result[i] = inclusion.Event.ToDomain()
	}

	return result, nil
}

func (r *GormEventPacketInclusionRepository) GetEventPacketsByEventID(ctx context.Context, eventID int) ([]*domain.EventPacket, error) {
	var gormInclusions []gormmodel.GormEventPacketInclusion

	if err := r.DB.WithContext(ctx).Preload("Packet").
		Where("event_id = ?", eventID).
		Find(&gormInclusions).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to get event packets by event ID", Err: err}
	}

	result := make([]*domain.EventPacket, len(gormInclusions))
	for i, inclusion := range gormInclusions {
		result[i] = inclusion.Packet.ToDomain()
	}

	return result, nil
}

func (r *GormEventPacketInclusionRepository) Update(ctx context.Context, eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	result := r.DB.WithContext(ctx).Model(&gormmodel.GormEventPacketInclusion{}).
		Where("event_id = ? AND packet_id = ?", eventID, packetID).
		Updates(updates)

	if result.Error != nil {
		return nil, &domain.InternalError{Msg: "failed to update event-packet inclusion", Err: result.Error}
	}

	if result.RowsAffected == 0 {
		return nil, &domain.NotFoundError{ID: eventID}
	}

	var gormInclusion gormmodel.GormEventPacketInclusion
	if err := r.DB.WithContext(ctx).Preload("Event").Preload("Packet").
		Where("event_id = ? AND packet_id = ?", eventID, packetID).
		First(&gormInclusion).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to load updated event-packet inclusion", Err: err}
	}

	return gormInclusion.ToDomain(), nil
}

func (r *GormEventPacketInclusionRepository) Delete(ctx context.Context, eventID, packetID int) (*domain.EventPacketInclusion, error) {
	var gormInclusion gormmodel.GormEventPacketInclusion
	if err := r.DB.WithContext(ctx).Preload("Event").Preload("Packet").
		Where("event_id = ? AND packet_id = ?", eventID, packetID).
		First(&gormInclusion).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{ID: eventID}
		}
		return nil, &domain.InternalError{Msg: "failed to find event-packet inclusion", Err: err}
	}

	if err := r.DB.WithContext(ctx).Delete(&gormInclusion).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to delete event-packet inclusion", Err: err}
	}

	return gormInclusion.ToDomain(), nil
}
