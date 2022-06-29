package service

import (
	"bookman/entity"
	"bookman/repository"
	"time"
)

type RentalService interface {
	StartRentBook(bookID, userID int) (*entity.RentalStatus, error)
	GetRentStatus(bookID int) (*entity.RentalStatus, error)
	EndRentBook(bookID, userID int) error
}

type rentalService struct {
	rentalRepo repository.RentalRepo
}

func NewRentalService(rentalRepo repository.RentalRepo) RentalService {
	return &rentalService{rentalRepo: rentalRepo}
}

func (s *rentalService) StartRentBook(bookID, userID int) (*entity.RentalStatus, error) {
	rental := entity.RentalStatus{
		BookID: bookID,
		UserID: userID,
		Status: "대여",
		Start:  time.Now(),
	}
	return s.rentalRepo.SaveRental(&rental)
}

func (s *rentalService) GetRentStatus(bookID int) (*entity.RentalStatus, error) {
	return s.rentalRepo.FetchRentalStatusByBookID(bookID)
}

func (s *rentalService) EndRentBook(bookID, userID int) error {
	rental := entity.RentalStatus{
		BookID: bookID,
		UserID: userID,
	}

	if _, err := s.rentalRepo.FetchRentalStatus(&rental); err != nil {
		return err
	}

	return s.rentalRepo.DeleteRental(&rental)
}
