package gormrepository

import (
	"eventManager/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"gorm.io/gorm"
)

type GormEventRepository struct {
	DB *gorm.DB
}

func (r *GormEventRepository) Create(event *domain.Event) (*domain.Event, error) {
	gormEvent, err := gormmodel.FromDomain(event)
	if err != nil {
		return nil, err
	}
	gormEvent.ID = 0
	if err = r.DB.Create(gormEvent).Error; err != nil {
		return nil, err
	}

	event.ID = gormEvent.ID
	ret, err := gormEvent.ToDomain()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *GormEventRepository) Update(event *domain.Event) (*domain.Event, error) {
	return nil, nil
}

func (r *GormEventRepository) GetByID(id int) (*domain.Event, error) {
	return nil, nil
}

func (r *GormEventRepository) Delete(event *domain.Event) (*domain.Event, error) {
	return nil, nil

}
