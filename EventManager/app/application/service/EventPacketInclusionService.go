package service

import (
	"eventManager/application/repository"
	"eventManager/domain"
)

type eventPacketInclusionService struct {
	repo repository.EventPacketInclusionRepository
}

func NewEventPacketInclusionService(repo repository.EventPacketInclusionRepository) *eventPacketInclusionService {
	return &eventPacketInclusionService{repo: repo}
}

func (service *eventPacketInclusionService) CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	return nil, nil

}

func (service *eventPacketInclusionService) GetEventsInPacketbyID(id int) (*domain.EventPacketInclusion, error) {
	return nil, nil
}
func (service *eventPacketInclusionService) GetEventPacketsByEventID(id int) (*domain.EventPacketInclusion, error) {
	return nil, nil
}

func (service *eventPacketInclusionService) DeleteEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	return nil, nil
}

//Update(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
