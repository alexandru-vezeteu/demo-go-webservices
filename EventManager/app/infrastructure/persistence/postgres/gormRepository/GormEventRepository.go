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

type GormEventRepository struct {
	DB *gorm.DB
}

func (r *GormEventRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	gormEvent := gormmodel.FromEvent(event)

	if err := r.DB.WithContext(ctx).Create(gormEvent).Error; err != nil {
		
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "23505") {
			return nil, &domain.AlreadyExistsError{Name: gormEvent.Name}

		}
		return nil, &domain.InternalError{Msg: "failed to persist event", Err: err}
	}

	return gormEvent.ToDomain(), nil
}

func (r *GormEventRepository) GetByID(ctx context.Context, id int) (*domain.Event, error) {

	var ret gormmodel.GormEvent
	result := r.DB.WithContext(ctx).Where("id = ?", id).First(&ret)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.NotFoundError{ID: id}
		}
		return nil, &domain.InternalError{Msg: "could not find the event", Err: result.Error}
	}

	retDomain := ret.ToDomain()

	return retDomain, nil
}

func (r *GormEventRepository) Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.Event, error) {

	result := r.DB.WithContext(ctx).Model(&gormmodel.GormEvent{}).Clauses(clause.Returning{}).
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

func (r *GormEventRepository) Delete(ctx context.Context, id int) (*domain.Event, error) {
	var ret gormmodel.GormEvent
	var retDomain *domain.Event

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

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

func (r *GormEventRepository) FilterEvents(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error) {
	if filter == nil {
		return nil, &domain.ValidationError{Reason: "filter cannot be nil"}
	}
	if filter.Page == nil {
		return nil, &domain.ValidationError{Field: "page", Reason: "page cannot be nil"}
	}
	if filter.PerPage == nil {
		return nil, &domain.ValidationError{Field: "per_page", Reason: "per_page cannot be nil"}
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	var gormEvents []gormmodel.GormEvent
	query := r.DB.WithContext(ctx).Model(&gormmodel.GormEvent{})

	if filter.Name != nil {
		query = query.Where("name LIKE ?", "%"+*filter.Name+"%")
	}

	if filter.Location != nil {
		query = query.Where("location LIKE ?", "%"+*filter.Location+"%")
	}

	if filter.Description != nil {
		query = query.Where("description LIKE ?", "%"+*filter.Description+"%")
	}

	if filter.MinSeats != nil {
		query = query.Where("seats >= ?", *filter.MinSeats)
	}

	if filter.MaxSeats != nil {
		query = query.Where("seats <= ?", *filter.MaxSeats)
	}

	var sortTranslationMap = map[string]string{
		"name_asc":   "name asc",
		"name_desc":  "name desc",
		"seats_asc":  "seats asc",
		"seats_desc": "seats desc",
	}

	if filter.OrderBy != nil {
		if sqlSortString, ok := sortTranslationMap[*filter.OrderBy]; ok {
			query = query.Order(sqlSortString)
		}
	} else {
		query = query.Order("id asc")
	}

	limit := *filter.PerPage
	page := *filter.Page
	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&gormEvents).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to filter events", Err: err}
	}

	domainEvents := make([]*domain.Event, 0, len(gormEvents))
	for _, gormEvent := range gormEvents {
		domainEvents = append(domainEvents, gormEvent.ToDomain())
	}

	return domainEvents, nil

}

func (r *GormEventRepository) CountSoldTickets(ctx context.Context, eventID int) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&gormmodel.GormTicket{}).Where("event_id = ?", eventID).Count(&count).Error
	if err != nil {
		return 0, &domain.InternalError{Msg: "failed to count tickets", Err: err}
	}
	return int(count), nil
}
