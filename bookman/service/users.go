package service

import (
	"bookman/entity"
	"bookman/repository"
)

type UserService interface {
	SaveNewUser(user *entity.User) (*entity.User, error)
	ListUsers(pagination *entity.Pagination) ([]*entity.User, error)
	GetUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(user *entity.User) error
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) SaveNewUser(user *entity.User) (*entity.User, error) {
	return s.userRepo.SaveUser(user)
}

func (s *userService) ListUsers(pagination *entity.Pagination) ([]*entity.User, error) {
	return s.userRepo.ListUsers(pagination)
}

func (s *userService) GetUser(user *entity.User) (*entity.User, error) {
	return s.userRepo.FetchUser(user)
}

func (s *userService) UpdateUser(user *entity.User) (*entity.User, error) {
	_, err := s.userRepo.FetchUser(&entity.User{ID: user.ID})
	if err != nil {
		return nil, err
	}
	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(user *entity.User) error {
	_, err := s.userRepo.FetchUser(&entity.User{ID: user.ID})
	if err != nil {
		return err
	}
	return s.userRepo.DeleteUser(user)
}
