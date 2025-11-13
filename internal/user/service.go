package user

import (
	"fmt"

	"github.com/mthsgimenez/participe/internal/company"
)

type Repository interface {
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() (*[]User, error)
	Insert(u *User) (*User, error)
	Update(u *User) (*User, error)
	DeleteById(id int) error
	Exists(id int) (bool, error)
}

type Service struct {
	userRepo    Repository
	companyRepo company.Repository
}

func NewService(userRepo Repository, companyRepo company.Repository) *Service {
	return &Service{userRepo, companyRepo}
}

func (s *Service) GetUser(id int) (*User, error) {
	u, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("user_service: get user by id: %w", err)
	}

	comp, err := s.companyRepo.FindById(u.Company.Id)
	if err != nil {
		return nil, fmt.Errorf("user_service: get user by id (load company): %w", err)
	}

	u.Company = *comp
	return u, nil
}

func (s *Service) GetUserByEmail(email string) (*User, error) {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user_service: get user by email: %w", err)
	}

	comp, err := s.companyRepo.FindById(u.Company.Id)
	if err != nil {
		return nil, fmt.Errorf("user_service: get user by email (load company): %w", err)
	}

	u.Company = *comp
	return u, nil
}

func (s *Service) GetUsers() (*[]User, error) {
	uList, err := s.userRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("user_service: get users: %w", err)
	}

	for i, u := range *uList {
		comp, err := s.companyRepo.FindById(u.Company.Id)
		if err != nil {
			return nil, fmt.Errorf("user_service: get users (load company): %w", err)
		}
		(*uList)[i].Company = *comp
	}

	return uList, nil
}

func (s *Service) CreateUser(u *User) (*User, error) {
	newUser, err := s.userRepo.Insert(u)
	if err != nil {
		return nil, fmt.Errorf("user_service: create user: %w", err)
	}

	comp, err := s.companyRepo.FindById(newUser.Company.Id)
	if err != nil {
		return nil, fmt.Errorf("user_service: create user (load company): %w", err)
	}

	newUser.Company = *comp
	return newUser, nil
}

func (s *Service) UpdateUser(id int, newData *User) (*User, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("user_service: update user (find by id): %w", err)
	}

	user.Email = newData.Email
	user.Name = newData.Name
	user.Role = newData.Role
	user.Company = newData.Company
	user.hash = newData.hash

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("user_service: update user: %w", err)
	}

	comp, err := s.companyRepo.FindById(updatedUser.Company.Id)
	if err != nil {
		return nil, fmt.Errorf("user_service: update user (load company): %w", err)
	}

	updatedUser.Company = *comp
	return updatedUser, nil
}

func (s *Service) DeleteUser(id int) error {
	exists, err := s.userRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("user_service: delete user (check exists): %w", err)
	}

	if !exists {
		return fmt.Errorf("user_service: delete user: %w", ErrUserNotFound)
	}

	if err := s.userRepo.DeleteById(id); err != nil {
		return fmt.Errorf("user_service: delete user: %w", err)
	}

	return nil
}
