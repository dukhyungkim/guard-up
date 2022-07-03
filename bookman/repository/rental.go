package repository

import (
	"bookman/common"
	"bookman/config"
	"bookman/entity"
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var rentalRepoInstance RentalRepo

type RentalRepo interface {
	SaveRental(rental *entity.RentalStatus) (*entity.RentalStatus, error)
	FetchRentalStatusByBookID(bookID int) (*entity.RentalStatus, error)
	FetchRentalStatus(rental *entity.RentalStatus) (*entity.RentalStatus, error)
	DeleteRental(rental *entity.RentalStatus) error
}

type rentalRepo struct {
	repo repo[entity.RentalStatus]
}

func NewRentalRepo(cfg *config.RDB) (RentalRepo, error) {
	if rentalRepoInstance != nil {
		return rentalRepoInstance, nil
	}

	if err := initClient(cfg); err != nil {
		return nil, err
	}

	rentalRepoInstance = &rentalRepo{}
	return rentalRepoInstance, nil
}

func (r *rentalRepo) SaveRental(rental *entity.RentalStatus) (*entity.RentalStatus, error) {
	newRental, err := r.repo.Save(rental)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case fkConstraint:
				return nil, common.ErrNotFoundBookOrUser(err)
			case duplicateKey:
				return nil, common.ErrStartRent(err)
			default:
				return nil, common.ErrInternal(err)
			}
		}
		return nil, err
	}
	return newRental, nil
}

func (r *rentalRepo) FetchRentalStatusByBookID(bookID int) (*entity.RentalStatus, error) {
	var rentalStatus entity.RentalStatus
	if err := db.Where(map[string]any{"book_id": bookID}).First(&rentalStatus).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundRentalStatus(err)
		}
		return nil, err
	}
	return &rentalStatus, nil
}

func (r *rentalRepo) FetchRentalStatus(rental *entity.RentalStatus) (*entity.RentalStatus, error) {
	found, err := r.repo.First(rental)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundRentalStatus(err)
		}
		return nil, err
	}
	return found, nil
}

func (r *rentalRepo) DeleteRental(rental *entity.RentalStatus) error {
	return r.repo.Delete(rental)
}
