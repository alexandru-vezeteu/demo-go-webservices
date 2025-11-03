package gormrepository

import (
	"errors"
	"eventManager/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormEventRepository struct {
	DB *gorm.DB
}

func (r *GormEventRepository) Create(event *domain.Event) (*domain.Event, error) {
	gormEvent := gormmodel.FromEvent(event)

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

func (r *GormEventRepository) GetByID(id int) (*domain.Event, error) {

	var ret gormmodel.GormEvent
	result := r.DB.Where("id = ?", id).First(&ret)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &domain.NotFoundError{ID: id}
	} else if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not find the event", Err: result.Error}
	}

	retDomain := ret.ToDomain()

	return retDomain, nil
}

func (r *GormEventRepository) Update(id int, updates map[string]interface{}) (*domain.Event, error) {

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

func (r *GormEventRepository) Delete(id int) (*domain.Event, error) {
	var ret gormmodel.GormEvent
	var retDomain *domain.Event

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		result := tx.Where("id = ?", id).First(&ret)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &domain.NotFoundError{ID: id}
		} else if result.Error != nil {
			return &domain.InternalError{Msg: "could not find the event", Err: result.Error}
		}

		deleteResult := tx.Where("id = ?", id).Delete(&gormmodel.GormEvent{})

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
