package gormrepository

import (
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

func (r *GormEventPacketRepository) Create(event *domain.EventPacket) (*domain.EventPacket, error) {

	gormEvent := gormmodel.FromEventPacket(event)

	if err := r.DB.Create(gormEvent).Error; err != nil {
		//postgres err code
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "23505") {
			return nil, &domain.AlreadyExistsError{Name: gormEvent.Name}

		}
		return nil, &domain.InternalError{Msg: "failed to persist event", Err: err}
	}

	return gormEvent.ToDomain(), nil
}

func (r *GormEventPacketRepository) GetByID(id int) (*domain.EventPacket, error) {

	var ret gormmodel.GormEventPacket
	result := r.DB.Where("id = ?", id).First(&ret)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &domain.NotFoundError{ID: id}
	} else if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not find the event", Err: result.Error}
	}

	retDomain := ret.ToDomain()

	return retDomain, nil
}

func (r *GormEventPacketRepository) Update(id int, updates map[string]interface{}) (*domain.EventPacket, error) {
	result := r.DB.Model(&gormmodel.GormEvent{}).Clauses(clause.Returning{}).
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

	return r.GetByID(id)
}

func (r *GormEventPacketRepository) Delete(id int) (*domain.EventPacket, error) {
	var ret gormmodel.GormEventPacket
	var retDomain *domain.EventPacket

	err := r.DB.Transaction(func(tx *gorm.DB) error {

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
