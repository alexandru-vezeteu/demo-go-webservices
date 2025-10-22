package gormrepository

import (
	"errors"
	"eventManager/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"gorm.io/gorm"
)

type GormEventRepository struct {
	DB *gorm.DB
}

func (r *GormEventRepository) Create(event *domain.Event) (*domain.Event, error) {
	gormEvent := gormmodel.FromDomain(event)

	if err := r.DB.Create(gormEvent).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.NewEventAlreadyExistsError(gormEvent.Name)
		}
		return nil, domain.NewInternalError("failed to persist event", err)
	}

	event.ID = gormEvent.ID

	ret := gormEvent.ToDomain()

	return ret, nil
}

func (r *GormEventRepository) GetByID(id int) (*domain.Event, error) {

	var ret gormmodel.GormEvent
	result := r.DB.First(&ret, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, domain.NewEventNotFoundError(id)
	} else if result.Error != nil {
		return nil, domain.NewInternalError("could not find the event", result.Error)
	}

	retDomain := ret.ToDomain()

	return retDomain, nil
}

func (r *GormEventRepository) Update(event *domain.Event) (*domain.Event, error) {
	return nil, nil
}

func (r *GormEventRepository) Delete(event *domain.Event) (*domain.Event, error) {
	return nil, nil

}
