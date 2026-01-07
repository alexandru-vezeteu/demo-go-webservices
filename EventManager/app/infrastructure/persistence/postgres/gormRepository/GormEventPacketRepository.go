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

func (r *GormEventPacketRepository) FilterEventPackets(ctx context.Context, filter *domain.EventPacketFilter) ([]*domain.EventPacket, error) {
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

	var gormPackets []gormmodel.GormEventPacket
	query := r.DB.WithContext(ctx).Model(&gormmodel.GormEventPacket{})

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
		query = query.Where("allocated_seats >= ?", *filter.MinSeats)
	}

	if filter.MaxSeats != nil {
		query = query.Where("allocated_seats <= ?", *filter.MaxSeats)
	}

	var sortTranslationMap = map[string]string{
		"name_asc":   "name asc",
		"name_desc":  "name desc",
		"seats_asc":  "allocated_seats asc",
		"seats_desc": "allocated_seats desc",
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

	if err := query.Find(&gormPackets).Error; err != nil {
		return nil, &domain.InternalError{Msg: "failed to filter event packets", Err: err}
	}

	domainPackets := make([]*domain.EventPacket, 0, len(gormPackets))
	for _, gormPacket := range gormPackets {
		domainPackets = append(domainPackets, gormPacket.ToDomain())
	}

	return domainPackets, nil
}
