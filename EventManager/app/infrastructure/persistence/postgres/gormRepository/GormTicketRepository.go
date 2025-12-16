package gormrepository

import (
	"context"
	"errors"
	"eventManager/application/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormTicketRepository struct {
	DB *gorm.DB
}

func (r *GormTicketRepository) CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	gormTicket := gormmodel.FromTicket(ticket)

	if err := r.DB.WithContext(ctx).Create(gormTicket).Error; err != nil {
		
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &domain.AlreadyExistsError{Name: gormTicket.Code}
		}
		return nil, &domain.InternalError{Msg: "failed to create ticket", Err: err}
	}

	return gormTicket.ToDomain(), nil
}

func (r *GormTicketRepository) GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error) {
	var ret gormmodel.GormTicket
	result := r.DB.WithContext(ctx).Where("code = ?", code).First(&ret)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &domain.NotFoundError{ID: 0} 
	} else if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not find the ticket", Err: result.Error}
	}

	return ret.ToDomain(), nil
}

func (r *GormTicketRepository) UpdateTicket(ctx context.Context, code string, updates map[string]interface{}) (*domain.Ticket, error) {
	result := r.DB.WithContext(ctx).Model(&gormmodel.GormTicket{}).
		Clauses(clause.Returning{}).
		Where("code = ?", code).
		Updates(updates)

	if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not update the ticket", Err: result.Error}
	}

	if result.RowsAffected == 0 {
		return nil, &domain.NotFoundError{ID: 0}
	}

	return r.GetTicketByCode(ctx, code)
}

func (r *GormTicketRepository) DeleteEvent(ctx context.Context, code string) (*domain.Ticket, error) {
	var ret gormmodel.GormTicket
	var retDomain *domain.Ticket

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("code = ?", code).First(&ret)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &domain.NotFoundError{ID: 0}
		} else if result.Error != nil {
			return &domain.InternalError{Msg: "could not find the ticket", Err: result.Error}
		}

		deleteResult := tx.Where("code = ?", code).Delete(&gormmodel.GormTicket{})

		if deleteResult.Error != nil {
			return &domain.InternalError{Msg: "could not delete the ticket", Err: deleteResult.Error}
		}

		retDomain = ret.ToDomain()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return retDomain, nil
}
