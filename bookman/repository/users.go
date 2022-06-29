package repository

import (
	"bookman/common"
	"bookman/config"
	"bookman/entity"
	"errors"

	"gorm.io/gorm"
)

type UserRepo interface {
	SaveUser(user *entity.User) (*entity.User, error)
	FetchUser(user *entity.User) (*entity.User, error)
	ListUsers(pagination *entity.Pagination) ([]*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(user *entity.User) error
}

type userRepo struct {
	repo repo[entity.User]
}

func NewUserRepo(cfg *config.RDB) (UserRepo, error) {
	if err := initClient(cfg); err != nil {
		return nil, err
	}
	return &userRepo{}, nil
}

func (r *userRepo) SaveUser(user *entity.User) (*entity.User, error) {
	return r.repo.Save(user)
}

func (r *userRepo) FetchUser(user *entity.User) (*entity.User, error) {
	user, err := r.repo.First(user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundUser(err)
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepo) ListUsers(pagination *entity.Pagination) ([]*entity.User, error) {
	return r.repo.Find(pagination)
}

func (r *userRepo) UpdateUser(user *entity.User) (*entity.User, error) {
	return r.repo.Update(user)
}

func (r *userRepo) DeleteUser(user *entity.User) error {
	return r.repo.Delete(user)
}
