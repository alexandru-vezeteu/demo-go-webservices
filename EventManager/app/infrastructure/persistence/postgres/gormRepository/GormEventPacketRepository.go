package gormrepository

import (
	"context"
	"errors"
	"eventManager/application/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormEventPacketRepository struct {
	DB *gorm.DB
}

func (r *GormEventPacketRepository) Create(ctx context.Context, event *domain.EventPacket) (*domain.EventPacket, error) {

	gormEvent := gormmodel.FromEventPacket(event)

	if err := r.DB.WithContext(ctx).Create(gormEvent).Error; err != nil {
		
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "23505") {
			return nil, &domain.AlreadyExistsError{Name: gormEvent.Name}

		}
		return nil, &domain.InternalError{Msg: "failed to persist event", Err: err}
	}

	return gormEvent.ToDomain(), nil
}

func (r *GormEventPacketRepository) GetByID(ctx context.Context, id int) (*domain.EventPacket, error) {

	var ret gormmodel.GormEventPacket
	result := r.DB.WithContext(ctx).Where("id = ?", id).First(&ret)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &domain.NotFoundError{ID: id}
	} else if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not find the event", Err: result.Error}
	}

	retDomain := ret.ToDomain()

	return retDomain, nil
}

func (r *GormEventPacketRepository) Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error) {
	result := r.DB.WithContext(ctx).Model(&gormmodel.GormEventPacket{}).Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {

		if strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "23505") {

			return nil, &domain.UniqueNameError{Msg: updates["name"].(string)}
		}
		return nil, &domain.InternalError{Msg: "could not update the event", Err: result.Error}
	}

	if result.RowsAffected == 0 {
		return nil, &domain.NotFoundError{ID: id}
	}

	return r.GetByID(ctx, id)
}

func (r *GormEventPacketRepository) Delete(ctx context.Context, id int) (*domain.EventPacket, error) {
	var ret gormmodel.GormEventPacket
	var retDomain *domain.EventPacket

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		result := tx.Where("id = ?", id).First(&ret)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &domain.NotFoundError{ID: id}
		} else if result.Error != nil {
			return &domain.InternalError{Msg: "could not find the event", Err: result.Error}
		}

		deleteResult := tx.Where("id = ?", id).Delete(&gormmodel.GormEventPacket{})

		if deleteResult.Error != nil {
			return &domain.InternalError{Msg: "could not delete the event", Err: deleteResult.Error}
		}

		retDomain = ret.ToDomain()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return retDomain, nil

}

func (r *GormEventPacketRepository) CountSoldTickets(ctx context.Context, packetID int) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&gormmodel.GormTicket{}).Where("packet_id = ?", packetID).Count(&count).Error
	if err != nil {
		return 0, &domain.InternalError{Msg: "failed to count tickets", Err: err}
	}
	return int(count), nil
}
