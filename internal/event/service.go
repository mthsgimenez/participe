package event

import (
	"fmt"

	"github.com/mthsgimenez/participe/internal/user"
)

type Repository interface {
	FindById(id int) (*Event, error)
	FindAll() (*[]Event, error)
	Insert(e *Event) (*Event, error)
	Update(e *Event) (*Event, error)
	DeleteById(id int) error
	Exists(id int) (bool, error)
	FindUpcoming() (*[]Event, error)
	CheckinUser(e *Event, u *user.User) error
	FindCheckedUsers(e *Event) (*[]user.User, error)
}

type Service struct {
	eventRepo Repository
}

func NewService(eventRepo Repository) *Service {
	return &Service{eventRepo}
}

func (s *Service) GetEvent(id int) (*Event, error) {
	e, err := s.eventRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("event_service: get event by id: %w", err)
	}

	return e, nil
}

func (s *Service) GetEvents() (*[]Event, error) {
	eList, err := s.eventRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("event_service: get events: %w", err)
	}

	return eList, nil
}

func (s *Service) GetUpcomingEvents() (*[]Event, error) {
	eList, err := s.eventRepo.FindUpcoming()
	if err != nil {
		return nil, fmt.Errorf("event_service: get upcoming events: %w", err)
	}

	return eList, nil
}

func (s *Service) CreateEvent(e *Event) (*Event, error) {
	newEvent, err := s.eventRepo.Insert(e)
	if err != nil {
		return nil, fmt.Errorf("event_service: create event: %w", err)
	}

	return newEvent, nil
}

func (s *Service) UpdateEvent(id int, newData *Event) (*Event, error) {
	event, err := s.eventRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("event_service: update event: %w", err)
	}

	event.Description = newData.Description
	event.Name = newData.Name
	event.Date = newData.Date

	updatedEvent, err := s.eventRepo.Update(event)
	if err != nil {
		return nil, fmt.Errorf("event_service: update event: %w", err)
	}

	return updatedEvent, nil
}

func (s *Service) DeleteEvent(id int) error {
	exists, err := s.eventRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("event_service: delete event: %w", err)
	}

	if !exists {
		return fmt.Errorf("event_service: delete event: %w", ErrEventNotFound)
	}

	if err := s.eventRepo.DeleteById(id); err != nil {
		return fmt.Errorf("event_service: delete event: %w", err)
	}

	return nil
}

func (s *Service) CheckinUserInEvent(e *Event, u *user.User) error {
	if err := s.eventRepo.CheckinUser(e, u); err != nil {
		return fmt.Errorf("event_service: %w", err)
	}

	return nil
}

func (s *Service) GetCheckedUsers(e *Event) (*[]user.User, error) {
	uList, err := s.eventRepo.FindCheckedUsers(e)
	if err != nil {
		return nil, fmt.Errorf("event_service: get checked users: %w", err)
	}

	return uList, nil
}
