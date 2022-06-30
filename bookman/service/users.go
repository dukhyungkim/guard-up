package service

import (
	"bookman/entity"
	"bookman/events"
	"bookman/repository"
)

type UserService interface {
	SaveNewUser(user *entity.User) (*entity.User, error)
	ListUsers(pagination *entity.Pagination) ([]*entity.User, error)
	GetUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(userID int) error
}

type userService struct {
	userRepo    repository.UserRepo
	eventSender events.EventSender
}

func NewUserService(userRepo repository.UserRepo, eventSender events.EventSender) UserService {
	return &userService{
		userRepo:    userRepo,
		eventSender: eventSender,
	}
}

func (s *userService) SaveNewUser(user *entity.User) (*entity.User, error) {
	saveUser, err := s.userRepo.SaveUser(user)
	if err != nil {
		return nil, err
	}

	defer func() {
		go s.eventSender(events.EventAddedUser, saveUser)
	}()

	return saveUser, nil
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

	updateUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	defer func() {
		go s.eventSender(events.EventUpdatedUser, updateUser)
	}()

	return updateUser, nil
}

func (s *userService) DeleteUser(userID int) error {
	user := &entity.User{ID: userID}
	foundUser, err := s.userRepo.FetchUser(user)
	if err != nil {
		return err
	}

	err = s.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}

	defer func() {
		go s.eventSender(events.EventDeletedUser, foundUser)
	}()

	return nil
}
