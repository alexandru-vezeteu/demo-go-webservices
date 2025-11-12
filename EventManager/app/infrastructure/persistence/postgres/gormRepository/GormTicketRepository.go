package gormrepository

import (
	"errors"
	"eventManager/application/domain"
	gormmodel "eventManager/infrastructure/persistence/postgres/gormModel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormTicketRepository struct {
	DB *gorm.DB
}

func (r *GormTicketRepository) CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error) {
	gormTicket := gormmodel.FromTicket(ticket)

	if err := r.DB.Create(gormTicket).Error; err != nil {
		// Check for duplicate key (primary key violation on code)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &domain.AlreadyExistsError{Name: gormTicket.Code}
		}
		return nil, &domain.InternalError{Msg: "failed to create ticket", Err: err}
	}

	return gormTicket.ToDomain(), nil
}

func (r *GormTicketRepository) GetTicketByCode(code string) (*domain.Ticket, error) {
	var ret gormmodel.GormTicket
	result := r.DB.Where("code = ?", code).First(&ret)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &domain.NotFoundError{ID: 0} // Using 0 since code is string
	} else if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not find the ticket", Err: result.Error}
	}

	return ret.ToDomain(), nil
}

func (r *GormTicketRepository) UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error) {
	result := r.DB.Model(&gormmodel.GormTicket{}).
		Clauses(clause.Returning{}).
		Where("code = ?", code).
		Updates(updates)

	if result.Error != nil {
		return nil, &domain.InternalError{Msg: "could not update the ticket", Err: result.Error}
	}

	if result.RowsAffected == 0 {
		return nil, &domain.NotFoundError{ID: 0}
	}

	return r.GetTicketByCode(code)
}

func (r *GormTicketRepository) DeleteEvent(code string) (*domain.Ticket, error) {
	var ret gormmodel.GormTicket
	var retDomain *domain.Ticket

	err := r.DB.Transaction(func(tx *gorm.DB) error {
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
